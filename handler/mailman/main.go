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
func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	buf, _ := json.Marshal(event)
	log.Println(string(buf))
	log.Println(event)

	// verify request is from cloudfront
	cfHeaderValue, cfHeaderOk := event.Headers[verifyHeader]
	if cfHeaderOk && cfHeaderValue == verifyValue {
		// Handle routes
		switch event.RouteKey {
		case "GET /.well-known/openid-configuration":
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
		case "GET /api/auth/jwks.json":
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
		case "GET /api/{userID}/emails":
			emails, err := email.ListEmails(ctx, s3Client, mailboxBucket, mailboxPrefix, event.PathParameters["userID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				emails,
			), nil

		case "GET /api/{userID}/email/{emailID}":
			email, err := email.GetEmailByID(ctx, s3Client, mailboxBucket, mailboxPrefix, event.PathParameters["userID"], event.PathParameters["emailID"])
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}

			return buildOKResponse(ctx, false, map[string]string{
				"Content-Type": "application/json",
			},
				email,
			), nil
		case "POST /api/auth/login":
			payload := event.Body
			if event.IsBase64Encoded {
				decoded, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					return buildErrorResponse(ctx, err), err
				}
				payload = string(decoded)
			}

			token, err := auth.Login(ctx, domain, event.Headers, payload)
			if err != nil {
				log.Println(err)
				return buildErrorResponse(ctx, err), err
			}
			buf, _ := json.Marshal(token)

			res := events.APIGatewayV2HTTPResponse{
				StatusCode: 200,
				Body:       string(buf),
			}

			buf, _ = json.Marshal(res)
			log.Println(string(buf))

			return res, nil
		default:
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Body:       "Error: 404 Not found. " + domain + " path: " + event.RawPath,
			}, nil
		}
	} else if !cfHeaderOk {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 502,
			Body:       fmt.Sprintf("Error: %s", "Unable to verify if request came from CloudFront"),
		}, nil
	} else {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 403,
			Body:       fmt.Sprintf("Headers[%s]: %s == %s", verifyHeader, cfHeaderValue, verifyValue),
		}, nil
	}
}

func buildErrorResponse(ctx context.Context, err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf("%s", err),
	}
}

func buildOKResponse(ctx context.Context, cache bool, headers map[string]string, body string) events.APIGatewayV2HTTPResponse {
	finalHeaders := headers
	if !cache {
		finalHeaders["Cache-Control"] = "private,s-maxage=15"
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers:    finalHeaders,
		Body:       body,
	}
}
