package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

func TestClient(t *testing.T) {
	log.Printf("%s, %s, %s", "tests", "r", "go")
	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	cfg.EndpointResolver = aws.ResolveWithEndpointURL("http://localhost:4572")
	cfg.Region = "us-east-2"

	input := &s3.CreateBucketInput{
		Bucket: aws.String("examplebucket"),
	}

	s3Svc := s3.New(cfg)
	s3Svc.ForcePathStyle = true

	req := s3Svc.CreateBucketRequest(input)
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

	region, err := s3manager.GetBucketRegion(context.Background(), cfg, "examplebucket", endpoints.UsEast2RegionID)

	log.Printf("result %s", result)
	log.Printf("result %s", region)

}

func localResolver() func(service, region string) (aws.Endpoint, error) {
	defaultResolver := endpoints.NewDefaultResolver()
	myCustomResolver := func(service, region string) (aws.Endpoint, error) {
		if service == s3.EndpointsID {
			return aws.Endpoint{
				URL:           "http://localhost:4572", // localstack
				SigningRegion: endpoints.UsEast2RegionID,
			}, nil

		}

		return defaultResolver.ResolveEndpoint(service, region)
	}

	return myCustomResolver
}

func older() {

	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	creds, err := cfg.Credentials.Retrieve()
	if err != nil {
		exitErrorf("failed to load plugin provider, %+v", err)
	}
	log.Printf("%s, %s, %s", creds.AccessKeyID, creds.SecretAccessKey, cfg.Region)

	cfg.EndpointResolver = aws.EndpointResolverFunc(localResolver())
	//	cfg.S3ForcePathStyle = aws.Bool(true)

	s3Svc := s3.New(cfg)

	input := &s3.CreateBucketInput{
		Bucket: aws.String("examplebucket"),
	}
	req := s3Svc.CreateBucketRequest(input)
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

	log.Printf("result %s", result)
}
