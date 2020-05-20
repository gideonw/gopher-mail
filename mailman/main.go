package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var invokeCount = 0
var s3Client *s3.Client
var addressRegex *regexp.Regexp

// ENV variables
var domain string
var mailboxBucket string
var mailboxPrefix string
var verifyHeader string
var verifyValue string

const pathPrefix = "/api"

func init() {
	// Set ENV variables on init
	domain = os.Getenv("DOMAIN")
	mailboxBucket = os.Getenv("MAILBOX_BUCKET")
	mailboxPrefix = os.Getenv("MAILBOX_PREFIX")
	verifyHeader = os.Getenv("CF_VERIFY_HEADER")
	verifyValue = os.Getenv("CF_VERIFY_VALUE")

	addressRegex = regexp.MustCompile(`[^a-zA-Z0-9\-_()*'.].*`)

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	// cfg.Region = endpoints.UsWest2RegionID

	s3Client = s3.New(cfg)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	buf, _ := json.Marshal(event)
	log.Println(string(buf))

	// verify request is from cloudfront
	if event.Headers[verifyHeader] != verifyValue {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil
	}

	switch event.HTTPMethod {
	case "GET":
		switch event.Resource {
		case "/api/emails":
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Hello world! From: " + domain + "\n\n" + event.Path,
			}, nil
		case "/api/email/{emailID}":
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Hello world! From: " + domain + "\n\n" + event.Path,
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Hello world! Go Boom From: " + domain + " path: " + event.Path,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
