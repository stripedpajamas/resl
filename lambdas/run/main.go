package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stripedpajamas/resl/models"
)

// RequestBody represents the model of the incoming request body
type RequestBody struct {
	TriggerID   string `json:"trigger_id"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

var languageConfig models.LanguageConfig

func processRequestBody(request events.APIGatewayProxyRequest) (RequestBody, error) {
	fmt.Printf("Incoming request body: %s", request.Body)

	var body RequestBody

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		fmt.Println(err)
		return body, err
	}

	return body, nil
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
	languages, err := models.ParseLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfig = languages

	lambda.Start(handleRequest)
}
