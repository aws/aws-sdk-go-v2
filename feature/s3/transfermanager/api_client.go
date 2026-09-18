package transfermanager

import (
	"math"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const userAgentKey = "s3-transfer"

// defaultMaxUploadParts is the maximum allowed number of parts in a multi-part upload
// on Amazon S3.
const defaultMaxUploadParts = 10000

// defaultPartSizeBytes is the default part size when transferring objects to/from S3
const defaultPartSizeBytes = 1024 * 1024 * 8

// defaultMultipartUploadThreshold is the default size threshold in bytes indicating when to use multipart upload.
const defaultMultipartUploadThreshold = 1024 * 1024 * 16

// defaultTransferConcurrency is the default concurrency for clients created
// with New.
const defaultTransferConcurrency = 5

// defaultConfigTransferConcurrency is the high-performance default concurrency
// for clients created with NewFromConfig.
const defaultConfigTransferConcurrency = 128

const defaultPartBodyMaxRetries = 3

const defaultGetBufferSize = 1024 * 1024 * 50

// Client provides the API client to make operations call for Amazon Simple
// Storage Service's Transfer Manager
// It is safe to call Client methods concurrently across goroutines.
type Client struct {
	options Options
}

// NewFromConfig returns an initialized Client from the provided AWS config.
// It creates the underlying S3 client and gives SDK buildable HTTP clients an
// independently owned transport configured for high-concurrency transfers.
func NewFromConfig(cfg *aws.Config, optFns ...func(*Options)) *Client {
	cfgCopy := cfg.Copy()

	httpClient, ok := cfg.HTTPClient.(*awshttp.BuildableClient)
	if cfg.HTTPClient == nil {
		httpClient = awshttp.NewBuildableClient()
		ok = true
	}
	if ok {
		cfgCopy.HTTPClient = httpClient.WithTransportOptions(func(transport *http.Transport) {
			transport.MaxIdleConns = 0
			// Zero means http.DefaultMaxIdleConnsPerHost, not unlimited.
			transport.MaxIdleConnsPerHost = math.MaxInt
		})
	}

	options := make([]func(*Options), 0, len(optFns)+1)
	options = append(options, func(o *Options) {
		o.Concurrency = defaultConfigTransferConcurrency
	})
	options = append(options, optFns...)

	return New(s3.NewFromConfig(cfgCopy), options...)
}

// New returns an initialized Client from the client Options. Provide
// more functional options to further configure the Client
func New(s3Client S3APIClient, optFns ...func(*Options)) *Client {
	opts := Options{
		S3: s3Client,
	}
	for _, fn := range optFns {
		fn(&opts)
	}

	resolveConcurrency(&opts)
	resolvePartSizeBytes(&opts)
	resolveRequestChecksumCalculation(&opts)
	resolveMultipartUploadThreshold(&opts)
	resolveGetObjectType(&opts)
	resolvePartBodyMaxRetries(&opts)
	resolveGetBufferSize(&opts)
	resolveMaxUploadParts(&opts)

	return &Client{
		options: opts,
	}
}
