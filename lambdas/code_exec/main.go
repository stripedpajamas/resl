package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stripedpajamas/resl/models"
)

var languageConfigs models.LanguageConfig

func handleRequest(ctx context.Context, input models.CodeProcessRequest) (models.CodeProcessResponse, error) {
	fmt.Printf("Incoming request details, %s %s", input.Code, input.Language)

	languageConfig, ok := languageConfigs[input.Language]
	if !ok {
		return models.CodeProcessResponse{}, nil
	}

	file, err := writeCodeFile(languageConfig.FileName, []byte(input.Code))
	if err != nil {
		return models.CodeProcessResponse{
			Error: err,
		}, err
	}

	languageConfig.FileName = file

	output, err := runCode(languageConfig)
	if err != nil {
		return models.CodeProcessResponse{
			Error: err,
		}, err
	}

	return models.CodeProcessResponse{
		Output: output,
		Error:  nil,
	}, nil
}

func main() {
	languages, err := models.ParseLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfigs = languages

	lambda.Start(handleRequest)
}
