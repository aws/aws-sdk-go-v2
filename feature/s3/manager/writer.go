package manager

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ObjectWriter is an io.WriteCloser implementation for an s3 Object
type ObjectWriter struct {
	ctx           context.Context
	s3            UploadAPIClient
	wr            *io.PipeWriter
	partSize      int64
	concurrency   int
	alc           types.ObjectCannedACL
	clientOptions []func(*s3.Options)
	input         *s3.PutObjectInput

	mux        sync.Mutex
	wg         sync.WaitGroup
	parts      []types.CompletedPart
	closingErr chan error
}

// NewWriter returns a new object writer
//
// Note you MUST close the io.WriteCloser in order to complete an upload
// The close method can return any errors that occured during the write if the
// write method hasn't returned them along with any complete upload / abort upload errors
func (u *Uploader) NewWriter(ctx context.Context, input *s3.PutObjectInput) io.WriteCloser {
	wr := &ObjectWriter{
		ctx:           ctx,
		s3:            u.S3,
		input:         input,
		partSize:      u.PartSize,
		concurrency:   u.Concurrency,
		clientOptions: u.ClientOptions,

		closingErr: make(chan error, 1),
	}

	return wr
}

// Write is the io.Writer implementation of the ObjectWriter
//
// The object is stored when the Close method is called.
func (w *ObjectWriter) Write(p []byte) (int, error) {
	if w.wr == nil {
		if err := w.preWrite(); err != nil {
			return 0, err
		}
	}

	return w.wr.Write(p)
}

// Close completes the write opperation.
//
// If the byte size is less than writer's chunk size then a simply PutObject opperation is preformed.
// Otherwise a multipart upload complete opperation is preformed.
// The error returned is the error from this store opperation.
//
// If an error occured while uploading parts this error might also be a upload part error joined with
// a AbortMultipartUpload error.
func (w *ObjectWriter) Close() error {
	w.wr.CloseWithError(io.EOF)

	err := <-w.closingErr

	return err
}

func (w *ObjectWriter) preWrite() error {
	ctx := w.ctx
	rd, wr := io.Pipe()

	w.wr = wr
	cl := newConcurrencyLock(w.concurrency)

	w.wg.Add(1)
	go w.writeChunk(ctx, rd, cl, nil, 1)

	return nil
}

func (w *ObjectWriter) writeChunk(ctx context.Context, rd *io.PipeReader, cl *concurrencyLock, uploadID *string, partNr int32) {
	defer w.wg.Done()

	select {
	case <-ctx.Done():
		cl.Close()
		return
	default:
		cl.Lock()
		defer cl.Unlock()

		by, err := io.ReadAll(io.LimitReader(rd, int64(w.partSize)))
		if err != nil {
			w.closeWithErr(ctx, err, rd, cl, uploadID)
			return
		}

		size := len(by)
		if partNr == 1 {
			if int64(size) < w.partSize { // For small uploads
				err = w.putObject(ctx, by)
				w.closeWithErr(ctx, err, rd, cl, uploadID)
				return
			}

			uploadID, err = w.createMultipartUpload(ctx)
			if err != nil {
				w.closeWithErr(ctx, err, rd, cl, uploadID)
				return
			}

		}

		if int64(len(by)) < w.partSize { // EOF
			go w.completeUpload(ctx, uploadID)
		} else {
			w.wg.Add(1)
			go w.writeChunk(ctx, rd, cl, uploadID, partNr+1)
		}

		part, err := w.uploadPart(ctx, uploadID, partNr, by)
		if err != nil {
			w.closeWithErr(ctx, err, rd, cl, uploadID)
			return
		}

		w.mux.Lock()
		defer w.mux.Unlock()

		w.parts = append(w.parts, part)
	}
}

func (w *ObjectWriter) closeWithErr(ctx context.Context, err error, rd *io.PipeReader, cl *concurrencyLock, uploadID *string) {
	defer close(w.closingErr)
	defer cl.Close()

	rd.CloseWithError(err)
	if uploadID != nil {
		err = errors.Join(err, w.abortUpload(ctx, uploadID))
	}

	w.closingErr <- err
}

func (w *ObjectWriter) putObject(ctx context.Context, by []byte) error {
	input := w.input
	input.Body = bytes.NewReader(by)

	_, err := w.s3.PutObject(ctx, input, w.clientOptions...)

	return err
}

