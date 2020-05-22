package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func main() {
	lambda.Start(Handler)
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
		case "/api/{userID}/emails":
			emails, err := listEmails(ctx, event.PathParameters["userID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				emails,
			), nil

		case "/api/{userID}/email/{emailID}":
			email, err := getEmailByID(ctx, event.PathParameters["userID"], event.PathParameters["emailID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				email,
			), nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Hello world! Go Boom From: " + domain + " path: " + event.Path,
	}, nil
}

func buildErrorResponse(ctx context.Context, err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf("%s", err),
	}
}

func buildOKResponse(ctx context.Context, cache bool, headers map[string]string, body string) events.APIGatewayProxyResponse {
	finalHeaders := headers
	if !cache {
		finalHeaders["Cache-Control"] = "private,s-maxage=15"
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    finalHeaders,
		Body:       body,
	}
}
