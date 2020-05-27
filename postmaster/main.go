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

var s3Client *s3.Client
var addressRegex *regexp.Regexp
var domain string
var mailboxBucket string
var mailboxPrefix string
var postOfficePrefix string

type emailToSort struct {
	MessageID string

	SourceBucket    string
	SourceObjectKey string

	DestPrefixes  []string
	DestObjectKey string

	Errored bool
}

func init() {
	domain = os.Getenv("DOMAIN")
	mailboxBucket = os.Getenv("MAILBOX_BUCKET")
	postOfficePrefix = os.Getenv("POST_OFFICE_PREFIX")
	mailboxPrefix = os.Getenv("MAILBOX_PREFIX")

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
func Handler(ctx context.Context, event events.SNSEvent) error {
	var lastErr error
	emailsToProcess := []emailToSort{}

	for i := range event.Records {
		record := event.Records[i]

		// Accumulate emails in the triggering event
		email, err := parseEvent(ctx, record)
		if err != nil {
			log.Println(err)
			lastErr = err
		}
		emailsToProcess = append(emailsToProcess, email)
	}

	// Retrieve all of the emails in the `_errored` folder for reprocessing
	erroredEmails, err := loadErroredEmails(ctx)
	if err != nil {
		log.Println(err)
		lastErr = err
	} else {
		emailsToProcess = append(emailsToProcess, erroredEmails...)
	}

	// Sort the emails into their mailboxes
	for i := range emailsToProcess {
		err = sortEmailIntoMailbox(ctx, emailsToProcess[i])
		if err != nil {
			log.Println(err)
			lastErr = err
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
