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

func createErrorResponse(code int, err error, message string) (events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = "Error found"
	}
	log.Printf("%s: %s\n", message, err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: code,
	}, err
}

func parseText(text string) (string, string) {
	trimmedText := strings.Trim(text, " ")
	spaceIdx := strings.IndexRune(trimmedText, ' ')

	if spaceIdx < 0 {
		return trimmedText, ""
	}

	return trimmedText[0:spaceIdx], trimmedText[spaceIdx+1:]
}

func getCodePayloadFromRequestBody(requestBody slack.Request) (models.CodeProcessRequest, error) {
	log.Printf("Request Body: %+v\n", requestBody)

	if requestBody.Text == "" {
		return models.CodeProcessRequest{
			ResponseURL: requestBody.ResponseURL,
			Code:        "",
			Props:       models.LanguageProperties{},
		}, nil
	}

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
		UserID:      requestBody.UserID,
	}, nil
}

func createRequestBodyFromModalPayload(payload slack.ModalRequest) (slack.Request, error) {
	formData := payload.View.State.Values

	codeElementVal, ok := formData[slack.CodeBlockName]
	if !ok {
		return slack.Request{}, errors.New("Code block not found")
	}

	language := payload.View.PrivateMetadata

	languageElementVal, ok := formData[slack.LanguageBlockName]
	if language == "" && !ok {
		return slack.Request{}, errors.New("No language provided")
	} else if ok && language == "" {
		languageInputVal, ok := languageElementVal[slack.LanguageActionID]
		if !ok {
			return slack.Request{}, errors.New("Language action not found")
		}

		languageInputJSON, err := json.Marshal(languageInputVal)
		if err != nil {
			return slack.Request{}, err
		}

		var languageInput slack.StaticSelectElement
		err = json.Unmarshal([]byte(languageInputJSON), &languageInput)
		if err != nil {
			return slack.Request{}, err
		}

		language = languageInput.SelectedOption.Value
	}

	codeinputVal, ok := codeElementVal[slack.CodeActionID]
	if !ok {
		return slack.Request{}, errors.New("Code action not found")
	}

	inputJSON, err := json.Marshal(codeinputVal)
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
		Text:        fmt.Sprintf("%s %s", language, codeInput.Value),
		ResponseURL: payload.ResponseURLS[0].URL,
		TriggerID:   payload.TriggerID,
		UserID:      payload.User.UserID,
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

	log.Printf("Parsed Body: %+v\n", body)

	var modalBody slack.ModalRequest
	isModal := false

	if body.ModalPayload != "" {
		isModal = true
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

	// fire a modal back since no code was there and modal is not alreay present
	if !isModal && codeProcessRequest.Code == "" {
		err = slack.SendModal(body.TriggerID, codeProcessRequest.Props.Name, codeProcessRequest.Props.ShortName, codeProcessRequest.Props.Placeholder)

		if err != nil {
			return createErrorResponse(500, err, "Failed to send modal")
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	codeProcessRequest.Modal = isModal

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

	var res []byte

	if isModal {
		res, err = slack.ClearModal()
	} else {
		res, err = slack.PublicAcknowledgement()
	}

	if err != nil {
		return createErrorResponse(500, err, "")
	}

	log.Printf("Responding to slack: %s", res)

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
