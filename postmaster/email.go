package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func sortEmailIntoMailbox(ctx context.Context, email emailToSort) error {
	var errList []error

	for _, prefix := range email.DestPrefixes {
		destObjectKey := mailboxPrefix + "/" + prefix + "/" + email.DestObjectKey
		err := processEmail(ctx, email.MessageID, email.SourceBucket, email.SourceObjectKey, destObjectKey)
		if err != nil {
			// Log the error and mark the email as errored
			log.Println(err)
			email.Errored = true
		}
	}

	// Attempt to copy errored emails to _errored bucket.
	if email.Errored {
		err := copyErroredEmail(ctx, email.SourceBucket, email.SourceObjectKey, email.DestObjectKey)
		if err != nil {
			log.Println("Failed to copy errored email to the '_errored' mailbox, skipping delete")
			log.Println(email)

			// TODO: IDK, returning here will cause trash to pile up in the post-office folder,
			return err
		}
	}

	// if we don't have an error copying the object we can delete the old one
	delInput := &s3.DeleteObjectInput{
		Bucket: aws.String(email.SourceBucket),
		Key:    aws.String(email.SourceObjectKey),
	}
	log.Printf("Deleting object \"%s\"\n", email.SourceBucket+"/"+email.SourceObjectKey)

	delResp, err := s3Client.DeleteObjectRequest(delInput).Send(ctx)
	if err != nil {
		log.Println(delResp)
		return err
	}

	if len(errList) != 0 {
		return errList[0]
	}

	log.Println("Finished processing " + email.SourceObjectKey)

	return nil
}

func processEmail(ctx context.Context, messageID, srcBucket, srcObjectKey, destObjectKey string) error {
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    aws.String(srcObjectKey),
	}
	log.Printf("Getting raw email from \"%s\"\n", srcBucket+"/"+srcObjectKey)

	result, err := s3Client.GetObjectRequest(getInput).Send(ctx)
	if checkAwsErr(err) != nil {
		return err
	}
	bodyBuf, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println("Failed to read the contents of the raw email")
		return err
	}

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(mailboxBucket),
		Key:    aws.String(destObjectKey),

		Body:        bytes.NewReader(bodyBuf),
		ContentType: aws.String("application/octet-stream"),
	}
	log.Printf("Writing raw email to \"%s\"\n", mailboxBucket+"/"+destObjectKey)

	putResp, err := s3Client.PutObjectRequest(putInput).Send(ctx)
	if checkAwsErr(err) != nil {
		log.Println(putResp)
		return err
	}

	// Create a payload with the messageID and a nested object for the email
	email, err := parsemail.Parse(bytes.NewReader(bodyBuf))
	if err != nil {
		log.Println(err)
		return err
	}

	//nest the email struct onto the json object
	emailMeta := make(map[string]interface{})
	emailMeta["MessageID"] = messageID
	emailMeta["Email"] = email

	buf, err := json.Marshal(emailMeta)
	if err != nil {
		log.Println(err)
		return err
	}
	bodyReader := strings.NewReader(string(buf))

	putInput = &s3.PutObjectInput{
		Bucket: aws.String(mailboxBucket),
		Key:    aws.String(destObjectKey + ".json"),

		Body:        bodyReader,
		ContentType: aws.String("application/json"),
	}
	log.Printf("Writing json email to \"%s\"\n", mailboxBucket+"/"+destObjectKey+".json")

	putResp, err = s3Client.PutObjectRequest(putInput).Send(ctx)
	if checkAwsErr(err) != nil {
		log.Println(putResp)
		return err
	}

	return nil
}

func loadErroredEmails(ctx context.Context) ([]emailToSort, error) {
	ret := []emailToSort{}

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(mailboxBucket),
		Prefix: aws.String(mailboxPrefix + "/_errored/"),
	}
	log.Printf("Loading '_errored' emails from \"%s/%s/_errored/\"\n", mailboxBucket, mailboxPrefix)

	result, err := s3Client.ListObjectsV2Request(listInput).Send(ctx)
	if checkAwsErr(err) != nil {
		log.Println(result)
		return ret, err
	}

	for i := range result.Contents {
		// TODO: Load errored email, parse it for the To addresses, and then set the Dest fields on emailToSort

		ret = append(ret, emailToSort{
			MessageID: *result.Contents[i].Key,

			SourceBucket:    mailboxBucket,
			SourceObjectKey: *result.Contents[i].Key,
		})
	}

	return ret, nil
}

// copyErroredEmail is called when an email has failed and we want to copy it into the `_errored` mailbox
func copyErroredEmail(ctx context.Context, srcBucket, srcObjectKey, destKey string) error {
	sourcePath := srcBucket + "/" + srcObjectKey
	destBucketPath := mailboxPrefix + "/_errored/" + destKey

	copyInput := &s3.CopyObjectInput{
		CopySource: aws.String(url.PathEscape(sourcePath)),

		Bucket: aws.String(mailboxBucket),
		Key:    aws.String(destBucketPath),

		ContentType: aws.String("application/octet-stream"),
	}
	log.Printf("Copying from \"%s\" to \"%s\"\n", sourcePath, mailboxBucket+"/"+destBucketPath)

	copyResp, err := s3Client.CopyObjectRequest(copyInput).Send(ctx)
	if checkAwsErr(err) != nil {
		log.Println(copyResp)
		return err
	}

	return nil
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
