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
	"github.com/gorilla/schema"
	"github.com/stripedpajamas/resl/models"
	"github.com/stripedpajamas/resl/slack"
)

var languageConfig models.LanguageConfig

var decoder = schema.NewDecoder()

// SlackRequestBody represents the incoming request body from Slack
type SlackRequestBody struct {
	APIAppID            string `schema:"api_app_id"`
	ChannelID           string `schema:"channel_id"`
	ChannelName         string `schema:"channel_name"`
	AppCommand          string `schema:"command"`
	IsEnterpriseInstall bool   `schema:"is_enterprise_install"`
	ResponseURL         string `schema:"response_url"`
	TeamDomain          string `schema:"team_domain"`
	TeamID              string `schema:"team_id"`
	Text                string `schema:"text"`
	Token               string `schema:"token"`
	TriggerID           string `schema:"trigger_id"`
	UserID              string `schema:"user_id"`
	UserName            string `schema:"user_name"`
}

func getCodePayloadFromRequestBody(requestBody SlackRequestBody) ([]byte, error) {
	trimmedText := strings.Trim(requestBody.Text, " ")
	spaceIdx := strings.IndexRune(trimmedText, ' ')

	if spaceIdx < 0 {
		return []byte{}, errors.New("failed to parse language and code")
	}

	language := trimmedText[0:spaceIdx]
	code := trimmedText[spaceIdx+1:]

	// confirm this language is supported
	var props models.LanguageProperties
	if _, found := languageConfig[language]; !found {
		return []byte{}, errors.New("language not supported")
	}

	props = languageConfig[language]

	// clean up slack's auto replacements
	code = strings.ReplaceAll(code, "&amp;", "&")
	code = strings.ReplaceAll(code, "&lt;", "<")
	code = strings.ReplaceAll(code, "&gt;", ">")

	// remove backticks from code block
	i := 0
	j := len(code) - 1

	for i <= j && code[i] == code[j] && code[i] == '`' {
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
		ResponseURL: requestBody.ResponseURL,
		Code:        code,
		Props:       props,
	})
}

func parseFormData(body string) (SlackRequestBody, error) {
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return SlackRequestBody{}, err
	}

	form, err := url.ParseQuery(string(decoded))
	if err != nil {
		return SlackRequestBody{}, err
	}

	log.Printf("Request form %s\n", form)

	var payload SlackRequestBody
	err = decoder.Decode(&payload, form)
	if err != nil {
		return SlackRequestBody{}, err
	}

	return payload, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := request.Headers
	log.Printf("Headers %s\n", headers)

	log.Printf("Request body: %s\n", request.Body)

	body, err := parseFormData(request.Body)
	if err != nil {
		log.Printf("Error while parsing request: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	payload, err := getCodePayloadFromRequestBody(body)
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

	lambda.Start(authorizeRequest(handleRequest))
}
