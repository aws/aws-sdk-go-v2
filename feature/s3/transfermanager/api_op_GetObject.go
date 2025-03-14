package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
)

type errReadingBody struct {
	err error
}

func (e *errReadingBody) Error() string {
	return fmt.Sprintf("failed to read part body: %v", e.err)
}

type errInvalidRange struct {
	max int64
}

func (e *errInvalidRange) Error() string {
	return fmt.Sprintf("invalid input range, must be between 0 and %d", e.max)
}

// GetObjectInput represents a request to the GetObject() or DownloadObject() call. It contains common fields
// of s3 GetObject input
type GetObjectInput struct {
	// Bucket where the object is downloaded from
	Bucket string

	// Key of the object to get.
	Key string

	Reader *ConcurrentReader

	// To retrieve the checksum, this mode must be enabled.
	//
	// General purpose buckets - In addition, if you enable checksum mode and the
	// object is uploaded with a [checksum]and encrypted with an Key Management Service (KMS)
	// key, you must have permission to use the kms:Decrypt action to retrieve the
	// checksum.
	//
	// [checksum]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_Checksum.html
	ChecksumMode types.ChecksumMode

	// The account ID of the expected bucket owner. If the account ID that you provide
	// does not match the actual owner of the bucket, the request fails with the HTTP
	// status code 403 Forbidden (access denied).
	ExpectedBucketOwner string

	// Return the object only if its entity tag (ETag) is the same as the one
	// specified in this header; otherwise, return a 412 Precondition Failed error.
	//
	// If both of the If-Match and If-Unmodified-Since headers are present in the
	// request as follows: If-Match condition evaluates to true , and;
	// If-Unmodified-Since condition evaluates to false ; then, S3 returns 200 OK and
	// the data requested.
	//
	// For more information about conditional requests, see [RFC 7232].
	//
	// [RFC 7232]: https://tools.ietf.org/html/rfc7232
	IfMatch string

	// Return the object only if it has been modified since the specified time;
	// otherwise, return a 304 Not Modified error.
	//
	// If both of the If-None-Match and If-Modified-Since headers are present in the
	// request as follows: If-None-Match condition evaluates to false , and;
	// If-Modified-Since condition evaluates to true ; then, S3 returns 304 Not
	// Modified status code.
	//
	// For more information about conditional requests, see [RFC 7232].
	//
	// [RFC 7232]: https://tools.ietf.org/html/rfc7232
	IfModifiedSince time.Time

	// Return the object only if its entity tag (ETag) is different from the one
	// specified in this header; otherwise, return a 304 Not Modified error.
	//
	// If both of the If-None-Match and If-Modified-Since headers are present in the
	// request as follows: If-None-Match condition evaluates to false , and;
	// If-Modified-Since condition evaluates to true ; then, S3 returns 304 Not
	// Modified HTTP status code.
	//
	// For more information about conditional requests, see [RFC 7232].
	//
	// [RFC 7232]: https://tools.ietf.org/html/rfc7232
	IfNoneMatch string

	// Return the object only if it has not been modified since the specified time;
	// otherwise, return a 412 Precondition Failed error.
	//
	// If both of the If-Match and If-Unmodified-Since headers are present in the
	// request as follows: If-Match condition evaluates to true , and;
	// If-Unmodified-Since condition evaluates to false ; then, S3 returns 200 OK and
	// the data requested.
	//
	// For more information about conditional requests, see [RFC 7232].
	//
	// [RFC 7232]: https://tools.ietf.org/html/rfc7232
	IfUnmodifiedSince time.Time

	// Part number of the object being read. This is a positive integer between 1 and
	// 10,000. Effectively performs a 'ranged' GET request for the part specified.
	// Useful for downloading just a part of an object.
	PartNumber int32

	// Downloads the specified byte range of an object. For more information about the
	// HTTP Range header, see [https://www.rfc-editor.org/rfc/rfc9110.html#name-range].
	//
	// Amazon S3 doesn't support retrieving multiple ranges of data per GET request.
	//
	// [https://www.rfc-editor.org/rfc/rfc9110.html#name-range]: https://www.rfc-editor.org/rfc/rfc9110.html#name-range
	Range string

	// Confirms that the requester knows that they will be charged for the request.
	// Bucket owners need not specify this parameter in their requests. If either the
	// source or destination S3 bucket has Requester Pays enabled, the requester will
	// pay for corresponding charges to copy the object. For information about
	// downloading objects from Requester Pays buckets, see [Downloading Objects in Requester Pays Buckets]in the Amazon S3 User
	// Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Downloading Objects in Requester Pays Buckets]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
	RequestPayer types.RequestPayer

	// Sets the Cache-Control header of the response.
	ResponseCacheControl string

	// Sets the Content-Disposition header of the response.
	ResponseContentDisposition string

	// Sets the Content-Encoding header of the response.
	ResponseContentEncoding string

	// Sets the Content-Language header of the response.
	ResponseContentLanguage string

	// Sets the Content-Type header of the response.
	ResponseContentType string

	// Sets the Expires header of the response.
	ResponseExpires time.Time

	// Specifies the algorithm to use when decrypting the object (for example, AES256 ).
	//
	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3, then when you
	// GET the object, you must use the following headers:
	//
	//   - x-amz-server-side-encryption-customer-algorithm
	//
	//   - x-amz-server-side-encryption-customer-key
	//
	//   - x-amz-server-side-encryption-customer-key-MD5
	//
	// For more information about SSE-C, see [Server-Side Encryption (Using Customer-Provided Encryption Keys)] in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Server-Side Encryption (Using Customer-Provided Encryption Keys)]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html
	SSECustomerAlgorithm string

	// Specifies the customer-provided encryption key that you originally provided for
	// Amazon S3 to encrypt the data before storing it. This value is used to decrypt
	// the object when recovering it and must match the one used when storing the data.
	// The key must be appropriate for use with the algorithm specified in the
	// x-amz-server-side-encryption-customer-algorithm header.
	//
	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3, then when you
	// GET the object, you must use the following headers:
	//
	//   - x-amz-server-side-encryption-customer-algorithm
	//
	//   - x-amz-server-side-encryption-customer-key
	//
	//   - x-amz-server-side-encryption-customer-key-MD5
	//
	// For more information about SSE-C, see [Server-Side Encryption (Using Customer-Provided Encryption Keys)] in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Server-Side Encryption (Using Customer-Provided Encryption Keys)]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html
	SSECustomerKey string

	// Specifies the 128-bit MD5 digest of the customer-provided encryption key
	// according to RFC 1321. Amazon S3 uses this header for a message integrity check
	// to ensure that the encryption key was transmitted without error.
	//
	// If you encrypt an object by using server-side encryption with customer-provided
	// encryption keys (SSE-C) when you store the object in Amazon S3, then when you
	// GET the object, you must use the following headers:
	//
	//   - x-amz-server-side-encryption-customer-algorithm
	//
	//   - x-amz-server-side-encryption-customer-key
	//
	//   - x-amz-server-side-encryption-customer-key-MD5
	//
	// For more information about SSE-C, see [Server-Side Encryption (Using Customer-Provided Encryption Keys)] in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Server-Side Encryption (Using Customer-Provided Encryption Keys)]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html
	SSECustomerKeyMD5 string

	// Version ID used to reference a specific version of the object.
	//
	// By default, the GetObject operation returns the current version of an object.
	// To return a different version, use the versionId subresource.
	//
	//   - If you include a versionId in your request header, you must have the
	//   s3:GetObjectVersion permission to access a specific version of an object. The
	//   s3:GetObject permission is not required in this scenario.
	//
	//   - If you request the current version of an object without a specific versionId
	//   in the request header, only the s3:GetObject permission is required. The
	//   s3:GetObjectVersion permission is not required in this scenario.
	//
	//   - Directory buckets - S3 Versioning isn't enabled and supported for directory
	//   buckets. For this API operation, only the null value of the version ID is
	//   supported by directory buckets. You can only specify null to the versionId
	//   query parameter in the request.
	//
	// For more information about versioning, see [PutBucketVersioning].
	//
	// [PutBucketVersioning]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketVersioning.html
	VersionID string
}

