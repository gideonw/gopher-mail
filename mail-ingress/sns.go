package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func parseEvent(ctx context.Context, event events.SNSEvent) error {
	for i := range event.Records {
		record := event.Records[i]
		log.Printf("[parseEvent] %s - %s", record.SNS.MessageID, record.SNS.Subject)

		paths, filename := getS3StoragePath(ctx, record)
		for _, path := range paths {

			err := saveToS3(ctx, path+"/"+filename, record.SNS.Message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func saveToS3(ctx context.Context, filename string, email string) error {
	stringReader := strings.NewReader(email)

	input := &s3.PutObjectInput{
		Bucket:      aws.String("gideonw-com-email"),
		Key:         aws.String(filename),
		Body:        stringReader,
		ContentType: aws.String("application/json"),
	}
	req := s3Client.PutObjectRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return err
	}

	log.Println(res)

	return nil
}

func getS3StoragePath(ctx context.Context, event events.SNSEventRecord) ([]string, string) {
	defaultPath := []string{"errored"}
	defaultFilename := time.Now().UTC().Format(time.RFC3339Nano) + ".json"

	paths := []string{}
	filename := ""

	var msg map[string]interface{}

	msg, err := parseJSON(ctx, event.SNS.Message)
	if err != nil {
		log.Println("Error parsing message as JSON")
		return defaultPath, defaultFilename
	}

	mailBody, ok := msg["mail"].(map[string]interface{})
	if !ok {
		log.Println("Error asserting mailBody to map[string]interface{}")
		return defaultPath, defaultFilename
	}

	// msg, err = parseJSON(ctx, mailBody)
	// if err != nil {
	// 	log.Println("Error parsing mailBody as JSON")
	// 	return defaultPath, defaultFilename
	// }

	timestamp, ok := mailBody["timestamp"].(string)
	if !ok {
		log.Println("Error asserting timestamp")
		return defaultPath, defaultFilename
	}

	messageID, ok := mailBody["messageId"].(string)
	if !ok {
		log.Println("Error asserting messageID")
		return defaultPath, defaultFilename
	}

	filename = timestamp + "_" + messageID[0:8] + ".json"

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
			return defaultPath, filename
		}
	} else {
		log.Println("Error asserting destination", mailBody["destination"])
		return defaultPath, filename
	}

	return paths, filename
}

func parseJSON(ctx context.Context, body string) (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
