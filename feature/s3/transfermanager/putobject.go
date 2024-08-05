package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	smithyhttp "github.com/aws/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
)

// A MultiUploadFailure wraps a failed S3 multipart upload. An error returned
// will satisfy this interface when a multi part upload failed to upload all
// chucks to S3. In the case of a failure the UploadID is needed to operate on
// the chunks, if any, which were uploaded.
//
// Example:
//
//	u := manager.NewUploader(client)
//	output, err := u.upload(context.Background(), input)
//	if err != nil {
//		var multierr manager.MultiUploadFailure
//		if errors.As(err, &multierr) {
//			fmt.Printf("upload failure UploadID=%s, %s\n", multierr.UploadID(), multierr.Error())
//		} else {
//			fmt.Printf("upload failure, %s\n", err.Error())
//		}
//	}
type MultiUploadFailure interface {
	error

	// UploadID returns the upload id for the S3 multipart upload that failed.
	UploadID() string
}

// A multiUploadError wraps the upload ID of a failed s3 multipart upload.
// Composed of BaseError for code, message, and original error
//
// Should be used for an error that occurred failing a S3 multipart upload,
// and a upload ID is available. If an uploadID is not available a more relevant
type multiUploadError struct {
	err error

	// ID for multipart upload which failed.
	uploadID string
}

// batchItemError returns the string representation of the error.
//
// # See apierr.BaseError ErrorWithExtra for output format
//
// Satisfies the error interface.
func (m *multiUploadError) Error() string {
	var extra string
	if m.err != nil {
		extra = fmt.Sprintf(", cause: %s", m.err.Error())
	}
	return fmt.Sprintf("upload multipart failed, upload id: %s%s", m.uploadID, extra)
}

// Unwrap returns the underlying error that cause the upload failure
func (m *multiUploadError) Unwrap() error {
	return m.err
}

// UploadID returns the id of the S3 upload which failed.
func (m *multiUploadError) UploadID() string {
	return m.uploadID
}

// PutObjectInput represents a request to the Upload() call.
type PutObjectInput struct {
	// Bucket the object is uploaded into
	Bucket *string

	// Object key for which the PUT action was initiated
	Key *string

	// Object data
	Body io.Reader

	// Indicates the algorithm used to create the checksum for the object
	ChecksumAlgorithm types.ChecksumAlgorithm

	// Specifies the algorithm to use when encrypting the object (for example, AES256 ).
	//
	// This functionality is not supported for directory buckets.
	SSECustomerAlgorithm *string

	// Specifies the customer-provided encryption key for Amazon S3 to use in
	// encrypting data. This value is used to store the object and then it is
	// discarded; Amazon S3 does not store the encryption key. The key must be
	// appropriate for use with the algorithm specified in the
	// x-amz-server-side-encryption-customer-algorithm header.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKey *string

	// Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321.
	// Amazon S3 uses this header for a message integrity check to ensure that the
	// encryption key was transmitted without error.
	//
	// This functionality is not supported for directory buckets.
	SSECustomerKeyMD5 *string

	// If x-amz-server-side-encryption has a valid value of aws:kms or aws:kms:dsse ,
	// this header specifies the ID (Key ID, Key ARN, or Key Alias) of the Key
	// Management Service (KMS) symmetric encryption customer managed key that was used
	// for the object. If you specify x-amz-server-side-encryption:aws:kms or
	// x-amz-server-side-encryption:aws:kms:dsse , but do not provide
	// x-amz-server-side-encryption-aws-kms-key-id , Amazon S3 uses the Amazon Web
	// Services managed key ( aws/s3 ) to protect the data. If the KMS key does not
	// exist in the same account that's issuing the command, you must use the full ARN
	// and not just the ID.
	//
	// This functionality is not supported for directory buckets.
	SSEKMSKeyID *string

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

	// The account ID of the expected bucket owner. If the account ID that you provide
	// does not match the actual owner of the bucket, the request fails with the HTTP
	// status code 403 Forbidden (access denied).
	ExpectedBucketOwner *string

	// The server-side encryption algorithm that was used when you store this object
	// in Amazon S3 (for example, AES256 , aws:kms , aws:kms:dsse ).
	//
	// General purpose buckets - You have four mutually exclusive options to protect
	// data using server-side encryption in Amazon S3, depending on how you choose to
	// manage the encryption keys. Specifically, the encryption key options are Amazon
	// S3 managed keys (SSE-S3), Amazon Web Services KMS keys (SSE-KMS or DSSE-KMS),
	// and customer-provided keys (SSE-C). Amazon S3 encrypts data with server-side
	// encryption by using Amazon S3 managed keys (SSE-S3) by default. You can
	// optionally tell Amazon S3 to encrypt data at rest by using server-side
	// encryption with other key options. For more information, see [Using Server-Side Encryption]in the Amazon S3
	// User Guide.
	//
	// Directory buckets - For directory buckets, only the server-side encryption with
	// Amazon S3 managed keys (SSE-S3) ( AES256 ) value is supported.
	//
	// [Using Server-Side Encryption]: https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html
	ServerSideEncryption types.ServerSideEncryption

	// A standard MIME type describing the format of the contents. For more
	// information, see [https://www.rfc-editor.org/rfc/rfc9110.html#name-content-type].
	//
	// [https://www.rfc-editor.org/rfc/rfc9110.html#name-content-type]: https://www.rfc-editor.org/rfc/rfc9110.html#name-content-type
	ContentType *string
}