func (i GetObjectInput) mapGetObjectInput(enableChecksumValidation bool) *s3.GetObjectInput {
	input := &s3.GetObjectInput{
		Bucket: aws.String(i.Bucket),
		Key:    aws.String(i.Key),
	}

	if i.ChecksumMode != "" {
		input.ChecksumMode = s3types.ChecksumMode(i.ChecksumMode)
	} else if enableChecksumValidation {
		input.ChecksumMode = s3types.ChecksumModeEnabled
	}

	if i.RequestPayer != "" {
		input.RequestPayer = s3types.RequestPayer(i.RequestPayer)
	}

	input.ExpectedBucketOwner = nzstring(i.ExpectedBucketOwner)
	input.IfMatch = nzstring(i.IfMatch)
	input.IfNoneMatch = nzstring(i.IfNoneMatch)
	input.IfModifiedSince = nztime(i.IfModifiedSince)
	input.IfUnmodifiedSince = nztime(i.IfUnmodifiedSince)
	input.ResponseCacheControl = nzstring(i.ResponseCacheControl)
	input.ResponseContentDisposition = nzstring(i.ResponseContentDisposition)
	input.ResponseContentEncoding = nzstring(i.ResponseContentEncoding)
	input.ResponseContentLanguage = nzstring(i.ResponseContentLanguage)
	input.ResponseContentType = nzstring(i.ResponseContentType)
	input.ResponseExpires = nztime(i.ResponseExpires)
	input.SSECustomerAlgorithm = nzstring(i.SSECustomerAlgorithm)
	input.SSECustomerKey = nzstring(i.SSECustomerKey)
	input.SSECustomerKeyMD5 = nzstring(i.SSECustomerKeyMD5)
	input.VersionId = nzstring(i.VersionID)

	return input
}

