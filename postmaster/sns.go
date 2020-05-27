package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func parseEvent(ctx context.Context, record events.SNSEventRecord) (emailToSort, error) {
	log.Printf("[parseEvent] %s - %s\n", record.SNS.MessageID, record.SNS.Subject)
	eventEmail := emailToSort{
		Errored: true,
	}

	msg, err := parseJSON(ctx, record.SNS.Message)
	if err != nil {
		log.Println("Error parsing mailBody as JSON")
		return eventEmail, fmt.Errorf("%s", "Error parsing email event")
	}

	eventEmail.MessageID, err = getMessageID(ctx, msg)
	if err != nil {
		log.Println("Error failed to parse incoming new email SNS event")
		log.Println(record.SNS.Message)
		log.Println(err)
		return eventEmail, err
	}

	eventEmail.SourceBucket, eventEmail.SourceObjectKey, err = getS3SourcePath(ctx, msg)
	if err != nil {
		log.Println("Error failed to parse incoming new email SNS event")
		log.Println(record.SNS.Message)
		log.Println(err)
		return eventEmail, err
	}

	eventEmail.DestPrefixes, eventEmail.DestObjectKey, err = getS3DestinationPath(ctx, msg)
	if err != nil {
		log.Println(err)
		return eventEmail, err
	}

	eventEmail.Errored = false
	return eventEmail, nil
}

func parseJSON(ctx context.Context, body string) (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func getMessageID(ctx context.Context, msg map[string]interface{}) (string, error) {
	mail, ok := msg["mail"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("%s", "Error parsing mail object while retrieving messageId")
		log.Println(err)
		return "", err
	}

	messageID, ok := mail["messageId"].(string)
	if !ok {
		err := fmt.Errorf("%s", "Error asserting messageId")
		log.Println(err)
		return "", err
	}

	return messageID, nil
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
