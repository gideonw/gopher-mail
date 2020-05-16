package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var invokeCount = 0
var s3Client *s3.Client

func init() {
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
}
