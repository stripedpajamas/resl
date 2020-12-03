package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CodeInput struct {
	Text     string
	Language string
}

type CodeOutput struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}

func handleRequest(ctx context.Context, input CodeInput) (CodeOutput, error) {
	return CodeOutput{}
}

func main() {
	lambda.Start(handleRequest)
}