// PutObjectOutput represents a response from the Upload() call.
type PutObjectOutput struct {
	// The URL where the object was uploaded to.
	Location string

	// The ID for a multipart upload to S3. In the case of an error the error
	// can be cast to the MultiUploadFailure interface to extract the upload ID.
	// Will be empty string if multipart upload was not used, and the object
	// was uploaded as a single PutObject call.
	UploadID string

	// The list of parts that were uploaded and their checksums. Will be empty
	// if multipart upload was not used, and the object was uploaded as a
	// single PutObject call.
	CompletedParts []types.CompletedPart

	// Indicates whether the uploaded object uses an S3 Bucket Key for server-side
	// encryption with Amazon Web Services KMS (SSE-KMS).
	BucketKeyEnabled bool

	// The base64-encoded, 32-bit CRC32 checksum of the object.
	ChecksumCRC32 *string

	// The base64-encoded, 32-bit CRC32C checksum of the object.
	ChecksumCRC32C *string

	// The base64-encoded, 160-bit SHA-1 digest of the object.
	ChecksumSHA1 *string

	// The base64-encoded, 256-bit SHA-256 digest of the object.
	ChecksumSHA256 *string

	// Entity tag for the uploaded object.
	ETag *string

	// If the object expiration is configured, this will contain the expiration date
	// (expiry-date) and rule ID (rule-id). The value of rule-id is URL encoded.
	Expiration *string

	// The bucket where the newly created object is put
	Bucket *string

	// The object key of the newly created object.
	Key *string

	// If present, indicates that the requester was successfully charged for the
	// request.
	RequestCharged types.RequestCharged

	// If present, specifies the ID of the Amazon Web Services Key Management Service
	// (Amazon Web Services KMS) symmetric customer managed customer master key (CMK)
	// that was used for the object.
	SSEKMSKeyID *string

	// If you specified server-side encryption either with an Amazon S3-managed
	// encryption key or an Amazon Web Services KMS customer master key (CMK) in your
	// initiate multipart upload request, the response includes this header. It
	// confirms the encryption algorithm that Amazon S3 used to encrypt the object.
	ServerSideEncryption types.ServerSideEncryption

	// The version of the object that was uploaded. Will only be populated if
	// the S3 Bucket is versioned. If the bucket is not versioned this field
	// will not be set.
	VersionID *string
}

// Upload uploads an object to S3, intelligently buffering large
// files into smaller chunks and sending them in parallel across multiple
// goroutines. You can configure the chunk size and concurrency through the
// Options parameters.
//
// Additional functional options can be provided to configure the individual
// upload. These options are copies of the Uploader instance Upload is called from.
// Modifying the options will not impact the original Uploader instance.
//
// It is safe to call this method concurrently across goroutines.
func (c Client) Upload(ctx context.Context, input *PutObjectInput, opts ...func(*Options)) (*PutObjectOutput, error) {
	i := uploader{in: input, cfg: c, ctx: ctx}
	// Copy ClientOptions
	clientOptions := make([]func(*s3.Options), 0, len(c.options.PutClientOptions)+1)
	clientOptions = append(clientOptions, func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions,
			middleware.AddSDKAgentKey(middleware.FeatureMetadata, userAgentKey),
			addFeatureUserAgent, // yes, there are two of these
			func(s *smithymiddleware.Stack) error {
				return s.Finalize.Insert(&setS3ExpressDefaultChecksum{}, "ResolveEndpointV2", smithymiddleware.After)
			},
		)
	})
	clientOptions = append(clientOptions, c.options.PutClientOptions...)
	i.cfg.options.PutClientOptions = clientOptions
	for _, opt := range opts {
		opt(&i.cfg.options)
	}

	return i.upload()
}

