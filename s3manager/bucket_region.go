package s3manager

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

const bucketRegionHeader = "X-Amz-Bucket-Region"

// GetBucketRegion will attempt to get the region for a bucket using the
// client's configured region to determine which AWS partition to perform the query on.
//
// The request will not be signed, and will not use your AWS credentials.
//
// A "bucket not found" error will be returned if the bucket does not exist in the
// AWS partition the client region belongs to.
//
// For example to get the region of a bucket which exists in "eu-central-1"
// you could provide a region hint of "us-west-2".
//
//	cfg := config.LoadDefaultConfig()
//
//	bucket := "my-bucket"
//	region, err := s3manager.GetBucketRegion(ctx, s3.NewFromConfig(cfg), bucket)
//	if err != nil {
//		var bnf BucketNotFound
//		if errors.As(err, &bnf) {
//			fmt.Fprintf(os.Stderr, "unable to find bucket %s's region\n", bucket)
//		}
//	}
//	fmt.Printf("Bucket %s is in %s region\n", bucket, region)
//
// By default the request will be made to the Amazon S3 endpoint using the Path
// style addressing.
//
//	s3.us-west-2.amazonaws.com/bucketname
//
// This is not compatible with Amazon S3's FIPS endpoints. To override this
// behavior to use Virtual Host style addressing, provide a functional option
// that will set the client options UsePathStyle to false.
//
//	region, err := s3manager.GetBucketRegion(s3.NewFromConfig(cfg), "bucketname", func(o *s3.Options) {
//		o.UsePathStyle = false
//	})
//
// To configure the GetBucketRegion to make a request via the Amazon
// S3 FIPS endpoints directly when a FIPS region name is not available, (e.g.
// fips-us-gov-west-1) set the Config.Endpoint on the Session, or client the
// utility is called with. The hint region will be ignored if an endpoint URL
// is configured on the session or client.
//
//	cfg, err := config.LoadDefaultConfig(config.WithEndpointResolver{
//		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
//			return aws.Endpoint{URL: "https://s3-fips.us-west-2.amazonaws.com"}, nil
//		}),
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	region, err := s3manager.GetBucketRegion(context.Background(), s3.NewFromConfig(cfg), "bucketname")
//	if err != nil {
//		panic(err)
//	}
func GetBucketRegion(ctx context.Context, client HeadBucketAPIClient, bucket string, optFns ...func(*s3.Options)) (string, error) {
	var captureBucketRegion deserializeBucketRegion

	clientOptionFns := make([]func(*s3.Options), 0, len(optFns)+1)
	clientOptionFns = append(clientOptionFns, func(options *s3.Options) {
		options.UsePathStyle = true
		options.Credentials = aws.AnonymousCredentials{}

		// Disable HTTP redirects to prevent an invalid 301 from eating the response
		// because Go's HTTP client will fail, and drop the response if an 301 is
		// received without a location header. S3 will return a 301 without the
		// location header for HeadObject API calls.
		// TODO: log warning if we can't configure the client for not following redirect
		if buildableHTTPClient, ok := options.HTTPClient.(*aws.BuildableHTTPClient); ok {
			options.HTTPClient = buildableHTTPClient.WithCheckRedirect(func(redirect *func(req *http.Request, via []*http.Request) error) {
				orig := *redirect
				*redirect = func(req *http.Request, via []*http.Request) error {
					err := orig(req, via)
					if err == nil {
						return http.ErrUseLastResponse
					}
					return err
				}
			})
		}

		options.APIOptions = append(options.APIOptions, captureBucketRegion.RegisterMiddleware)
	})
	copy(clientOptionFns[1:], optFns)

	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}, clientOptionFns...)
	if len(captureBucketRegion.BucketRegion) == 0 && err != nil {
		var httpStatusErr interface {
			HTTPStatusCode() int
		}
		if !errors.As(err, &httpStatusErr) {
			return "", err
		}

		if httpStatusErr.HTTPStatusCode() == http.StatusNotFound {
			return "", &bucketNotFound{}
		}

		return "", err
	}

	bucketRegion := normalizeBucketLocation(captureBucketRegion.BucketRegion)

	return bucketRegion, nil
}

type deserializeBucketRegion struct {
	BucketRegion string
}

func (d *deserializeBucketRegion) RegisterMiddleware(stack *middleware.Stack) error {
	return stack.Deserialize.Add(d, middleware.After)
}

func (d *deserializeBucketRegion) ID() string {
	return "DeserializeBucketRegion"
}

func (d *deserializeBucketRegion) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}

	resp, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", out.RawResponse)
	}

	d.BucketRegion = resp.Header.Get(bucketRegionHeader)

	return out, metadata, err
}

// normalizeBucketLocation is a utility function which will update the
// passed in value to always be a region ID. Generally this would be used
// with GetBucketLocation API operation.
//
// Replaces empty string with "us-east-1", and "EU" with "eu-west-1".
//
// See http://docs.aws.amazon.com/AmazonS3/latest/API/RESTBucketGETlocation.html
// for more information on the values that can be returned.
func normalizeBucketLocation(loc string) string {
	switch loc {
	case "":
		loc = "us-east-1"
	case "EU":
		loc = "eu-west-1"
	}

	return loc
}

// BucketNotFound indicates the bucket was not found in the partition when calling GetBucketRegion.
type BucketNotFound interface {
	error

	isBucketNotFound()
}

type bucketNotFound struct{}

func (b *bucketNotFound) Error() string {
	return "bucket not found"
}

func (b *bucketNotFound) isBucketNotFound() {}
