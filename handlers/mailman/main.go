package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gideonw/gopher-mail/auth"
	"github.com/gideonw/gopher-mail/email"
)

var invokeCount = 0
var s3Client *s3.Client

// ENV variables
var domain string
var mailboxBucket string
var mailboxPrefix string
var verifyHeader string
var verifyValue string

const pathPrefix = "/api"

func init() {
	// Set ENV variables on init
	domain = os.Getenv("DOMAIN")
	mailboxBucket = os.Getenv("MAILBOX_BUCKET")
	mailboxPrefix = os.Getenv("MAILBOX_PREFIX")
	verifyHeader = os.Getenv("CF_VERIFY_HEADER")
	verifyValue = os.Getenv("CF_VERIFY_VALUE")

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	// cfg.Region = endpoints.UsWest2RegionID

	s3Client = s3.New(cfg)
}

func main() {
	lambda.Start(Handler)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	buf, _ := json.Marshal(event)
	log.Println(string(buf))

	// verify request is from cloudfront
	if event.Headers[verifyHeader] != verifyValue {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil
	}

	switch event.HTTPMethod {
	case "GET":
		switch event.Resource {
		case "/.well-known/openid-configuration":
			payload, err := auth.WellKnownOpenIDConfig()
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				payload,
			), nil
		case "/auth/jwks.json":
			payload, err := auth.WellKnownJWKSJSON()
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				payload,
			), nil
		case "/api/{userID}/emails":
			emails, err := email.ListEmails(ctx, event.PathParameters["userID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				emails,
			), nil

		case "/api/{userID}/email/{emailID}":
			email, err := email.GetEmailByID(ctx, event.PathParameters["userID"], event.PathParameters["emailID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				email,
			), nil
		}
	case "POST":
		switch event.Resource {
		case "/api/auth/login":
			payload := event.Body
			if event.IsBase64Encoded {
				decoded, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					return buildErrorResponse(ctx, err), err
				}
				payload = string(decoded)
			}

			cookies, err := auth.Login(ctx, event.Headers, payload)
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			res := events.APIGatewayProxyResponse{
				StatusCode:        200,
				MultiValueHeaders: cookies,
			}

			fmt.Printf("%#v", res)

			return res, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Hello world! Go Boom From: " + domain + " path: " + event.Path,
	}, nil
}

func buildErrorResponse(ctx context.Context, err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf("%s", err),
	}
}

func buildOKResponse(ctx context.Context, cache bool, headers map[string]string, body string) events.APIGatewayProxyResponse {
	finalHeaders := headers
	if !cache {
		finalHeaders["Cache-Control"] = "private,s-maxage=15"
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    finalHeaders,
		Body:       body,
	}
}
