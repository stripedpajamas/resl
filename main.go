package main

import (
        "fmt"
        "context"
        "github.com/aws/aws-lambda-go/lambda"
        "github.com/aws/aws-lambda-go/events"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
    fmt.Printf("Body size = %d.\n", len(request.Body))

    fmt.Println("Headers:")
    for key, value := range request.Headers {
        fmt.Printf("    %s: %s\n", key, value)
    }

    return events.APIGatewayProxyResponse {
            Body: request.Body,
            StatusCode: 200,
    }, nil
}

func main() {
        lambda.Start(handleRequest)
}
