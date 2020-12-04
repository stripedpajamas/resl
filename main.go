package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
)

// RequestBody represents the model of the incoming request body
type RequestBody struct {
	TriggerID   int    `json:"trigger_id"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

// CodePayload represents the request payload sent to the invoked language lambda ** MOVE TO SHARED PKG **
type CodePayload struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

func processRequestBody(request events.APIGatewayProxyRequest) (RequestBody, error) {
	fmt.Printf("Incoming request body: %s", request.Body)

	var body RequestBody

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		fmt.Println(err)
		return body, err
	}

	return body, nil
}

func getCodePayloadFromRequestBody(body RequestBody) ([]byte, error) {
	var payload = parseText(body.Text)

	return json.Marshal(payload)
}

func parseCode(payload *CodePayload) {
	// TODO
	// clean code string
}

func parseLanguage(payload *CodePayload) {
	// TODO
	// parse language
	// validate that it is supported
}

func parseText(text string) CodePayload {
	text = strings.Trim(text, " ")
	chars := []rune(text)
	textLen := len(chars)
	code := ""
	lang := text

	for idx, c := range text {
		if c == ' ' {
			return CodePayload{
				Code:     string(chars[idx+1 : textLen]),
				Language: string(chars[0:idx]),
			}
		}
	}

	return CodePayload{
		Code:     code,
		Language: lang,
	}
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body, err := processRequestBody(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	payload, err := getCodePayloadFromRequestBody(body)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	input := lambdaClient.InvokeInput{
		FunctionName: aws.String("resl-lang"),
		Payload:      payload,
	}
	output, err := client.Invoke(&input)

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(output.Payload),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