// GetObjectOutput represents a response from GetObject() or DownloadObject() call. It contains common fields
// of s3 GetObject output
type GetObjectOutput struct {
	// Indicates that a range of bytes was specified in the request.
	AcceptRanges string

	// Object data.
	Body io.ReadCloser

	// Indicates whether the object uses an S3 Bucket Key for server-side encryption
	// with Key Management Service (KMS) keys (SSE-KMS).
	BucketKeyEnabled bool

	// Specifies caching behavior along the request/reply chain.
	CacheControl string

	// Specifies if the response checksum validation is enabled
	ChecksumMode types.ChecksumMode

	// The base64-encoded, 32-bit CRC-32 checksum of the object. This will only be
	// present if it was uploaded with the object. For more information, see [Checking object integrity]in the
	// Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32 string

	// The base64-encoded, 32-bit CRC-32C checksum of the object. This will only be
	// present if it was uploaded with the object. For more information, see [Checking object integrity]in the
	// Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32C string

	// The base64-encoded, 160-bit SHA-1 digest of the object. This will only be
	// present if it was uploaded with the object. For more information, see [Checking object integrity]in the
	// Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA1 string

	// The base64-encoded, 256-bit SHA-256 digest of the object. This will only be
	// present if it was uploaded with the object. For more information, see [Checking object integrity]in the
	// Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA256 string

	// Specifies presentational information for the object.
	ContentDisposition string

	// Indicates what content encodings have been applied to the object and thus what
	// decoding mechanisms must be applied to obtain the media-type referenced by the
	// Content-Type header field.
	ContentEncoding string

	// The language the content is in.
	ContentLanguage string

	// Size of the body in bytes.
	ContentLength int64

	// The portion of the object returned in the response.
	ContentRange string

	// A standard MIME type describing the format of the object data.
	ContentType string

	// Indicates whether the object retrieved was (true) or was not (false) a Delete
	// Marker. If false, this response header does not appear in the response.
	//
	//   - If the current version of the object is a delete marker, Amazon S3 behaves
	//   as if the object was deleted and includes x-amz-delete-marker: true in the
	//   response.
	//
	//   - If the specified version in the request is a delete marker, the response
	//   returns a 405 Method Not Allowed error and the Last-Modified: timestamp
	//   response header.
	DeleteMarker bool

	// An entity tag (ETag) is an opaque identifier assigned by a web server to a
	// specific version of a resource found at a URL.
	ETag string

	// If the object expiration is configured (see [PutBucketLifecycleConfiguration]PutBucketLifecycleConfiguration ),
	// the response includes this header. It includes the expiry-date and rule-id
	// key-value pairs providing object expiration information. The value of the
	// rule-id is URL-encoded.
	//
	// This functionality is not supported for directory buckets.
	//
	// [PutBucketLifecycleConfiguration]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketLifecycleConfiguration.html
	Expiration string

	// The date and time at which the object is no longer cacheable.
	//
	// Deprecated: This field is handled inconsistently across AWS SDKs. Prefer using
	// the ExpiresString field which contains the unparsed value from the service
	// response.
	Expires time.Time

	// The unparsed value of the Expires field from the service response. Prefer use
	// of this value over the normal Expires response field where possible.
	ExpiresString string

	// Date and time when the object was last modified.
	//
	// General purpose buckets - When you specify a versionId of the object in your
	// request, if the specified version in the request is a delete marker, the
	// response returns a 405 Method Not Allowed error and the Last-Modified: timestamp
	// response header.
	LastModified time.Time

	// A map of metadata to store with the object in S3.
	//
	// Map keys will be normalized to lower-case.
	Metadata map[string]string

	// This is set to the number of metadata entries not returned in the headers that
	// are prefixed with x-amz-meta- . This can happen if you create metadata using an
	// API like SOAP that supports more flexible metadata than the REST API. For
	// example, using SOAP, you can create metadata whose values are not legal HTTP
	// headers.
	//
	// This functionality is not supported for directory buckets.
	MissingMeta int32

	// Indicates whether this object has an active legal hold. This field is only
	// returned if you have permission to view an object's legal hold status.
	//
	// This functionality is not supported for directory buckets.
	ObjectLockLegalHoldStatus types.ObjectLockLegalHoldStatus

	// The Object Lock mode that's currently in place for this object.
	//
	// This functionality is not supported for directory buckets.
	ObjectLockMode types.ObjectLockMode

	// The date and time when this object's Object Lock will expire.
	//
	// This functionality is not supported for directory buckets.
	ObjectLockRetainUntilDate time.Time

	// The count of parts this object has. This value is only returned if you specify
	// partNumber in your request and the object was uploaded as a multipart upload.
	PartsCount int32

	// Amazon S3 can return this if your request involves a bucket that is either a
	// source or destination in a replication rule.
	//
	// This functionality is not supported for directory buckets.
	ReplicationStatus types.ReplicationStatus

	// If present, indicates that the requester was successfully charged for the
	// request.
	//
	// This functionality is not supported for directory buckets.
	RequestCharged types.RequestCharged

	// Provides information about object restoration action and expiration time of the
	// restored object copy.
	//
	// This functionality is not supported for directory buckets. Only the S3 Express
	// One Zone storage class is supported by directory buckets to store objects.
	Restore string

	// If server-side encryption with a customer-provided encryption key was
	// requested, the response will include this header to confirm the encryption
	// algorithm that's used.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerAlgorithm string

	// If server-side encryption with a customer-provided encryption key was
	// requested, the response will include this header to provide the round-trip
	// message integrity verification of the customer-provided encryption key.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKeyMD5 string

	// If present, indicates the ID of the KMS key that was used for object encryption.
	SSEKMSKeyID string

	// The server-side encryption algorithm used when you store this object in Amazon
	// S3.
	ServerSideEncryption types.ServerSideEncryption

	// Provides storage class information of the object. Amazon S3 returns this header
	// for all objects except for S3 Standard storage class objects.
	//
	// Directory buckets - Only the S3 Express One Zone storage class is supported by
	// directory buckets to store objects.
	StorageClass types.StorageClass

	// The number of tags, if any, on the object, when you have the relevant
	// permission to read object tags.
	//
	// You can use [GetObjectTagging] to retrieve the tag set associated with an object.
	//
	// This functionality is not supported for directory buckets.
	//
	// [GetObjectTagging]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectTagging.html
	TagCount int32

	// Version ID of the object.
	//
	// This functionality is not supported for directory buckets.
	VersionID string

	// If the bucket is configured as a website, redirects requests for this object to
	// another object in the same bucket or to an external URL. Amazon S3 stores the
	// value of this header in the object metadata.
	//
	// This functionality is not supported for directory buckets.
	WebsiteRedirectLocation string

	// Metadata pertaining to the operation's result.
	ResultMetadata smithymiddleware.Metadata
}

