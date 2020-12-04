package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type CodeInput struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type CodeOutput struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}

func handleRequest(ctx context.Context, input CodeInput) (CodeOutput, error) {
	fmt.Printf("Incoming request details, %s %s", input.Code, input.Language)
	return CodeOutput{}, nil
}

func main() {
	lambda.Start(handleRequest)
}
