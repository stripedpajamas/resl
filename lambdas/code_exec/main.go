package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/stripedpajamas/resl/models"
)

var languageConfigs models.LanguageConfig

func writeCodeFile(fileName string, code []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := path.Join(dir, fileName)
	err = ioutil.WriteFile(filePath, code, 0755)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created file to execute code: %s", filePath)

	return filePath, nil
}

func runCode(fileName string, languageConfig models.LanguageProperties) (string, error) {
	cmd := fmt.Sprintf(languageConfig.RunCommand, fileName)
	fmt.Printf("Running command %s", cmd)

	runCmd := exec.Command(cmd)

	runOut, err := runCmd.Output()
	if err != nil {
		return "", err
	}

	return string(runOut), nil
}

func handleRequest(ctx context.Context, input models.CodeProcessRequest) (models.CodeProcessResponse, error) {
	fmt.Printf("Incoming request details, %s %s", input.Code, input.Language)

	languageConfig, ok := languageConfigs[input.Language]
	if !ok {
		return models.CodeProcessResponse{}, nil
	}

	file, err := writeCodeFile(languageConfig.FileName, []byte(input.Code))
	if err != nil {
		return models.CodeProcessResponse{}, err
	}

	output, err := runCode(file, languageConfig)
	if err != nil {
		return models.CodeProcessResponse{}, err
	}

	return models.CodeProcessResponse{
		Output: output,
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
