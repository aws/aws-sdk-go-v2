package transfermanager

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const userAgentKey = "s3-transfer"

// DefaultMaxUploadParts is the maximum allowed number of parts in a multi-part upload
// on Amazon S3.
const DefaultMaxUploadParts int32 = 10000

// DefaultPartSizeBytes is the default part size when transferring objects to/from S3
const DefaultPartSizeBytes int64 = 1024 * 1024 * 8

// DefaultMPUThreshold is the default size threshold in bytes indicating when to use multipart upload.
const DefaultMPUThreshold int64 = 1024 * 1024 * 16

// DefaultTransferConcurrency is the default number of goroutines to spin up when
// using PutObject().
const DefaultTransferConcurrency = 5

// Client provides the API client to make operations call for Amazon Simple
// Storage Service's Transfer Manager
type Client struct {
	options Options
}

// New returns an initialized Client from the client Options. Provide
// more functional options to further configure the Client
func New(opts Options, optFns ...func(*Options)) *Client {
	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		options: opts,
	}
}

// NewFromConfig returns a new Client from the provided s3 config
func NewFromConfig(cfg aws.Config, optFns ...func(*Options)) *Client {
	return New(Options{
		S3: s3.NewFromConfig(cfg),
	}, optFns...)
}

// Options provides params needed for transfer api calls
type Options struct {
	// The client to use when uploading to S3.
	S3 S3APIClient

	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultUploadPartSize value will be used.
	PartSizeBytes int64

	// the threshold bytes to decide when the file should be multi-uploaded
	MPUThreshold int64

	DisableChecksum bool

	ChecksumAlgorithm types.ChecksumAlgorithm

	MaxUploadParts int32

	// The number of goroutines to spin up in parallel per call to Upload when
	// sending parts. If this is set to zero, the DefaultUploadConcurrency value
	// will be used.
	//
	// The concurrency pool is not shared between calls to Upload.
	Concurrency int

	// List of request options that will be passed down to individual PutObject
	// operation requests
	PutClientOptions []func(*s3.Options)

	// partPool allows for the re-usage of streaming payload part buffers between upload calls
	partPool byteSlicePool
}
