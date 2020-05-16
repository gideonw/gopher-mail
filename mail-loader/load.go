package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func parseEvent(ctx context.Context, event events.SNSEvent) error {
	for i := range event.Records {
		record := event.Records[i]
		log.Printf("[parseEvent] %s - %s", record.SNS.MessageID, record.SNS.Subject)
		// log.Printf("[parseEvent] %s - %s\n%s", record.SNS.MessageID, record.SNS.Subject, record.SNS.Message)

		parseEmail(ctx, record.SNS.Message)
	}

	return nil
}
func parseEmail(ctx context.Context, body []byte) error {

	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	if content, ok := msg["content"]; ok {
		if contentString, ok := content.(string); ok {
			stringReader := strings.NewReader(contentString)
			email, err := parsemail.Parse(stringReader)
			if err != nil {
				return err
			}

			err = saveToS3(ctx, email)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func saveToS3(ctx context.Context, email []byte) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("gideonw-com-email"),
	}
	req := s3Client.ListObjectsV2Request(input)
	res, err := req.Send(ctx)
	if err != nil {
		return err
	}

	objs := res.Contents
	log.Println(objs)

	return nil
}
