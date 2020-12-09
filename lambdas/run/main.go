package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
)

var languageConfig LanguageConfig

// LanguageConfig is a map of language to props
type LanguageConfig map[string]LanguageProperties

// LanguageProperties represents properties for running each supported language
type LanguageProperties struct {
	Name           string `json:"langName"`
	Extension      string `json:"extension"`
	Placeholder    string `json:"placeholder"`
	RunCommand     string `json:"runCmd"`
	CompileCommand string `json:"compileCmd"`
}

// RequestBody represents the model of the incoming request body
type RequestBody struct {
	TriggerID   string `json:"trigger_id"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

// CodeProcessRequest represents the payload sent to the code runner lambda
type CodeProcessRequest struct {
	Code     string             `json:"code"`
	Language string             `json:"language"`
	Props    LanguageProperties `json:"props"`
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
	var props LanguageProperties
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
	return json.Marshal(CodeProcessRequest{
		Language: language,
		Code:     code,
		Props:    props,
	})
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body RequestBody

	if err := json.Unmarshal([]byte(request.Body), &body); err != nil {
		log.Printf("Error while parsing request: %s\n", err.Error())
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
		log.Printf("Error while parsing language and code from request: %s\n", err.Error())
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
		log.Printf("Error while invoking code runner: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(output.Payload),
		StatusCode: 200,
	}, nil
}

func importLanguageConfig(filePath string) (LanguageConfig, error) {
	var config LanguageConfig

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path.Join(dir, filePath))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	languages, err := importLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfig = languages

	lambda.Start(handleRequest)
}
