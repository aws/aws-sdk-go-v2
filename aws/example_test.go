package aws_test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/modeledendpoints"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func ExampleEndpointResolverFunc() {
	defaultResolver := modeledendpoints.NewDefaultResolver()
	myCustomResolver := func(service, region string) (aws.Endpoint, error) {
		if service == modeledendpoints.S3ServiceID {
			return aws.Endpoint{
				URL:           "s3.custom.endpoint.com",
				SigningRegion: "custom-signing-region",
			}, nil
		}

		return defaultResolver.ResolveEndpoint(service, region)
	}

	cfg := defaults.Config()
	cfg.Region = modeledendpoints.UsWest2RegionID
	cfg.EndpointResolver = aws.EndpointResolverFunc(myCustomResolver)

	// Create the S3 service client with the shared config. This will
	// automatically use the S3 custom endpoint configured in the custom
	// endpoint resolver wrapping the default endpoint resolver.
	s3Svc := s3.New(cfg)
	// Operation calls will be made to the custom endpoint.
	s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("myBucket"),
		Key:    aws.String("myObjectKey"),
	})

	// Create the SQS service client with the shared cfg. This will
	// fallback to the default endpoint resolver because the customization
	// passes any non S3 service endpoint resolve to the default resolver.
	sqsSvc := sqs.New(cfg)
	// Operation calls will be made to the default endpoint for SQS for the
	// region configured.
	sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("my-queue-url"),
	})
}
