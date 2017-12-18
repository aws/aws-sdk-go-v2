package aws_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func ExampleEndpointResolverFunc() {
	defaultResolver := endpoints.NewDefaultResolver()
	myCustomResolver := func(service, region string) (aws.Endpoint, error) {
		if service == endpoints.S3ServiceID {
			return aws.Endpoint{
				URL:           "s3.custom.endpoint.com",
				SigningRegion: "custom-signing-region",
			}, nil
		}

		return defaultResolver.ResolveEndpoint(service, region)
	}

	cfg := defaults.Config()
	cfg.Region = endpoints.UsWest2RegionID
	cfg.EndpointResolver = aws.EndpointResolverFunc(myCustomResolver)

	// Create the S3 service client with the shared config. This will
	// automatically use the S3 custom endpoint configured in the custom
	// endpoint resolver wrapping the default endpoint resolver.
	s3Svc := s3.New(cfg)
	// Operation calls will be made to the custom endpoint.
	objReq := s3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String("myObjectKey"),
	})
	objResp, err := objReq.Send()
	if err != nil {
		panic("S3 Get Object error, " + err.Error())
	}
	fmt.Println("S3 Get object", objResp)

	// Create the SQS service client with the shared cfg. This will
	// fallback to the default endpoint resolver because the customization
	// passes any non S3 service endpoint resolve to the default resolver.
	sqsSvc := sqs.New(cfg)
	// Operation calls will be made to the default endpoint for SQS for the
	// region configured.
	msgReq := sqsSvc.ReceiveMessageRequest(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("my-queue-url"),
	})
	msgResp, err := msgReq.Send()
	if err != nil {
		panic("SQS Receive Message error, " + err.Error())
	}
	fmt.Println("SQS Receive Message", msgResp)
}