type uploader struct {
	ctx context.Context
	cfg Client

	in *PutObjectInput

	//totalSize int64 // set to -1 if the size is not known
}

func (u *uploader) upload() (*PutObjectOutput, error) {
	if err := u.init(); err != nil {
		return nil, fmt.Errorf("unable to initialize upload: %w", err)
	}
	defer u.cfg.options.partPool.Close()

	r, _, cleanUp, err := u.nextReader()

	if err == io.EOF {
		return u.singleUpload(r, cleanUp)
	} else if err != nil {
		cleanUp()
		return nil, err
	}

	mu := multiUploader{
		uploader: u,
	}
	return mu.upload(r, cleanUp)
}

func (u *uploader) init() error {
	if err := validateSupportedARNType(aws.ToString(u.in.Bucket)); err != nil {
		return err
	}

	o := &u.cfg.options
	if o.S3 == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return fmt.Errorf("error while creating default s3 cfg: %q", err)
		}
		o.S3 = s3.NewFromConfig(cfg)
	}

	if o.Concurrency == 0 {
		o.Concurrency = DefaultTransferConcurrency
	}
	if o.PartSizeBytes == 0 {
		o.PartSizeBytes = DefaultPartSizeBytes
	} else if o.PartSizeBytes < DefaultPartSizeBytes {
		return fmt.Errorf("part size must be at least %d bytes", DefaultPartSizeBytes)
	}
	if o.MaxUploadParts == 0 {
		o.MaxUploadParts = DefaultMaxUploadParts
	} else if o.MaxUploadParts > DefaultMaxUploadParts {
		return fmt.Errorf("max upload parts must be at most %d bytes", DefaultMaxUploadParts)
	}
	if err := u.initSize(); err != nil {
		return err
	}
	if o.ChecksumAlgorithm == "" {
		o.ChecksumAlgorithm = types.ChecksumAlgorithmCrc32
	}

	// If PartSize was changed or partPool was never setup then we need to allocated a new pool
	// so that we return []byte slices of the correct size
	poolCap := o.Concurrency + 1
	if o.partPool == nil || o.partPool.SliceSize() != o.PartSizeBytes {
		o.partPool = newByteSlicePool(o.PartSizeBytes)
		o.partPool.ModifyCapacity(poolCap)
	} else {
		o.partPool = &returnCapacityPoolCloser{byteSlicePool: o.partPool}
		o.partPool.ModifyCapacity(poolCap)
	}

	return nil
}

// initSize tries to detect the total stream size, setting u.totalSize. If
// the size is not known, totalSize is set to -1.
func (u *uploader) initSize() error {

	switch r := u.in.Body.(type) {
	case io.Seeker:
		totalSize, err := seekerLen(r)
		if err != nil {
			return err
		}

		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		if totalSize/u.cfg.options.PartSizeBytes >= int64(u.cfg.options.MaxUploadParts) {
			// Add one to the part size to account for remainders
			// during the size calculation. e.g odd number of bytes.
			u.cfg.options.PartSizeBytes = (totalSize / int64(u.cfg.options.MaxUploadParts)) + 1
		}
	}

	return nil
}

func getLen(in *PutObjectInput) (int64, error) {
	b := &bytes.Buffer{}
	n, err := io.Copy(b, in.Body)
	if err != nil {
		return 0, err
	}

	in.Body = bytes.NewReader(b.Bytes())
	return n, nil
}

