package main

import "context"

func getEmailByID(ctx context.Context, ID string) (string, error) {
	// getInput := &s3.CopyObjectInput{
	// 	CopySource: aws.String(url.PathEscape(sourcePath)),

	// 	Bucket: aws.String(srcBucket),
	// 	Key:    aws.String(destKey),

	// 	ContentType: aws.String("application/octet-stream"),
	// }
	// log.Printf("Copying from \"%s\" to \"%s\"\n", sourcePath, srcBucket+"/"+destKey)

	// copyResp, err := s3Client.CopyObjectRequest(copyInput).Send(ctx)
	return "", nil
}

func listEmails(ctx context.Context, user string) ([]string, error) {

	return nil, nil
}
