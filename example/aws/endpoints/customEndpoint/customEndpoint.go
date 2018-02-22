// +build example

package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	defaultResolver := endpoints.NewDefaultResolver()
	s3CustResolverFn := func(service, region string) (aws.Endpoint, error) {
		if service == "s3" {
			return aws.Endpoint{
				URL:           "s3.custom.endpoint.com",
				SigningRegion: "custom-signing-region",
			}, nil
		}

		return defaultResolver.ResolveEndpoint(service, region)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	cfg.Region = "us-west-2"
	cfg.EndpointResolver = aws.EndpointResolverFunc(s3CustResolverFn)

	// Create the S3 service client with the shared config. This will
	// automatically use the S3 custom endpoint configured in the custom
	// endpoint resolver wrapping the default endpoint resolver.
	s3Svc := s3.New(cfg)
	// Operation calls will be made to the custom endpoint.
	getReq := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String("myObjectKey"),
	})
	getReq.Send()

	// Create the SQS service client with the shared config. This will
	// fallback to the default endpoint resolver because the customization
	// passes any non S3 service endpoint resolve to the default resolver.
	sqsSvc := sqs.New(cfg)
	// Operation calls will be made to the default endpoint for SQS for the
	// region configured.
	msgReq := sqsSvc.ReceiveMessageRequest(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("my-queue-url"),
	})
	msgReq.Send()

	// Create a DynamoDB service client that will use a custom endpoint
	// resolver that overrides the shared config's. This is useful when
	// custom endpoints are generated, or multiple endpoints are switched on
	// by a region value.
	ddbCustResolverFn := func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "dynamodb.custom.endpoint",
			SigningRegion: "custom-signing-region",
		}, nil
	}

	cfgCp := cfg.Copy()
	cfgCp.EndpointResolver = aws.EndpointResolverFunc(ddbCustResolverFn)

	ddbSvc := dynamodb.New(cfgCp)
	// Operation calls will be made to the custom endpoint set in the
	// ddCustResolverFn.
	listReq := ddbSvc.ListTablesRequest(&dynamodb.ListTablesInput{})
	listReq.Send()

	// Setting Config's Endpoint will override the EndpointResolver. Forcing
	// the service clien to make all operation to the endpoint specified
	// the in the config.
	cfgCp = cfg.Copy()
	cfgCp.EndpointResolver = aws.ResolveWithEndpointURL("http://localhost:8088")

	ddbSvcLocal := dynamodb.New(cfgCp)
	listReq = ddbSvcLocal.ListTablesRequest(&dynamodb.ListTablesInput{})
	listReq.Send()
}