func (u *uploader) singleUpload(r io.Reader, cleanUp func()) (*PutObjectOutput, error) {
	defer cleanUp()

	var params s3.PutObjectInput
	awsutil.Copy(&params, u.in)
	params.SSEKMSKeyId = u.in.SSEKMSKeyID
	params.Body = r

	var locationRecorder recordLocationClient
	out, err := u.cfg.options.S3.PutObject(u.ctx, &params, append(u.cfg.options.PutClientOptions, locationRecorder.WrapClient())...)
	if err != nil {
		return nil, err
	}

	return &PutObjectOutput{
		Location: locationRecorder.location,

		BucketKeyEnabled:     aws.ToBool(out.BucketKeyEnabled),
		ChecksumCRC32:        out.ChecksumCRC32,
		ChecksumCRC32C:       out.ChecksumCRC32C,
		ChecksumSHA1:         out.ChecksumSHA1,
		ChecksumSHA256:       out.ChecksumSHA256,
		ETag:                 out.ETag,
		Expiration:           out.Expiration,
		Bucket:               params.Bucket,
		Key:                  params.Key,
		RequestCharged:       out.RequestCharged,
		SSEKMSKeyID:          out.SSEKMSKeyId,
		ServerSideEncryption: out.ServerSideEncryption,
		VersionID:            out.VersionId,
	}, nil
}

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type recordLocationClient struct {
	httpClient
	location string
}

func (c *recordLocationClient) WrapClient() func(*s3.Options) {
	return func(o *s3.Options) {
		c.httpClient = o.HTTPClient
		o.HTTPClient = c
	}
}

func (c *recordLocationClient) Do(r *http.Request) (resp *http.Response, err error) {
	resp, err = c.httpClient.Do(r)
	if err != nil {
		return resp, err
	}

	if resp.Request != nil && resp.Request.URL != nil {
		url := *resp.Request.URL
		url.RawQuery = ""
		c.location = url.String()
	}
	return resp, err
}

// nextReader reads the next chunk of data from input Body
func (u *uploader) nextReader() (io.Reader, int, func(), error) {
	part, err := u.cfg.options.partPool.Get(u.ctx)
	if err != nil {
		return nil, 0, func() {}, err
	}

	n, err := readFillBuf(u.in.Body, *part)

	cleanup := func() {
		u.cfg.options.partPool.Put(part)
	}
	return bytes.NewReader((*part)[0:n]), n, cleanup, err
}

func readFillBuf(r io.Reader, b []byte) (offset int, err error) {
	for offset < len(b) && err == nil {
		var n int
		n, err = r.Read(b[offset:])
		offset += n
	}
	return offset, err
}

type multiUploader struct {
	*uploader
	wg       sync.WaitGroup
	m        sync.Mutex
	err      error
	uploadID *string
	parts    completedParts
}

type ulChunk struct {
	buf     io.Reader
	partNum *int32
	cleanup func()
}

type completedParts []types.CompletedPart

func (cp completedParts) Len() int {
	return len(cp)
}

func (cp completedParts) Less(i, j int) bool {
	return aws.ToInt32(cp[i].PartNumber) < aws.ToInt32(cp[j].PartNumber)
}

func (cp completedParts) Swap(i, j int) {
	cp[i], cp[j] = cp[j], cp[i]
}

