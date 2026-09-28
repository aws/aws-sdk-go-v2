package transfermanager

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
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

// defaultConfigTransferConcurrency is the high-throughput default concurrency
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
// Unlike New, NewFromConfig constructs and owns the underlying S3 client and
// applies defaults intended for high-throughput transfers. When the config
// uses a buildable HTTP client, NewFromConfig gives the S3 client an
// independently owned transport with larger idle connection pools. A
// non-buildable HTTP client is preserved unchanged.
//
// NewFromConfig defaults object downloads to ranged GETs. Ranged GETs provide
// fixed byte ranges, allowing downloaded responses to be written at known
// offsets and enabling optimized direct-I/O write paths when alignment permits.
// GetObjectParts downloads do not use direct I/O because multipart part sizes
// may be unequal and their write boundaries cannot be validated in advance.
func NewFromConfig(cfg aws.Config, optFns ...func(*Options)) *Client {
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		buildable, ok := o.HTTPClient.(*awshttp.BuildableClient)
		if !ok {
			return
		}

		o.HTTPClient = buildable.WithTransportOptions(func(tr *http.Transport) {
			// match these to our default max conns per host so we always keepalive
			tr.MaxIdleConns = awshttp.DefaultHTTPTransportMaxConnsPerHost
			tr.MaxIdleConnsPerHost = awshttp.DefaultHTTPTransportMaxConnsPerHost
		})
	})

	opts := Options{
		S3:            s3Client,
		Concurrency:   defaultConfigTransferConcurrency,
		GetObjectType: types.GetObjectRanges,
	}
	for _, fn := range optFns {
		fn(&opts)
	}

	return newClient(opts)
}

// New returns an initialized Client using the caller-provided S3 client.
// This is the advanced constructor and does not apply the high-throughput
// client and transport defaults used by NewFromConfig. Provide additional
// functional options to configure the Client.
func New(s3Client S3APIClient, optFns ...func(*Options)) *Client {
	opts := Options{
		S3: s3Client,
	}
	for _, fn := range optFns {
		fn(&opts)
	}

	return newClient(opts)
}

func newClient(opts Options) *Client {
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
