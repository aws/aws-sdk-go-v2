package transfermanager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	internalio "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/io"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
)

var (
	createDownloadFileFn = internalio.Create
	renameFileFn         = os.Rename
	removeFileFn         = os.Remove
	syncDirFn            = internalio.SyncDir
)

// DownloadFileInput represents a request to the DownloadFile() call. It mirrors the
// common fields of an S3 GetObject request, but instead of a caller-supplied
// io.WriterAt the object is written to a local file at FilePath.
type DownloadFileInput struct {
	// Bucket where the object is downloaded from.
	Bucket *string

	// Key of the object to get.
	Key *string

	// FilePath is the local destination path the object is written to. Required.
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
// The destination is replaced atomically: data is written to a temporary file
// alongside input.FilePath, synced to stable storage, and renamed into place on
// success. On failure the temporary file is removed and an existing destination
// is not modified.
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

	tmp, err := tempPath(input.FilePath)
	if err != nil {
		return nil, err
	}

	f, err := createDownloadFileFn(tmp)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	d := downloader{in: input.toDownloadObjectInput(f), options: options}
	out, err := d.download(ctx)
	if err != nil {
		_ = f.Close()
		removeTemp(tmp)
		return out, fmt.Errorf("download: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		removeTemp(tmp)
		return out, fmt.Errorf("sync: %w", err)
	}
	if err := f.Close(); err != nil {
		removeTemp(tmp)
		return out, fmt.Errorf("close: %w", err)
	}
	if err := renameFileFn(tmp, input.FilePath); err != nil {
		removeTemp(tmp)
		return out, fmt.Errorf("rename: %w", err)
	}
	// the destination now holds the complete object, but the rename itself
	// isn't durable until the directory is synced
	if err := syncDirFn(filepath.Dir(input.FilePath)); err != nil {
		return out, fmt.Errorf("sync dir: %w", err)
	}

	return out, nil
}

// The temp file has to live in the destination's directory so the final rename
// stays on one filesystem and is atomic.
func tempPath(path string) (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate temp file name: %w", err)
	}
	return fmt.Sprintf("%s.%s.tmp", path, hex.EncodeToString(b[:])), nil
}

// On Linux the file is opened lazily after the first response, so the temp file
// may not exist if the download failed early.
func removeTemp(path string) {
	_ = removeFileFn(path)
}
