package transfermanager

import (
	"context"
	"fmt"
	"io"
	"time"

	internalio "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/io"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
)

// DownloadFileInput represents a request to the DownloadFile() call. It mirrors the
// common fields of an S3 GetObject request, but instead of a caller-supplied
// io.WriterAt the object is written to a local file at FilePath.
type DownloadFileInput struct {
	// Bucket where the object is downloaded from.
	Bucket *string

	// Key of the object to get.
	Key *string

	// FilePath is the local destination path the object is written to. The file is
	// created (or truncated if it exists). Required.
	FilePath string

	// To retrieve the checksum, this mode must be enabled.
	ChecksumMode types.ChecksumMode

	// The account ID of the expected bucket owner.
	ExpectedBucketOwner *string

	// Return the object only if its entity tag (ETag) is the same as the one
	// specified; otherwise, return a 412 Precondition Failed error.
	IfMatch *string

	// Return the object only if it has been modified since the specified time;
	// otherwise, return a 304 Not Modified error.
	IfModifiedSince *time.Time

	// Return the object only if its entity tag (ETag) is different from the one
	// specified; otherwise, return a 304 Not Modified error.
	IfNoneMatch *string

	// Return the object only if it has not been modified since the specified time;
	// otherwise, return a 412 Precondition Failed error.
	IfUnmodifiedSince *time.Time

	// Downloads the specified byte range of an object. Only applies when
	// GetObjectType is GetObjectRanges.
	Range *string

	// Confirms that the requester knows that they will be charged for the request.
	RequestPayer types.RequestPayer

	// Sets the Cache-Control header of the response.
	ResponseCacheControl *string

	// Sets the Content-Disposition header of the response.
	ResponseContentDisposition *string

	// Sets the Content-Encoding header of the response.
	ResponseContentEncoding *string

	// Sets the Content-Language header of the response.
	ResponseContentLanguage *string

	// Sets the Content-Type header of the response.
	ResponseContentType *string

	// Sets the Expires header of the response.
	ResponseExpires *time.Time

	// Specifies the algorithm to use when decrypting the object (for example, AES256).
	SSECustomerAlgorithm *string

	// Specifies the customer-provided encryption key for Amazon S3 to decrypt the
	// object.
	SSECustomerKey *string

	// Specifies the 128-bit MD5 digest of the customer-provided encryption key.
	SSECustomerKeyMD5 *string

	// Version ID used to reference a specific version of the object.
	VersionID *string
}

// toDownloadObjectInput builds the DownloadObjectInput the shared downloader
// consumes, wiring the destination file as the WriterAt.
func (i *DownloadFileInput) toDownloadObjectInput(w io.WriterAt) *DownloadObjectInput {
	return &DownloadObjectInput{
		Bucket:                     i.Bucket,
		Key:                        i.Key,
		WriterAt:                   w,
		ChecksumMode:               i.ChecksumMode,
		ExpectedBucketOwner:        i.ExpectedBucketOwner,
		IfMatch:                    i.IfMatch,
		IfModifiedSince:            i.IfModifiedSince,
		IfNoneMatch:                i.IfNoneMatch,
		IfUnmodifiedSince:          i.IfUnmodifiedSince,
		Range:                      i.Range,
		RequestPayer:               i.RequestPayer,
		ResponseCacheControl:       i.ResponseCacheControl,
		ResponseContentDisposition: i.ResponseContentDisposition,
		ResponseContentEncoding:    i.ResponseContentEncoding,
		ResponseContentLanguage:    i.ResponseContentLanguage,
		ResponseContentType:        i.ResponseContentType,
		ResponseExpires:            i.ResponseExpires,
		SSECustomerAlgorithm:       i.SSECustomerAlgorithm,
		SSECustomerKey:             i.SSECustomerKey,
		SSECustomerKeyMD5:          i.SSECustomerKeyMD5,
		VersionID:                  i.VersionID,
	}
}

// DownloadFile downloads an object from S3 to a local file at input.FilePath,
// splitting it into byte ranges fetched in parallel.
//
// For write-to-disk use cases, prefer DownloadFile over DownloadObject, since
// DownloadFile has exclusive ownership of the file handle it can apply various
// optimizations based on the downloaded size and platform.
func (c *Client) DownloadFile(ctx context.Context, input *DownloadFileInput, opts ...func(*Options)) (*DownloadObjectOutput, error) {
	if input == nil || input.FilePath == "" {
		return nil, fmt.Errorf("FilePath is required")
	}

	options := c.options.Copy()
	for _, opt := range opts {
		opt(&options)
	}

	f, err := internalio.Create(input.FilePath)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	defer f.Close()

	d := downloader{in: input.toDownloadObjectInput(f), options: options}
	out, err := d.download(ctx)
	if err != nil {
		return out, fmt.Errorf("download: %w", err)
	}
	if err := f.Close(); err != nil {
		return out, fmt.Errorf("close: %w", err)
	}

	return out, nil
}
