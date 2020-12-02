package main

import (
        "fmt"
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-lambda-go/events"
)

type MyEvent struct {
        Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (events.APIGatewayProxyResponse, error) {
        return events.APIGatewayProxyResponse {
                Body: fmt.Sprintf("Hello %s!", name.Name ),
                StatusCode: 200,
                Headers: make(map[string]string),
                MultiValueHeaders: make(map[string][]string),
                IsBase64Encoded: false,
        }, nil
}

func main() {
        lambda.Start(HandleRequest)
}
