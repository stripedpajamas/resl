package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stripedpajamas/resl/models"
	"github.com/stripedpajamas/resl/slack"
)

var languageConfig models.LanguageConfig

func getCodePayloadFromRequestBody(bodyValues url.Values) ([]byte, error) {
	if _, found := bodyValues["text"]; !found {
		return []byte{}, errors.New("failed to parse language and code: no text")
	}
	if _, found := bodyValues["response_url"]; !found {
		return []byte{}, errors.New("failed to parse language and code: no response url")
	}

	trimmedText := strings.Trim(bodyValues["text"][0], " ")
	spaceIdx := strings.IndexRune(trimmedText, ' ')

	if spaceIdx < 0 {
		return []byte{}, errors.New("failed to parse language and code")
	}

	language := trimmedText[0:spaceIdx]
	code := trimmedText[spaceIdx+1:]

	// confirm this language is supported
	var props models.LanguageProperties
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
	return json.Marshal(models.CodeProcessRequest{
		ResponseURL: bodyValues["response_url"][0],
		Code:        code,
		Props:       props,
	})
}

func parseFormData(body string) (url.Values, error) {
	form, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, err
	}
	values, err := url.ParseQuery(string(form))
	if err != nil {
		return nil, err
	}

	return values, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request body: %s\n", request.Body)

	bodyValues, err := parseFormData(request.Body)
	if err != nil {
		log.Printf("Error while parsing request: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	payload, err := getCodePayloadFromRequestBody(bodyValues)
	if err != nil {
		log.Printf("Error while parsing language and code from request: %s\n", err.Error())
		responseBody, serializationErr := slack.PrivateAcknowledgement(err.Error())

		if serializationErr != nil {
			log.Printf("Failed to serialize parsing error for Slack: %s\n", err.Error())
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, serializationErr
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	input := lambdaClient.InvokeInput{
		FunctionName:   aws.String(os.Getenv("SLACK_RESP_ARN")),
		Payload:        payload,
		InvocationType: aws.String("Event"),
	}

	if _, err = client.Invoke(&input); err != nil {
		log.Printf("Error while invoking the code process lambda: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	res, err := slack.PublicAcknowledgement()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	languages, err := models.ImportLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfig = languages

	lambda.Start(handleRequest)
}
