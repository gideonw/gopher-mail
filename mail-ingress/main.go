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
func Handler(ctx context.Context, event events.SNSEvent) (int, error) {
	invokeCount = invokeCount + 1

	err := parseEvent(ctx, event)
	if err != nil {
		return invokeCount, err
	}

	return invokeCount, nil
}

func main() {
	lambda.Start(Handler)

	// b, err := ioutil.ReadFile("../test-sns-payload.json")
	// if err != nil {
	// 	log.Print(err)
	// }

	// var t events.SNSEvent
	// err = json.Unmarshal(b, &t)
	// if err != nil {
	// 	log.Print(err)
	// }

	// parseEvent(t)
}