func (o *GetObjectOutput) mapFromGetObjectOutput(out *s3.GetObjectOutput, checksumMode s3types.ChecksumMode) {
	o.AcceptRanges = aws.ToString(out.AcceptRanges)
	o.CacheControl = aws.ToString(out.CacheControl)
	o.ChecksumMode = types.ChecksumMode(checksumMode)
	o.ChecksumCRC32 = aws.ToString(out.ChecksumCRC32)
	o.ChecksumCRC32C = aws.ToString(out.ChecksumCRC32C)
	o.ChecksumSHA1 = aws.ToString(out.ChecksumSHA1)
	o.ChecksumSHA256 = aws.ToString(out.ChecksumSHA256)
	o.ContentDisposition = aws.ToString(out.ContentDisposition)
	o.ContentEncoding = aws.ToString(out.ContentEncoding)
	o.ContentLanguage = aws.ToString(out.ContentLanguage)
	o.ContentRange = aws.ToString(out.ContentRange)
	o.ContentType = aws.ToString(out.ContentType)
	o.ETag = aws.ToString(out.ETag)
	o.Expiration = aws.ToString(out.Expiration)
	o.ExpiresString = aws.ToString(out.ExpiresString)
	o.Restore = aws.ToString(out.Restore)
	o.SSECustomerAlgorithm = aws.ToString(out.SSECustomerAlgorithm)
	o.SSECustomerKeyMD5 = aws.ToString(out.SSECustomerKeyMD5)
	o.SSEKMSKeyID = aws.ToString(out.SSEKMSKeyId)
	o.VersionID = aws.ToString(out.VersionId)
	o.WebsiteRedirectLocation = aws.ToString(out.WebsiteRedirectLocation)
	o.BucketKeyEnabled = aws.ToBool(out.BucketKeyEnabled)
	o.DeleteMarker = aws.ToBool(out.DeleteMarker)
	o.MissingMeta = aws.ToInt32(out.MissingMeta)
	o.PartsCount = aws.ToInt32(out.PartsCount)
	o.TagCount = aws.ToInt32(out.TagCount)
	o.ContentLength = aws.ToInt64(out.ContentLength)
	o.Body = out.Body
	o.Expires = aws.ToTime(out.Expires)
	o.LastModified = aws.ToTime(out.LastModified)
	o.ObjectLockRetainUntilDate = aws.ToTime(out.ObjectLockRetainUntilDate)
	o.Metadata = out.Metadata
	o.ObjectLockLegalHoldStatus = types.ObjectLockLegalHoldStatus(out.ObjectLockLegalHoldStatus)
	o.ObjectLockMode = types.ObjectLockMode(out.ObjectLockMode)
	o.ReplicationStatus = types.ReplicationStatus(out.ReplicationStatus)
	o.RequestCharged = types.RequestCharged(out.RequestCharged)
	o.ServerSideEncryption = types.ServerSideEncryption(out.ServerSideEncryption)
	o.StorageClass = types.StorageClass(out.StorageClass)
	o.ResultMetadata = out.ResultMetadata.Clone()
}

