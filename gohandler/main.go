package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

type Event struct {
	Username string
}

type AppContext struct {
	s3api *s3.Client
	cfg   aws.Config
}

// https://github.com/aws/aws-sdk-go-v2

func (app AppContext) handler(e Event) (string, error) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String("examplebucket"),
	}
	req := app.s3api.CreateBucketRequest(input)
	result, err := req.Send(context.Background())

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

	}

	region, err := s3manager.GetBucketRegion(context.Background(), app.cfg, "examplebucket", endpoints.UsEast2RegionID)
	log.Printf("result %s", result)
	log.Printf("result %s", region)

	return fmt.Sprintf("parsed username=%s create bucket in region=%s", e.Username, region), nil
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		exitErrorf("failed to load config, %v", err)

	}
	local_stack_host := os.Getenv("local_stack_host")
	if local_stack_host == "" {
		os.Exit(1)
	}

	s := fmt.Sprintf("http://%s:4572", local_stack_host)
	log.Printf("result %s", s)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(s)
	cfg.Region = "us-east-2"

	s3Svc := s3.New(cfg)
	s3Svc.ForcePathStyle = true
	app := AppContext{s3Svc, cfg}

	lambda.Start(app.handler)
}

func blockForever() {
	select {}
}

func exitErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
