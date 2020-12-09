package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stripedpajamas/resl/models"
)

var languageConfig models.LanguageConfig

// RequestBody represents the model of the incoming request body
type RequestBody struct {
	TriggerID   string `json:"trigger_id"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

// SlackResponse represents the model slack expects
type SlackResponse struct {
	ResponseType string `json:"response_type"`
}

func getCodePayloadFromRequestBody(body RequestBody) ([]byte, error) {
	trimmedText := strings.Trim(body.Text, " ")
	spaceIdx := strings.IndexRune(trimmedText, ' ')

	if spaceIdx < 0 {
		return []byte{}, errors.New("failed to parse language and code")
	}

	language := trimmedText[0:spaceIdx]
	code := trimmedText[spaceIdx+1:]

	// confirm this language is supported
	var props models.LanguageProperties
	if langProps, found := languageConfig[language]; !found {
		return []byte{}, errors.New("language not supported")
	} else {
		props = langProps
	}

	// clean up slack's auto replacements
	code = strings.ReplaceAll(code, "&amp;", "&")
	code = strings.ReplaceAll(code, "&lt;", "<")
	code = strings.ReplaceAll(code, "&gt;", ">")

	// remove backticks from code block
	i := 0
	j := len(code) - 1

	for i <= j && code[i] == '`' && code[j] == code[i] {
		i++
		j--
	}

	if i == 1 || i == 3 {
		code = code[i : j+1]
	}

	log.Printf("Parsed Code: %s\n", code)
	log.Printf("Parsed Language: %s\n", language)

	// json stringify the result for the execution lambda
	return json.Marshal(models.CodeProcessRequest{
		ResponseURL: body.ResponseURL,
		Code:        code,
		Props:       props,
	})
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Printf("Error while parsing request: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	payload, err := getCodePayloadFromRequestBody(body)
	if err != nil {
		log.Printf("Error while parsing language and code from request: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, err
	}

	input := lambdaClient.InvokeInput{
		FunctionName:   aws.String("resl-lang"),
		Payload:        payload,
		InvocationType: aws.String("event"),
	}

	if _, err = client.Invoke(&input); err != nil {
		log.Printf("Error while invoking the code process lambda: %s\n", err.Error())
	}

	res, err := json.Marshal(SlackResponse{ResponseType: "in_channel"})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: 200,
	}, nil
}

func main() {
	languages, err := models.ImportLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfig = languages

	lambda.Start(handleRequest)
}
