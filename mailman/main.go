package main

import (
	"context"
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
var domain string

const pathPrefix = "/api"

func init() {
	domain = os.Getenv("DOMAIN")
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
	switch event.HTTPMethod {
	case "GET":
		switch event.Path {
		case pathPrefix + "/email":
			return events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "Hello world! From: " + domain,
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