// upload will perform a multipart upload using the firstBuf buffer containing
// the first chunk of data.
func (u *multiUploader) upload(firstBuf io.Reader, cleanup func()) (*PutObjectOutput, error) {
	var params s3.CreateMultipartUploadInput
	awsutil.Copy(&params, u.uploader.in)
	params.SSEKMSKeyId = u.uploader.in.SSEKMSKeyID

	// Create a multipart
	var locationRecorder recordLocationClient
	resp, err := u.uploader.cfg.options.S3.CreateMultipartUpload(u.ctx, &params, append(u.cfg.options.PutClientOptions, locationRecorder.WrapClient())...)
	if err != nil {
		cleanup()
		return nil, err
	}
	u.uploadID = resp.UploadId

	ch := make(chan ulChunk, u.cfg.options.Concurrency)
	for i := 0; i < u.cfg.options.Concurrency; i++ {
		// launch workers
		u.wg.Add(1)
		go u.readChunk(ch)
	}

	var partNum int32 = 1
	ch <- ulChunk{buf: firstBuf, partNum: aws.Int32(partNum), cleanup: cleanup}
	for u.geterr() == nil && err == nil {
		partNum++
		var (
			data         io.Reader
			nextChunkLen int
			ok           bool
		)
		data, nextChunkLen, cleanup, err = u.nextReader()
		ok, err = u.shouldContinue(partNum, nextChunkLen, err)
		if !ok {
			cleanup()
			if err != nil {
				u.seterr(err)
			}
			break
		}

		ch <- ulChunk{buf: data, partNum: aws.Int32(partNum), cleanup: cleanup}
	}

	// close the channel, wait for workers and complete upload
	close(ch)
	u.wg.Wait()
	completeOut := u.complete()

	if err := u.geterr(); err != nil {
		return nil, &multiUploadError{
			err:      err,
			uploadID: *u.uploadID,
		}
	}

	return &PutObjectOutput{
		Location:       locationRecorder.location,
		UploadID:       *u.uploadID,
		CompletedParts: u.parts,

		BucketKeyEnabled:     aws.ToBool(completeOut.BucketKeyEnabled),
		ChecksumCRC32:        completeOut.ChecksumCRC32,
		ChecksumCRC32C:       completeOut.ChecksumCRC32C,
		ChecksumSHA1:         completeOut.ChecksumSHA1,
		ChecksumSHA256:       completeOut.ChecksumSHA256,
		ETag:                 completeOut.ETag,
		Expiration:           completeOut.Expiration,
		Bucket:               params.Bucket,
		Key:                  completeOut.Key,
		RequestCharged:       completeOut.RequestCharged,
		SSEKMSKeyID:          completeOut.SSEKMSKeyId,
		ServerSideEncryption: completeOut.ServerSideEncryption,
		VersionID:            completeOut.VersionId,
	}, nil
}

func (u *multiUploader) shouldContinue(part int32, nextChunkLen int, err error) (bool, error) {
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read multipart upload data failed, %w", err)
	}

	if nextChunkLen == 0 {
		// No need to upload empty part, if file was empty to start
		// with empty single part would of been created and never
		// started multipart upload.
		return false, nil
	}

	// This upload exceeded maximum number of supported parts, error now.
	if part > u.cfg.options.MaxUploadParts || part > DefaultMaxUploadParts {
		var msg string
		if part > u.cfg.options.MaxUploadParts {
			msg = fmt.Sprintf("exceeded total allowed configured MaxUploadParts (%d). Adjust PartSize to fit in this limit",
				u.cfg.options.MaxUploadParts)
		} else {
			msg = fmt.Sprintf("exceeded total allowed S3 limit MaxUploadParts (%d). Adjust PartSize to fit in this limit",
				DefaultMaxUploadParts)
		}
		return false, fmt.Errorf(msg)
	}

	return true, err
}

// readChunk runs in worker goroutines to pull chunks off of the ch channel
// and send() them as UploadPart requests.
func (u *multiUploader) readChunk(ch chan ulChunk) {
	defer u.wg.Done()
	for {
		data, ok := <-ch

		if !ok {
			break
		}

		if u.geterr() == nil {
			if err := u.send(data); err != nil {
				u.seterr(err)
			}
		}

		data.cleanup()
	}
}

// send performs an UploadPart request and keeps track of the completed
// part information.
func (u *multiUploader) send(c ulChunk) error {
	params := &s3.UploadPartInput{
		Bucket:               u.in.Bucket,
		Key:                  u.in.Key,
		Body:                 c.buf,
		SSECustomerAlgorithm: u.in.SSECustomerAlgorithm,
		SSECustomerKey:       u.in.SSECustomerKey,
		SSECustomerKeyMD5:    u.in.SSECustomerKeyMD5,
		ExpectedBucketOwner:  u.in.ExpectedBucketOwner,
		RequestPayer:         u.in.RequestPayer,

		ChecksumAlgorithm: u.in.ChecksumAlgorithm,
		// Invalid to set any of the individual ChecksumXXX members from
		// PutObject as they are never valid for individual parts of a
		// multipart upload.

		PartNumber: c.partNum,
		UploadId:   u.uploadID,
	}

	resp, err := u.cfg.options.S3.UploadPart(u.ctx, params, u.cfg.options.PutClientOptions...)
	if err != nil {
		return err
	}

	var completed types.CompletedPart
	awsutil.Copy(&completed, resp)
	completed.PartNumber = c.partNum

	u.m.Lock()
	u.parts = append(u.parts, completed)
	u.m.Unlock()

	return nil
}