// GetObject downloads an object from S3, intelligently splitting large
// files into smaller parts/ranges according to config and getting them in parallel across
// multiple goroutines. You can configure the download type, chunk size and concurrency
// through the Options parameters.
//
// Additional functional options can be provided to configure the individual
// download. These options are copies of the original Options instance, the client of which GetObject is called from.
// Modifying the options will not impact the original Client and Options instance.
//
// Before calling GetObject to download object, you must create a ConcurrentReader and use that reader to
// copy response content to your final destination file or buffer. This new reader type implements io.Reader to
// concurrently download parts of large object while limiting the max local cache size during download to prevent
// too much memory space consumption when getting large objects up to multi-gigabytes. You could configure that buffer
// size by changing Options.GetBufferSize.
//
// Example of creating ConcurrentReader to call GetObject:
//
//	file, err := os.Create("your filename")
//	if err != nil {
//		log.Fatal("error when creating local file: ", err)
//	}
//	r := transfermanager.NewConcurrentReader()
//	var wg sync.WaitGroup
//	wg.Add(1)
//
// // You must read from the r in a separate goroutine to drive getter to get all parts.
//
//	go func() {
//		defer wg.Done()
//		_, err := io.Copy(file, r)
//		if err != nil {
//			log.Fatal("error when writing to local file: ", err)
//		}
//	}()
//
//	out, err := svc.GetObject(context.Background(), &transfermanager.GetObjectInput{
//		Bucket: "your-bucket",
//		Key:    "your-key",
//		Reader: r,
//	})
//
//	// must wait for r.Read() to finish
//	wg.Wait()
//	if err != nil {
//		log.Fatal("error when downloading file: ", err)
//	}
func (c *Client) GetObject(ctx context.Context, input *GetObjectInput, opts ...func(*Options)) (*GetObjectOutput, error) {
	i := getter{in: input, options: c.options.Copy(), r: input.Reader}
	for _, opt := range opts {
		opt(&i.options)
	}

	return i.get(ctx)
}

