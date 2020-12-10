package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func authorizeRequest(next lambdaHandlerFunc) lambdaHandlerFunc {
	return lambdaHandlerFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		timestamp := request.Headers["x-slack-request-timestamp"]
		signature := request.Headers["x-slack-signature"]

		now := time.Now()
		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}

		if math.Abs(float64(now.Sub(t).Milliseconds())) > 5000.0 {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
			}, nil
		}

		validationStr := fmt.Sprintf("v0:%s:%s", timestamp, request.Body)
		secret := os.Getenv("SIGNING_SECRET")

		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write([]byte(validationStr))
		sha := hex.EncodeToString(hash.Sum(nil))

		if "v0="+sha != signature {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
			}, nil
		}

		return next(ctx, request)
	})
}
