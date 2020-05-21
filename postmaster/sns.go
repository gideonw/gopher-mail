package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func parseEvent(ctx context.Context, record events.SNSEventRecord) error {
	log.Printf("[parseEvent] %s - %s\n", record.SNS.MessageID, record.SNS.Subject)

	msg, err := parseJSON(ctx, record.SNS.Message)
	if err != nil {
		log.Println("Error parsing mailBody as JSON")
		return fmt.Errorf("%s", "Error parsing email event")
	}

	srcBucket, srcObjectKey, err := getS3SourcePath(ctx, msg)
	if err != nil {
		log.Println("Error failed to parse incoming new email SNS event")
		log.Println(err)
		log.Println(record.SNS.Message)
		return err
	}

	destPaths, destObjectKey, err := getS3DestinationPath(ctx, msg)
	if err != nil {
		if destObjectKey != "" {
			destPaths = []string{
				"_errored",
			}
		} else {
			log.Println(err)
			return err
		}
	}

	err = sortEmailToMailbox(ctx, srcBucket, srcObjectKey, destPaths, destObjectKey)
	if err != nil {
		return err
	}

	log.Println("Finished processing " + srcObjectKey)

	return nil
}

func sortEmailToMailbox(ctx context.Context, srcBucket, srcObjectKey string, destPrefix []string, destObjectKey string) error {
	sourcePath := srcBucket + "/" + srcObjectKey
	var errList []error

	for _, prefix := range destPrefix {
		destKey := mailboxPrefix + "/" + prefix + "/" + destObjectKey

		copyInput := &s3.CopyObjectInput{
			CopySource: aws.String(url.PathEscape(sourcePath)),

			Bucket: aws.String(srcBucket),
			Key:    aws.String(destKey),

			ContentType: aws.String("application/octet-stream"),
		}
		log.Printf("Copying from \"%s\" to \"%s\"\n", sourcePath, srcBucket+"/"+destKey)

		copyResp, err := s3Client.CopyObjectRequest(copyInput).Send(ctx)
		if err != nil {
			log.Println(copyResp)
			log.Println(err)
			// don't stop on error, track th errors and also copy to the _errored folder if we can
			errList = append(errList, err)
			if len(errList) == 0 {
				destPrefix = append(destPrefix, "_errored")
			}
		}
	}

	// if we don't have an error copying the object we can delete the old one
	delInput := &s3.DeleteObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    aws.String(srcObjectKey),
	}
	log.Printf("Deleting object \"%s\"\n", sourcePath)

	delResp, err := s3Client.DeleteObjectRequest(delInput).Send(ctx)
	if err != nil {
		log.Println(delResp)
		return err
	}

	if len(errList) != 0 {
		return errList[0]
	}

	return nil
}

func getS3SourcePath(ctx context.Context, msg map[string]interface{}) (string, string, error) {
	// receipt.action.{bucketName, objectKey}
	receipt, ok := msg["receipt"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("%s", "Error parsing receipt object while retrieving src object info")
		log.Println(err)
		return "", "", err
	}

	action, ok := receipt["action"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("%s", "Error parsing action object while retrieving src object info")
		log.Println(err)
		return "", "", err
	}

	bucketName, ok := action["bucketName"].(string)
	if !ok {
		err := fmt.Errorf("%s", "Error asserting bucketName")
		log.Println(err)
		return "", "", err
	}

	objectKey, ok := action["objectKey"].(string)
	if !ok {
		err := fmt.Errorf("%s", "Error asserting objectKey")
		log.Println(err)
		return "", "", err
	}

	return bucketName, objectKey, nil

}

// getS3DestinationPath takes the message and extracts the fields required to compute the paths
// returns list of 'to' emails and new path
func getS3DestinationPath(ctx context.Context, msg map[string]interface{}) ([]string, string, error) {
	paths := []string{}
	filename := ""

	mailBody, ok := msg["mail"].(map[string]interface{})
	if !ok {
		log.Println("Error asserting mailBody to map[string]interface{}")
		return nil, "", fmt.Errorf("%s", "Error parsing email event while building filenames")
	}

	timestamp, ok := mailBody["timestamp"].(string)
	if !ok {
		log.Println("Error asserting timestamp")
		return nil, "", fmt.Errorf("%s", "Error parsing email event while building filenames")
	}

	messageID, ok := mailBody["messageId"].(string)
	if !ok {
		log.Println("Error asserting messageID")
		return nil, "", fmt.Errorf("%s", "Error parsing email event while building filenames")
	}

	// use the messageID as the file name since we want to use the ID to request the emails
	filename = messageID

	if emailAddresses, ok := mailBody["destination"].([]interface{}); ok {
		paths = []string{}
		for _, address := range emailAddresses {
			addressString, ok := address.(string)
			if !ok {
				continue
			}

			userLen := strings.LastIndex(addressString, "@")
			user := addressRegex.ReplaceAllString(addressString[0:userLen], "-")
			addressDomain := addressString[userLen+1:]
			if domain == addressDomain {
				paths = append(paths, user)
			}
		}

		if len(paths) == 0 {
			return nil, filename, fmt.Errorf("%s", "No emails match our root domain")
		}
	} else {
		return nil, filename, fmt.Errorf("%s", "Error asserting destination", mailBody["destination"])
	}

	return paths, filename, nil
}

func parseJSON(ctx context.Context, body string) (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
