package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getEmailByID(ctx context.Context, userID, ID string) (string, error) {
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(mailboxBucket),
		Key:    aws.String(mailboxPrefix + "/" + userID + "/" + ID + ".json"),
	}

	result, err := s3Client.GetObjectRequest(getInput).Send(ctx)
	if checkAwsErr(err) != nil {
		return "", err
	}

	buf, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

type emailStorage struct {
	MessageID string
	Email     parsemail.Email
}

// EmailMeta contians a snapshot of an email for the frontend
type EmailMeta struct {
	MessageID string
	Subject   string
	Date      time.Time
}

// EmailMetaList represents a simple json object with a single key holding the emails
type EmailMetaList map[string][]EmailMeta

// listEmails in the user's mailbox sitting in S3, as JSON
func listEmails(ctx context.Context, userID string) (string, error) {
	prefix := mailboxPrefix + "/" + userID

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(mailboxBucket),
		Prefix: aws.String(prefix),
	}

	result, err := s3Client.ListObjectsV2Request(listInput).Send(ctx)
	if checkAwsErr(err) != nil {
		return "", err
	}

	ret := make(EmailMetaList)
	ret["emails"] = []EmailMeta{}

	for i := range result.Contents {
		if strings.HasSuffix(*result.Contents[i].Key, ".json") {
			result, err := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(mailboxBucket),
				Key:    aws.String(*result.Contents[i].Key),
			}).Send(ctx)
			if checkAwsErr(err) != nil {
				return "", err
			}

			bodyBuf, err := ioutil.ReadAll(result.Body)
			if err != nil {
				return "", err
			}

			var email emailStorage
			err = json.Unmarshal(bodyBuf, &email)
			if err != nil {
				return "", err
			}

			ret["emails"] = append(ret["emails"], EmailMeta{
				MessageID: email.MessageID,
				Subject:   email.Email.Subject,
				Date:      email.Email.Date,
			})
		}
	}

	buf, err := json.Marshal(ret)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(buf), nil
}

func checkAwsErr(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}

		return err
	}

	return nil
}
