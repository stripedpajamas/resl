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

	"github.com/stripedpajamas/resl/models"
	"github.com/stripedpajamas/resl/slack"
)

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
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambdaClient.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	input := lambdaClient.InvokeInput{
		FunctionName: aws.String("resl-lang"),
		Payload:      request,
	}

	output, err := client.Invoke(&input)
	if err != nil {
		slack.SendChannelResponse(request.ResponseURL, "Sorry! Unable to setup execution environment :(")
		log.Printf("Error while invoking code runner: %s\n", err.Error())
		return err
	}

	slack.SendChannelResponse(request.ResponseURL, wrapString(escapeString(output.Payload)))

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
