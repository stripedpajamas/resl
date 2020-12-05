package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stripedpajamas/resl/models"
)

func handleRequest(ctx context.Context, input models.CodeProcessRequest) (models.CodeProcessResponse, error) {
	fmt.Printf("Incoming request details, %s %s", input.Code, input.Language)
	return models.CodeProcessResponse{}, nil
}

func main() {
	lambda.Start(handleRequest)
}
