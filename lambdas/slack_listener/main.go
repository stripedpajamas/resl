package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

// RequestBodyParser reprsents a function type for parsing the incoming request body
type RequestBodyParser func(body string) (slack.Request, error)

func parseText(text string) (string, string) {
	trimmedText := strings.Trim(text, " ")
	spaceIdx := strings.IndexRune(trimmedText, ' ')

	if spaceIdx < 0 {
		return trimmedText, ""
	}

	return trimmedText[0:spaceIdx], trimmedText[spaceIdx+1:]
}

func createErrorResponse(code int, err error, message string) (events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = "Error found"
	}
	log.Printf("%s: %s\n", message, err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: code,
	}, err
}

func getCodePayloadFromRequestBody(requestBody slack.Request) (models.CodeProcessRequest, error) {
	language, code := parseText(requestBody.Text)

	props, found := languageConfig[language]
	if !found {
		return models.CodeProcessRequest{}, errors.New("language not supported")
	}

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
	return models.CodeProcessRequest{
		ResponseURL: requestBody.ResponseURL,
		Code:        code,
		Props:       props,
	}, nil
}

func createRequestBodyFromModalPayload(payload slack.ModalRequest) (slack.Request, error) {
	formData := payload.View.State.Values

	val, ok := formData[slack.CodeBlockName]
	if !ok {
		return slack.Request{}, errors.New("Code block not found")
	}

	inputJSON, err := json.Marshal(val)
	if err != nil {
		return slack.Request{}, err
	}

	var codeInput slack.InputElement
	err = json.Unmarshal([]byte(inputJSON), &codeInput)
	if err != nil {
		return slack.Request{}, err
	}

	if codeInput.Value == "" {
		return slack.Request{}, errors.New("No code present in the modal input")
	}

	if payload.ResponseURLS == nil || len(payload.ResponseURLS) == 0 {
		return slack.Request{}, errors.New("No response urls available in modal request")
	}

	return slack.Request{
		Text:        fmt.Sprintf("%s %s", payload.View.PrivateMetadata, codeInput.Value),
		ResponseURL: payload.ResponseURLS[0].URL,
	}, nil
}

func parseFormRequest(body string) (slack.Request, error) {
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return slack.Request{}, err
	}

	form, err := url.ParseQuery(string(decoded))
	if err != nil {
		return slack.Request{}, err
	}

	log.Printf("Request form %s\n", form)

	var payload slack.Request
	err = decoder.Decode(&payload, form)
	if err != nil {
		return slack.Request{}, err
	}

	return payload, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body, err := parseFormRequest(request.Body)
	if err != nil {
		return createErrorResponse(500, err, "Error while parsing request")
	}

	log.Printf("%+v\n", body)

	log.Printf(body.ModalPayload)

	var modalBody slack.ModalRequest

	if body.ModalPayload != "" {
		err := json.Unmarshal([]byte(body.ModalPayload), &modalBody)
		if err != nil {
			return createErrorResponse(500, err, "Error while parsing modal body")
		}

		body, err = createRequestBodyFromModalPayload(modalBody)
		if err != nil {
			return createErrorResponse(400, err, "Error while processing modal body")
		}
	}

	codeProcessRequest, err := getCodePayloadFromRequestBody(body)
	if err != nil {
		log.Printf("Error while parsing language and code from request: %s\n", err.Error())
		responseBody, serializationErr := slack.PrivateAcknowledgement(err.Error())

		if serializationErr != nil {
			return createErrorResponse(500, err, "Failed to serialize parsing error for Slack")
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(responseBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// fire a modal back since no code was there
	if codeProcessRequest.Code == "" {
		err = slack.SendModal(body.TriggerID, codeProcessRequest.Props.Name, codeProcessRequest.Props.ShortName, codeProcessRequest.Props.Placeholder)

		if err != nil {
			return createErrorResponse(500, err, "Failed to send modal")
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	payload, err := json.Marshal(codeProcessRequest)
	if err != nil {
		return createErrorResponse(500, err, "Failed to serialize parsing error for Slack")
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
		return createErrorResponse(500, err, "Error while invoking the code process lambda")
	}

	res, err := slack.PublicAcknowledgement()
	if err != nil {
		return createErrorResponse(500, err, "")
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
	decoder.IgnoreUnknownKeys(true)
	languages, err := models.ImportLanguageConfig("languages.json")
	if err != nil {
		panic(err)
	}

	languageConfig = languages

	lambda.Start(authorizeRequest(handleRequest))
}
