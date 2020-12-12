package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"

	"github.com/stripedpajamas/resl/models"
	"github.com/stripedpajamas/resl/slack"
)

type CodeOutput struct {
	Output string `json:"output"`
}

// replaces backticks with \`
func escapeString(s string) string {
	return strings.ReplaceAll(s, "`", "\\`")
}

// wraps a string in ```<string>```
func wrapString(s string) string {
	var b strings.Builder
	b.WriteString("```")
	b.WriteString(s)
	b.WriteString("```")
	return b.String()
}

func handleRequest(ctx context.Context, request models.CodeProcessRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	input := lambdaClient.InvokeInput{
		FunctionName: aws.String(os.Getenv("CODE_EXEC_LAMBDA_ARN")),
		Payload:      payload,
	}

	log.Printf("Invoking code exec lambda...\n")
	output, err := client.Invoke(&input)
	if err != nil {
		slack.SendChannelResponse(request.ResponseURL, "Sorry! Unable to setup execution environment :(")
		log.Printf("Error while invoking code runner: %s\n", err.Error())
		return err
	}

	var codeOutput CodeOutput
	err = json.Unmarshal(output.Payload, &codeOutput)
	if err != nil {
		slack.SendChannelResponse(request.ResponseURL, "Sorry! Unable to setup execution environment :(")
		log.Printf("Error while deserializing code output: %s\n", err.Error())
		return err
	}

	fmt.Println(codeOutput.Output)

	log.Printf("Sending slack response...\n")

	if codeOutput.Output == "" {
		codeOutput.Output == "[No output]"
	}

	slack.SendChannelResponse(request.ResponseURL, wrapString(escapeString(string(codeOutput.Output))))

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
