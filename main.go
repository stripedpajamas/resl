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
)

type code struct {
	Text     string
	Language string
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	fmt.Println("Headers:")
	for key, value := range request.Headers {
		fmt.Printf("    %s: %s\n", key, value)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	payload, err := json.Marshal(code{Text: "console.log('hello')", Language: "javascript"})

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	output, err = client.Invoke(&lambdaClient.InvokeInput{FunctionName: aws.String("ReslLangLambda"), Payload: payload})

	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       output.Payload,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