func (w *ObjectWriter) createMultipartUpload(ctx context.Context) (*string, error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket:                    w.input.Bucket,
		Key:                       w.input.Key,
		ACL:                       w.input.ACL,
		BucketKeyEnabled:          w.input.BucketKeyEnabled,
		CacheControl:              w.input.CacheControl,
		ChecksumAlgorithm:         w.input.ChecksumAlgorithm,
		ContentDisposition:        w.input.ContentDisposition,
		ContentEncoding:           w.input.ContentEncoding,
		ContentLanguage:           w.input.ContentLanguage,
		ContentType:               w.input.ContentType,
		ExpectedBucketOwner:       w.input.ExpectedBucketOwner,
		Expires:                   w.input.Expires,
		GrantFullControl:          w.input.GrantFullControl,
		GrantRead:                 w.input.GrantRead,
		GrantReadACP:              w.input.GrantReadACP,
		GrantWriteACP:             w.input.GrantWriteACP,
		Metadata:                  w.input.Metadata,
		ObjectLockLegalHoldStatus: w.input.ObjectLockLegalHoldStatus,
		ObjectLockMode:            w.input.ObjectLockMode,
		ObjectLockRetainUntilDate: w.input.ObjectLockRetainUntilDate,
		RequestPayer:              w.input.RequestPayer,
		SSECustomerAlgorithm:      w.input.SSECustomerAlgorithm,
		SSECustomerKey:            w.input.SSECustomerKey,
		SSECustomerKeyMD5:         w.input.SSECustomerKeyMD5,
		SSEKMSEncryptionContext:   w.input.SSEKMSEncryptionContext,
		SSEKMSKeyId:               w.input.SSEKMSKeyId,
		ServerSideEncryption:      w.input.ServerSideEncryption,
		StorageClass:              w.input.StorageClass,
		Tagging:                   w.input.Tagging,
		WebsiteRedirectLocation:   w.input.WebsiteRedirectLocation,
	}

	res, err := w.s3.CreateMultipartUpload(ctx, input, w.clientOptions...)
	if err != nil {
		return nil, err
	}

	return res.UploadId, nil
}

func (w *ObjectWriter) uploadPart(ctx context.Context, uploadID *string, partNr int32, by []byte) (types.CompletedPart, error) {
	input := &s3.UploadPartInput{
		Bucket:     w.input.Bucket,
		Key:        w.input.Key,
		UploadId:   uploadID,
		PartNumber: &partNr,
		Body:       bytes.NewReader(by),
	}

	res, err := w.s3.UploadPart(ctx, input, w.clientOptions...)
	if err != nil {
		return types.CompletedPart{}, err
	}

	return types.CompletedPart{
		ChecksumCRC32:  res.ChecksumCRC32,
		ChecksumCRC32C: res.ChecksumCRC32C,
		ChecksumSHA1:   res.ChecksumSHA1,
		ChecksumSHA256: res.ChecksumSHA256,
		ETag:           res.ETag,
		PartNumber:     &partNr,
	}, nil
}

func (w *ObjectWriter) abortUpload(ctx context.Context, uploadID *string) error {
	input := &s3.AbortMultipartUploadInput{
		Bucket:              w.input.Bucket,
		Key:                 w.input.Key,
		UploadId:            uploadID,
		ExpectedBucketOwner: w.input.ExpectedBucketOwner,
		RequestPayer:        w.input.RequestPayer,
	}

	_, err := w.s3.AbortMultipartUpload(ctx, input, w.clientOptions...)

	return err
}

func (w *ObjectWriter) completeUpload(ctx context.Context, uploadID *string) {
	defer close(w.closingErr)

	w.wg.Wait()

	w.mux.Lock()
	defer w.mux.Unlock()

	parts := make([]types.CompletedPart, len(w.parts))
	copy(parts, w.parts)

	sort.Slice(parts, func(i, j int) bool {
		return *parts[i].PartNumber < *parts[j].PartNumber
	})

	input := &s3.CompleteMultipartUploadInput{
		Bucket:   w.input.Bucket,
		Key:      w.input.Key,
		UploadId: uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
		ExpectedBucketOwner:  w.input.ExpectedBucketOwner,
		RequestPayer:         w.input.RequestPayer,
		SSECustomerAlgorithm: w.input.SSECustomerAlgorithm,
		SSECustomerKey:       w.input.SSECustomerKey,
		SSECustomerKeyMD5:    w.input.SSECustomerKeyMD5,
	}

	_, err := w.s3.CompleteMultipartUpload(ctx, input, w.clientOptions...)

	w.closingErr <- err
}
