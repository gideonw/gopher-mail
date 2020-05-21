package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func getEmailByID(ctx context.Context, userID, ID string) (string, error) {
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(mailboxBucket),
		Key:    aws.String(mailboxPrefix + "/" + userID + "/" + ID),
	}

	result, err := s3Client.GetObjectRequest(getInput).Send(ctx)
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

		return "", err
	}

	// Create a payload with the messageID and a nested object for the email
	ret := make(map[string]interface{})
	ret["messageID"] = ID

	// body, err := ioutil.ReadAll(result.Body)
	email, err := parsemail.Parse(result.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//nest the email struct onto the returned json object
	ret["email"] = email

	buf, err := json.Marshal(ret)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(buf), nil
}

// EmailList represents a simple json object with a single key holding the emails
type EmailList map[string][]string

// listEmails in the user's mailbox sitting in S3.
func listEmails(ctx context.Context, userID string) (string, error) {
	prefix := mailboxPrefix + "/" + userID

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(mailboxBucket),
		Prefix: aws.String(prefix),
	}

	result, err := s3Client.ListObjectsV2Request(listInput).Send(ctx)
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

		return "", err
	}

	ret := make(EmailList)
	ret["emails"] = []string{}

	for i := range result.Contents {
		ret["emails"] = append(ret["emails"], strings.Replace(*result.Contents[i].Key, prefix+"/", "", 1))
	}

	buf, err := json.Marshal(ret)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(buf), nil
}