type getter struct {
	options Options
	in      *GetObjectInput
	out     *GetObjectOutput
	w       *types.WriteAtBuffer
	r       *ConcurrentReader

	wg sync.WaitGroup
	m  sync.Mutex

	offset     int64
	pos        int64
	totalBytes int64
	written    int64

	err error
}

func (g *getter) get(ctx context.Context) (out *GetObjectOutput, err error) {
	if err := g.init(ctx); err != nil {
		return nil, fmt.Errorf("unable to initialize download: %w", err)
	}

	clientOptions := []func(*s3.Options){
		func(o *s3.Options) {
			o.APIOptions = append(o.APIOptions,
				middleware.AddSDKAgentKey(middleware.FeatureMetadata, userAgentKey),
				addFeatureUserAgent,
			)
		}}

	defer close(g.r.ch)
	if g.in.PartNumber > 0 {
		return g.singleDownload(ctx, clientOptions...)
	}

	if g.options.GetObjectType == types.GetObjectParts {
		if g.in.Range != "" {
			return g.singleDownload(ctx, clientOptions...)
		}
		// must know the part size before creating stream reader
		out, err := g.options.S3.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket:     aws.String(g.in.Bucket),
			Key:        aws.String(g.in.Key),
			PartNumber: aws.Int32(1),
		}, clientOptions...)
		if err != nil {
			g.r.setErr(err)
			return nil, err
		}

		partsCount := int32(math.Max(float64(aws.ToInt32(out.PartsCount)), 1))
		partSize := int64(math.Max(float64(aws.ToInt64(out.ContentLength)), 1))
		sectionParts := int32(math.Max(1, float64(g.options.GetBufferSize/partSize)))
		capacity := sectionParts
		g.r.setPartSize(partSize)
		g.r.setCapacity(int32(math.Min(float64(capacity), float64(partsCount))))
		g.r.setPartsCount(partsCount)

		ch := make(chan getChunk, g.options.Concurrency)
		for i := 0; i < g.options.Concurrency; i++ {
			g.wg.Add(1)
			go g.downloadPart(ctx, ch, clientOptions...)
		}

		var i int32
		for i < partsCount {
			if g.getErr() != nil {
				break
			}

			if g.r.getRead() == capacity {
				capacity = int32(math.Min(float64(capacity+sectionParts), float64(partsCount)))
				g.r.setCapacity(capacity)
			}

			if i == capacity {
				continue
			}

			ch <- getChunk{start: g.pos - g.offset, part: i + 1, index: i}
			i++
			g.pos += partSize
		}

		close(ch)
		g.wg.Wait()
	} else {
		out, err := g.options.S3.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(g.in.Bucket),
			Key:    aws.String(g.in.Key),
		}, clientOptions...)
		if err != nil {
			g.r.setErr(err)
			return nil, err
		}
		if aws.ToInt64(out.ContentLength) == 0 {
			return g.singleDownload(ctx, clientOptions...)
		}
		g.totalBytes = aws.ToInt64(out.ContentLength)
		if g.in.Range != "" {
			start, totalBytes := g.getDownloadRange()
			if start < 0 || start >= g.totalBytes || totalBytes > g.totalBytes || start >= totalBytes {
				err := &errInvalidRange{
					max: g.totalBytes - 1,
				}
				g.r.setErr(err)
				return nil, err
			}
			g.pos = start
			g.totalBytes = totalBytes
			g.offset = start
		}
		total := g.totalBytes - g.offset

		partsCount := int32((total-1)/g.options.PartSizeBytes + 1)
		sectionParts := int32(math.Max(1, float64(g.options.GetBufferSize/g.options.PartSizeBytes)))
		capacity := sectionParts
		g.r.setPartSize(g.options.PartSizeBytes)
		g.r.setCapacity(int32(math.Min(float64(capacity), float64(partsCount))))
		g.r.setPartsCount(partsCount)

		ch := make(chan getChunk, g.options.Concurrency)
		for i := 0; i < g.options.Concurrency; i++ {
			g.wg.Add(1)
			go g.downloadPart(ctx, ch, clientOptions...)
		}

		var i int32
		for i < partsCount {
			if g.getErr() != nil {
				break
			}

			if g.r.getRead() == capacity {
				capacity = int32(math.Min(float64(capacity+sectionParts), float64(partsCount)))
				g.r.setCapacity(capacity)
			}

			if i == capacity {
				continue
			}

			ch <- getChunk{start: g.pos - g.offset, withRange: g.byteRange(), index: i}
			i++
			g.pos += g.options.PartSizeBytes
		}

		// Wait for completion
		close(ch)
		g.wg.Wait()
	}

	if g.err != nil {
		g.r.setErr(g.err)
		return nil, g.err
	}

	g.out.ContentLength = g.written
	g.out.ContentRange = fmt.Sprintf("bytes=%d-%d", g.offset, g.offset+g.written-1)
	return g.out, nil
}

