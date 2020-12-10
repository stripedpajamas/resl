package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func authorizeRequest(next lambdaHandlerFunc) lambdaHandlerFunc {
	return lambdaHandlerFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		timestamp := request.Headers["x-slack-request-timestamp"]
		signature := request.Headers["x-slack-signature"]

		now := time.Now().Unix()
		t, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}

		if math.Abs(float64(now)-float64(t)) > (60.0 * 5.0) {
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