// geterr is a thread-safe getter for the error object
func (u *multiUploader) geterr() error {
	u.m.Lock()
	defer u.m.Unlock()

	return u.err
}

// seterr is a thread-safe setter for the error object
func (u *multiUploader) seterr(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.err = e
}

// fail will abort the multipart unless LeavePartsOnError is set to true.
func (u *multiUploader) fail() {
	params := &s3.AbortMultipartUploadInput{
		Bucket:   u.in.Bucket,
		Key:      u.in.Key,
		UploadId: u.uploadID,
	}
	_, err := u.cfg.options.S3.AbortMultipartUpload(u.ctx, params, u.cfg.options.PutClientOptions...)
	if err != nil {
		//logMessage(u.cfg.S3, aws.LogDebug, fmt.Sprintf("failed to abort multipart upload, %v", err))
		u.seterr(fmt.Errorf("failed to abort multipart upload (%v), triggered after multipart upload failed: %v", err, u.geterr()))
	}
}

// complete successfully completes a multipart upload and returns the response.
func (u *multiUploader) complete() *s3.CompleteMultipartUploadOutput {
	if u.geterr() != nil {
		u.fail()
		return nil
	}

	// Parts must be sorted in PartNumber order.
	sort.Sort(u.parts)

	var params s3.CompleteMultipartUploadInput
	awsutil.Copy(&params, u.in)
	params.UploadId = u.uploadID
	params.MultipartUpload = &types.CompletedMultipartUpload{Parts: u.parts}

	resp, err := u.cfg.options.S3.CompleteMultipartUpload(u.ctx, &params, u.cfg.options.PutClientOptions...)
	if err != nil {
		u.seterr(err)
		u.fail()
	}

	return resp
}

// setS3ExpressDefaultChecksum defaults to CRC32 for S3Express buckets,
// which is required when uploading to those through transfer manager.
type setS3ExpressDefaultChecksum struct{}

func (*setS3ExpressDefaultChecksum) ID() string {
	return "setS3ExpressDefaultChecksum"
}

func (*setS3ExpressDefaultChecksum) HandleFinalize(
	ctx context.Context, in smithymiddleware.FinalizeInput, next smithymiddleware.FinalizeHandler,
) (
	out smithymiddleware.FinalizeOutput, metadata smithymiddleware.Metadata, err error,
) {
	const checksumHeader = "x-amz-checksum-algorithm"

	if internalcontext.GetS3Backend(ctx) != internalcontext.S3BackendS3Express {
		return next.HandleFinalize(ctx, in)
	}

	// If this is CreateMultipartUpload we need to ensure the checksum
	// algorithm header is present. Otherwise everything is driven off the
	// context setting and we can let it flow from there.
	if middleware.GetOperationName(ctx) == "CreateMultipartUpload" {
		r, ok := in.Request.(*smithyhttp.Request)
		if !ok {
			return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
		}

		if internalcontext.GetChecksumInputAlgorithm(ctx) == "" {
			r.Header.Set(checksumHeader, "CRC32")
		}
		return next.HandleFinalize(ctx, in)
	} else if internalcontext.GetChecksumInputAlgorithm(ctx) == "" {
		ctx = internalcontext.SetChecksumInputAlgorithm(ctx, string(types.ChecksumAlgorithmCrc32))
	}

	return next.HandleFinalize(ctx, in)
}

func addFeatureUserAgent(stack *smithymiddleware.Stack) error {
	ua, err := getOrAddRequestUserAgent(stack)
	if err != nil {
		return err
	}

	ua.AddUserAgentFeature(middleware.UserAgentFeatureS3Transfer)
	return nil
}

func getOrAddRequestUserAgent(stack *smithymiddleware.Stack) (*middleware.RequestUserAgent, error) {
	id := (*middleware.RequestUserAgent)(nil).ID()
	mw, ok := stack.Build.Get(id)
	if !ok {
		mw = middleware.NewRequestUserAgent()
		if err := stack.Build.Add(mw, smithymiddleware.After); err != nil {
			return nil, err
		}
	}

	ua, ok := mw.(*middleware.RequestUserAgent)
	if !ok {
		return nil, fmt.Errorf("%T for %s middleware did not match expected type", mw, id)
	}

	return ua, nil
}