func (g *getter) init(ctx context.Context) error {
	if g.options.PartSizeBytes < minPartSizeBytes {
		return fmt.Errorf("part size must be at least %d bytes", minPartSizeBytes)
	}
	if g.r == nil {
		return fmt.Errorf("concurrent reader is required in input")
	}

	g.r.ch = make(chan outChunk, g.options.Concurrency)
	g.totalBytes = -1

	return nil
}

func (g *getter) singleDownload(ctx context.Context, clientOptions ...func(*s3.Options)) (*GetObjectOutput, error) {
	params := g.in.mapGetObjectInput(!g.options.DisableChecksumValidation)
	out, err := g.options.S3.GetObject(ctx, params, clientOptions...)
	if err != nil {
		g.r.setErr(err)
		return nil, err
	}

	defer out.Body.Close()
	buf, err := io.ReadAll(out.Body)
	if err != nil {
		g.r.setErr(err)
		return nil, err
	}

	g.r.setPartSize(int64(math.Max(1, float64(aws.ToInt64(out.ContentLength)))))
	g.r.setCapacity(1)
	g.r.setPartsCount(1)
	g.r.ch <- outChunk{body: bytes.NewReader(buf), length: aws.ToInt64(out.ContentLength)}
	var output GetObjectOutput
	output.mapFromGetObjectOutput(out, params.ChecksumMode)
	return &output, nil
}

func (g *getter) downloadPart(ctx context.Context, ch chan getChunk, clientOptions ...func(*s3.Options)) {
	defer g.wg.Done()
	for {
		chunk, ok := <-ch
		if !ok {
			break
		}
		if g.getErr() != nil {
			continue
		}
		out, err := g.downloadChunk(ctx, chunk, clientOptions...)
		if err != nil {
			g.setErr(err)
		} else {
			g.setOutput(out)
		}
	}
}

// downloadChunk downloads the chunk from s3
func (g *getter) downloadChunk(ctx context.Context, chunk getChunk, clientOptions ...func(*s3.Options)) (*GetObjectOutput, error) {
	params := g.in.mapGetObjectInput(!g.options.DisableChecksumValidation)
	if chunk.part != 0 {
		params.PartNumber = aws.Int32(chunk.part)
	}
	if chunk.withRange != "" {
		params.Range = aws.String(chunk.withRange)
	}

	out, err := g.options.S3.GetObject(ctx, params, clientOptions...)
	if err != nil {
		return nil, err
	}

	defer out.Body.Close()
	buf, err := io.ReadAll(out.Body)
	g.incrWritten(int64(len(buf)))

	if err != nil {
		return nil, err
	}
	g.r.ch <- outChunk{body: bytes.NewReader(buf), index: chunk.index, length: aws.ToInt64(out.ContentLength)}

	output := &GetObjectOutput{}
	output.mapFromGetObjectOutput(out, params.ChecksumMode)
	return output, err
}

func (g *getter) setOutput(resp *GetObjectOutput) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.out != nil {
		return
	}
	g.out = resp
}

func (g *getter) incrWritten(n int64) {
	g.m.Lock()
	defer g.m.Unlock()

	g.written += n
}

func (g *getter) getDownloadRange() (int64, int64) {
	parts := strings.Split(strings.Split(g.in.Range, "=")[1], "-")

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		g.err = err
		return 0, 0
	}

	end, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		g.err = err
		return 0, 0
	}

	return start, end + 1
}

// byteRange returns an HTTP Byte-Range header value that should be used by the
// client to request a chunk range.
func (g *getter) byteRange() string {
	return fmt.Sprintf("bytes=%d-%d", g.pos, int64(math.Min(float64(g.totalBytes-1), float64(g.pos+g.options.PartSizeBytes-1))))
}

func (g *getter) getErr() error {
	g.m.Lock()
	defer g.m.Unlock()

	return g.err
}

func (g *getter) setErr(e error) {
	g.m.Lock()
	defer g.m.Unlock()

	g.err = e
}

