package aws

import (
	"net/http"
)

// A Config provides service configuration for service clients.
type Config struct {
	// The region to send requests to. This parameter is required and must
	// be configured globally or on a per-client basis unless otherwise
	// noted. A full list of regions is found in the "Regions and Endpoints"
	// document.
	//
	// @see http://docs.aws.amazon.com/general/latest/gr/rande.html
	//   AWS Regions and Endpoints
	Region string

	// The credentials object to use when signing requests. Defaults to a
	// chain of credential providers to search for credentials in environment
	// variables, shared credential file, and EC2 Instance Roles.
	Credentials CredentialsProvider

	// The resolver to use for looking up endpoints for AWS service clients
	// to use based on region.
	EndpointResolver EndpointResolver

	// The HTTP client to use when sending requests. Defaults to
	// `http.DefaultClient`.
	HTTPClient *http.Client

	// TODO document
	Handlers Handlers

	// Retryer guides how HTTP requests should be retried in case of
	// recoverable failures.
	//
	// When nil or the value does not implement the request.Retryer interface,
	// the client.DefaultRetryer will be used.
	//
	// When both Retryer and MaxRetries are non-nil, the former is used and
	// the latter ignored.
	//
	// To set the Retryer field in a type-safe manner and with chaining, use
	// the request.WithRetryer helper function:
	//
	//   cfg := request.WithRetryer(aws.NewConfig(), myRetryer)
	//
	Retryer Retryer

	// An integer value representing the logging level. The default log level
	// is zero (LogOff), which represents no logging. To enable logging set
	// to a LogLevel Value.
	LogLevel LogLevel

	// The logger writer interface to write logging messages to. Defaults to
	// standard out.
	Logger Logger

	// EnforceShouldRetryCheck is used in the AfterRetryHandler to always call
	// ShouldRetry regardless of whether or not if request.Retryable is set.
	// This will utilize ShouldRetry method of custom retryers. If EnforceShouldRetryCheck
	// is not set, then ShouldRetry will only be called if request.Retryable is nil.
	// Proper handling of the request.Retryable field is important when setting this field.
	EnforceShouldRetryCheck bool

	// Disables the computation of request and response checksums, e.g.,
	// CRC32 checksums in Amazon DynamoDB.
	DisableComputeChecksums bool

	// Set this to `true` to force the request to use path-style addressing,
	// i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client
	// will use virtual hosted bucket addressing when possible
	// (`http://BUCKET.s3.amazonaws.com/KEY`).
	//
	// @note This configuration option is specific to the Amazon S3 service.
	// @see http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
	//   Amazon S3: Virtual Hosting of Buckets
	S3ForcePathStyle bool

	// Set this to `true` to disable the SDK adding the `Expect: 100-Continue`
	// header to PUT requests over 2MB of content. 100-Continue instructs the
	// HTTP client not to send the body until the service responds with a
	// `continue` status. This is useful to prevent sending the request body
	// until after the request is authenticated, and validated.
	//
	// http://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPUT.html
	//
	// 100-Continue is only enabled for Go 1.6 and above. See `http.Transport`'s
	// `ExpectContinueTimeout` for information on adjusting the continue wait
	// timeout. https://golang.org/pkg/net/http/#Transport
	//
	// You should use this flag to disble 100-Continue if you experience issues
	// with proxies or third party S3 compatible services.
	S3Disable100Continue bool

	// Set this to `true` to enable S3 Accelerate feature. For all operations
	// compatible with S3 Accelerate will use the accelerate endpoint for
	// requests. Requests not compatible will fall back to normal S3 requests.
	//
	// The bucket must be enable for accelerate to be used with S3 client with
	// accelerate enabled. If the bucket is not enabled for accelerate an error
	// will be returned. The bucket name must be DNS compatible to also work
	// with accelerate.
	S3UseAccelerate bool

	// DisableRestProtocolURICleaning will not clean the URL path when making rest protocol requests.
	// Will default to false. This would only be used for empty directory names in s3 requests.
	//
	// Example:
	//    sess := session.Must(session.NewSession(&aws.Config{
	//         DisableRestProtocolURICleaning: aws.Bool(true),
	//    }))
	//
	//    svc := s3.New(sess)
	//    out, err := svc.GetObject(&s3.GetObjectInput {
	//    	Bucket: aws.String("bucketname"),
	//    	Key: aws.String("//foo//bar//moo"),
	//    })
	DisableRestProtocolURICleaning bool
}

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
func NewConfig() *Config {
	return &Config{}
}

// Copy will return a shallow copy of the Config object. If any additional
// configurations are provided they will be merged into the new config returned.
func (c Config) Copy() Config {
	cp := c
	cp.Handlers = cp.Handlers.Copy()

	return cp
}
