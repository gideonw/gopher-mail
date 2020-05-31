package email

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/DusanKasan/parsemail"
	"github.com/aws/aws-sdk-go-v2/aws"
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

// Meta contians a snapshot of an email for the frontend
type Meta struct {
	MessageID string
	Subject   string
	Date      time.Time
}

// MetaMultiMap represents a simple json object with a single key holding the emails
type MetaMultiMap map[string][]Meta

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

	ret := make(MetaMultiMap)
	ret["emails"] = []Meta{}

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

			ret["emails"] = append(ret["emails"], Meta{
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
