package main

import (
	"context"
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
var domain string
var postOfficeBucket string
var postOfficePrefix string

func init() {
	domain = os.Getenv("DOMAIN")
	postOfficeBucket = os.Getenv("POST_OFFICE_BUCKET")
	postOfficePrefix = os.Getenv("POST_OFFICE_PREFIX")

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
func Handler(ctx context.Context, event events.SNSEvent) (int, error) {
	invokeCount = invokeCount + 1
	var lastErr error

	for i := range event.Records {
		record := event.Records[i]
		err := parseEvent(ctx, record)
		if err != nil {
			log.Println(err)
			lastErr = err
		}
	}

	if lastErr != nil {
		return invokeCount, lastErr
	}

	return invokeCount, nil
}

func main() {
	lambda.Start(Handler)
}