type getChunk struct {
	start int64
	cur   int64

	part      int32
	withRange string

	index int32
}

type outChunk struct {
	body  io.Reader
	index int32

	length int64
	cur    int64
}

// ConcurrentReader receives object parts from working goroutines, composes those chunks in order and read
// to user's buffer. ConcurrentReader limits the max number of chunks it could receive and read at the same
// time so getter won't send following parts' request to s3 until user reads all current chunks, which avoids
// too much memory consumption when caching large object parts
type ConcurrentReader struct {
	ch  chan outChunk
	buf map[int32]*outChunk

	partsCount int32
	capacity   int32
	count      int32
	read       int32

	written  int64
	partSize int64

	m sync.Mutex

	err error
}

// NewConcurrentReader returns a ConcurrentReader used in GetObject input
func NewConcurrentReader() *ConcurrentReader {
	return &ConcurrentReader{
		buf:      make(map[int32]*outChunk),
		partSize: 1, // just a placeholder value
	}
}

// Read implements io.Reader to compose object parts in order and read to p.
// It will receive up to r.capacity chunks, read them to p if any chunk index
// fits into p scope, otherwise it will buffer those chunks and read them in
// following calls
func (r *ConcurrentReader) Read(p []byte) (int, error) {
	if cap(p) == 0 {
		return 0, nil
	}

	var written int

	for r.count < r.getCapacity() {
		if e := r.getErr(); e != nil && e != io.EOF {
			r.written += int64(written)
			r.clean()
			return written, r.getErr()
		}
		if written >= cap(p) {
			r.written += int64(written)
			return written, r.getErr()
		}

		oc, ok := <-r.ch
		if !ok {
			r.written += int64(written)
			return written, r.getErr()
		}

		r.count++
		index := r.getPartSize()*int64(oc.index) - r.written

		if index < int64(cap(p)) {
			n, err := oc.body.Read(p[index:])
			oc.cur += int64(n)
			written += n
			if err != nil && err != io.EOF {
				r.setErr(err)
				r.clean()
				r.written += int64(written)
				return written, r.getErr()
			}
		}
		if oc.cur < oc.length {
			r.buf[oc.index] = &oc
		} else {
			r.incrRead(1)
			if r.getRead() >= r.partsCount {
				r.setErr(io.EOF)
			}
		}
	}

	partSize := r.getPartSize()
	minIndex := int32(r.written / partSize)
	maxIndex := int32(math.Min(float64(((r.written + int64(cap(p)) - 1) / partSize)), float64(r.getCapacity()-1)))
	for i := minIndex; i <= maxIndex; i++ {
		if e := r.getErr(); e != nil && e != io.EOF {
			r.written += int64(written)
			r.clean()
			return written, r.getErr()
		}

		c, ok := r.buf[i]
		if ok {
			index := int64(i)*partSize + c.cur - r.written
			n, err := c.body.Read(p[index:])
			c.cur += int64(n)
			written += n
			if err != nil && err != io.EOF {
				r.setErr(err)
				r.clean()
				r.written += int64(written)
				return written, r.getErr()
			}
			if c.cur >= c.length {
				r.incrRead(1)
				delete(r.buf, i)
				if r.getRead() >= r.partsCount {
					r.setErr(io.EOF)
				}
			}
		}
	}

	r.written += int64(written)
	return written, r.getErr()
}

func (r *ConcurrentReader) setPartSize(n int64) {
	r.m.Lock()
	defer r.m.Unlock()

	r.partSize = n
}

func (r *ConcurrentReader) getPartSize() int64 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.partSize
}

func (r *ConcurrentReader) setCapacity(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.capacity = n
}

func (r *ConcurrentReader) getCapacity() int32 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.capacity
}

func (r *ConcurrentReader) setPartsCount(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.partsCount = n
}

func (r *ConcurrentReader) incrRead(n int32) {
	r.m.Lock()
	defer r.m.Unlock()

	r.read += n
}

func (r *ConcurrentReader) getRead() int32 {
	r.m.Lock()
	defer r.m.Unlock()

	return r.read
}

func (r *ConcurrentReader) setErr(err error) {
	r.m.Lock()
	defer r.m.Unlock()

	r.err = err
}

func (r *ConcurrentReader) getErr() error {
	r.m.Lock()
	defer r.m.Unlock()

	return r.err
}

func (r *ConcurrentReader) clean() {
	for {
		_, ok := <-r.ch
		if !ok {
			break
		}
	}
}
