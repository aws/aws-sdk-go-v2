package s3manager

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	// DefaultBatchSize is the batch size we initialize when constructing a batch delete client.
	// This value is used when calling DeleteObjects. This represents how many objects to delete
	// per DeleteObjects call.
	DefaultBatchSize = 100
)

// BatchError is a collection of indecent errors that occurred with a given batch.
type BatchError interface {
	error

	Errors() []error
}

// batchError will contain the key and bucket of the object that failed to
// either upload or download.
type batchError struct {
	failures errSlice
	message  string
}

func (b *batchError) Error() string {
	return fmt.Sprintf("%s: %s", b.message, b.failures.String())
}

func (b *batchError) Errors() []error {
	return b.failures
}

// errSlice is a typed alias for a slice of errors to satisfy the error
// interface.
type errSlice []error

func (errs errSlice) String() string {
	buf := bytes.NewBuffer(nil)
	for i, err := range errs {
		buf.WriteString(err.Error())
		if i+1 < len(errs) {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// BatchItemError is an error that occurred while processing a given bucket and key.
type BatchItemError interface {
	error

	Bucket() string
	Key() string
}

// batchItemError will contain the original error, bucket, and key of the operation that failed
// during batch operations.
type batchItemError struct {
	err    error
	bucket string
	key    string
}

func newBatchItemError(err error, bucket, key *string) *batchItemError {
	return &batchItemError{
		err:    err,
		bucket: aws.ToString(bucket),
		key:    aws.ToString(key),
	}
}

// Error is the error string
func (b *batchItemError) Error() string {
	origErr := ""
	if b.err != nil {
		origErr = ":\n" + b.err.Error()
	}
	return fmt.Sprintf("failed to perform batch operation on %q to %q%s", b.key, b.bucket, origErr)
}

func (b *batchItemError) Unwrap() error {
	return b.err
}

func (b *batchItemError) Bucket() string {
	return b.bucket
}

func (b *batchItemError) Key() string {
	return b.key
}

// BatchDeleteIterator is an interface that uses the scanner pattern to
// iterate through what needs to be deleted.
type BatchDeleteIterator interface {
	Next() bool
	Err() error
	DeleteObject() BatchDeleteObject
}

// DeleteListIterator is an alternative iterator for the BatchDelete client. This will
// iterate through a list of objects and delete the objects.
//
// Example:
//	iter := NewDeleteListIterator(client, &s3.ListObjectsV2Input{
//		Bucket: aws.String("bucket"),
//	})
//
//	batcher := s3manager.NewBatchDelete(svc)
//	if err := batcher.Delete(context.Background(), iter); err != nil {
//		return err
//	}
type DeleteListIterator struct {
	bucket    *string
	paginator *listObjectsV2Paginator
	objects   []*types.Object
	err       error
}

// NewDeleteListIterator will return a new DeleteListIterator.
func NewDeleteListIterator(client ListObjectsV2APIClient, input *s3.ListObjectsV2Input, opts ...func(*DeleteListIterator)) BatchDeleteIterator {
	iter := &DeleteListIterator{
		bucket:    input.Bucket,
		paginator: newListObjectsV2Paginator(client, input),
	}

	for _, opt := range opts {
		opt(iter)
	}
	return iter
}

// Next will use the S3API client to iterate through a list of objects.
func (iter *DeleteListIterator) Next() bool {
	if iter.err != nil {
		return false
	}

	if len(iter.objects) > 0 {
		iter.objects = iter.objects[1:]
	}

	if len(iter.objects) == 0 && iter.paginator.HasMorePages() {
		nextPage, err := iter.paginator.NextPage(context.TODO())
		if err != nil {
			iter.err = err
			return false
		}
		iter.objects = nextPage.Contents
	}

	return len(iter.objects) > 0
}

// Err will return the last known error from Next.
func (iter *DeleteListIterator) Err() error {
	return iter.err
}

// DeleteObject will return the current object to be deleted.
func (iter *DeleteListIterator) DeleteObject() BatchDeleteObject {
	return BatchDeleteObject{
		Object: &s3.DeleteObjectInput{
			Bucket: iter.bucket,
			Key:    iter.objects[0].Key,
		},
	}
}

// BatchDelete will use the s3 package's service client to perform a batch
// delete.
type BatchDelete struct {
	Client    DeleteObjectsAPIClient
	BatchSize int
}

// NewBatchDelete will return a new delete client that can delete a batched amount of
// objects.
//
// Example:
//	batcher := s3manager.NewBatchDelete(client)
//
//	objects := []BatchDeleteObject{
//		{
//			Object:	&s3.DeleteObjectInput {
//				Key: aws.String("key"),
//				Bucket: aws.String("bucket"),
//			},
//		},
//	}
//
//	if err := batcher.Delete(context.Background(), &s3manager.DeleteObjectsIterator{
//		Objects: objects,
//	}); err != nil {
//		return err
//	}
func NewBatchDelete(client DeleteObjectsAPIClient, options ...func(*BatchDelete)) *BatchDelete {
	svc := &BatchDelete{
		Client:    client,
		BatchSize: DefaultBatchSize,
	}

	for _, opt := range options {
		opt(svc)
	}

	return svc
}

// BatchDeleteObject is a wrapper object for calling the batch delete operation.
type BatchDeleteObject struct {
	Object *s3.DeleteObjectInput
	// After will run after each iteration during the batch process. This function will
	// be executed regardless whether or not the request was successful.
	After func() error
}

// DeleteObjectsIterator is an interface that uses the scanner pattern to iterate
// through a series of objects to be deleted.
type DeleteObjectsIterator struct {
	Objects []BatchDeleteObject
	index   int
	inc     bool
}

// Next will increment the default iterator's index and ensure that there
// is another object to iterator to.
func (iter *DeleteObjectsIterator) Next() bool {
	if iter.inc {
		iter.index++
	} else {
		iter.inc = true
	}
	return iter.index < len(iter.Objects)
}

// Err will return an error. Since this is just used to satisfy the BatchDeleteIterator interface
// this will only return nil.
func (iter *DeleteObjectsIterator) Err() error {
	return nil
}

// DeleteObject will return the BatchDeleteObject at the current batched index.
func (iter *DeleteObjectsIterator) DeleteObject() BatchDeleteObject {
	object := iter.Objects[iter.index]
	return object
}

// Delete will use the iterator to queue up objects that need to be deleted.
// Once the batch size is met, this will call the deleteBatch function.
func (d *BatchDelete) Delete(ctx context.Context, iter BatchDeleteIterator) error {
	var errs []error
	var objects []BatchDeleteObject
	var input *s3.DeleteObjectsInput

	for iter.Next() {
		o := iter.DeleteObject()

		if input == nil {
			input = initDeleteObjectsInput(o.Object)
		}

		parity := hasParity(input, o)
		if parity {
			input.Delete.Objects = append(input.Delete.Objects, &types.ObjectIdentifier{
				Key:       o.Object.Key,
				VersionId: o.Object.VersionId,
			})
			objects = append(objects, o)
		}

		if len(input.Delete.Objects) == d.BatchSize || !parity {
			if err := deleteBatch(ctx, d, input, objects); err != nil {
				errs = append(errs, err...)
			}

			objects = objects[:0]
			input = nil

			if !parity {
				objects = append(objects, o)
				input = initDeleteObjectsInput(o.Object)
				input.Delete.Objects = append(input.Delete.Objects, &types.ObjectIdentifier{
					Key:       o.Object.Key,
					VersionId: o.Object.VersionId,
				})
			}
		}
	}

	// iter.Next() could return false (above) plus populate iter.Err()
	if iter.Err() != nil {
		errs = append(errs, iter.Err())
	}

	if input != nil && len(input.Delete.Objects) > 0 {
		if err := deleteBatch(ctx, d, input, objects); err != nil {
			errs = append(errs, err...)
		}
	}

	if len(errs) > 0 {
		return &batchError{failures: errs, message: "some objects have failed to be deleted."}
	}
	return nil
}

func initDeleteObjectsInput(o *s3.DeleteObjectInput) *s3.DeleteObjectsInput {
	return &s3.DeleteObjectsInput{
		Bucket:       o.Bucket,
		MFA:          o.MFA,
		RequestPayer: o.RequestPayer,
		Delete:       &types.Delete{},
	}
}

const errDefaultDeleteBatchMessage = "failed to delete"

// deleteBatch will delete a batch of items in the objects parameters.
func deleteBatch(ctx context.Context, d *BatchDelete, input *s3.DeleteObjectsInput, objects []BatchDeleteObject) []error {
	var errs []error

	if result, err := d.Client.DeleteObjects(ctx, input); err != nil {
		for i := 0; i < len(input.Delete.Objects); i++ {
			errs = append(errs, newBatchItemError(err, input.Bucket, input.Delete.Objects[i].Key))
		}
	} else if len(result.Errors) > 0 {
		for i := 0; i < len(result.Errors); i++ {
			msg := errDefaultDeleteBatchMessage
			if result.Errors[i].Message != nil {
				msg = *result.Errors[i].Message
			}

			errs = append(errs, newBatchItemError(fmt.Errorf(msg), input.Bucket, result.Errors[i].Key))
		}
	}
	for _, object := range objects {
		if object.After == nil {
			continue
		}
		if err := object.After(); err != nil {
			errs = append(errs, newBatchItemError(err, object.Object.Bucket, object.Object.Key))
		}
	}

	return errs
}

func hasParity(o1 *s3.DeleteObjectsInput, o2 BatchDeleteObject) bool {
	if o1.Bucket != nil && o2.Object.Bucket != nil {
		if *o1.Bucket != *o2.Object.Bucket {
			return false
		}
	} else if o1.Bucket != o2.Object.Bucket {
		return false
	}

	if o1.MFA != nil && o2.Object.MFA != nil {
		if *o1.MFA != *o2.Object.MFA {
			return false
		}
	} else if o1.MFA != o2.Object.MFA {
		return false
	}

	if o1.RequestPayer != o2.Object.RequestPayer {
		return false
	}

	return true
}

// BatchDownloadIterator is an interface that uses the scanner pattern to iterate
// through a series of objects to be downloaded.
type BatchDownloadIterator interface {
	Next() bool
	Err() error
	DownloadObject() BatchDownloadObject
}

// BatchDownloadObject contains all necessary information to run a batch operation once.
type BatchDownloadObject struct {
	Object *s3.GetObjectInput
	Writer io.WriterAt
	// After will run after each iteration during the batch process. This function will
	// be executed whether or not the request was successful.
	After func() error
}

// DownloadObjectsIterator implements the BatchDownloadIterator interface and allows for batched
// download of objects.
type DownloadObjectsIterator struct {
	Objects []BatchDownloadObject
	index   int
	inc     bool
}

// Next will increment the default iterator's index and ensure that there
// is another object to iterator to.
func (b *DownloadObjectsIterator) Next() bool {
	if b.inc {
		b.index++
	} else {
		b.inc = true
	}
	return b.index < len(b.Objects)
}

// DownloadObject will return the BatchDownloadObject at the current batched index.
func (b *DownloadObjectsIterator) DownloadObject() BatchDownloadObject {
	object := b.Objects[b.index]
	return object
}

// Err will return an error. Since this is just used to satisfy the BatchDeleteIterator interface
// this will only return nil.
func (b *DownloadObjectsIterator) Err() error {
	return nil
}

// BatchUploadIterator is an interface that uses the scanner pattern to
// iterate through what needs to be uploaded.
type BatchUploadIterator interface {
	Next() bool
	Err() error
	UploadObject() BatchUploadObject
}

// UploadObjectsIterator implements the BatchUploadIterator interface and allows for batched
// upload of objects.
type UploadObjectsIterator struct {
	Objects []BatchUploadObject
	index   int
	inc     bool
}

// Next will increment the default iterator's index and ensure that there
// is another object to iterator to.
func (b *UploadObjectsIterator) Next() bool {
	if b.inc {
		b.index++
	} else {
		b.inc = true
	}
	return b.index < len(b.Objects)
}

// Err will return an error. Since this is just used to satisfy the BatchUploadIterator interface
// this will only return nil.
func (b *UploadObjectsIterator) Err() error {
	return nil
}

// UploadObject will return the BatchUploadObject at the current batched index.
func (b *UploadObjectsIterator) UploadObject() BatchUploadObject {
	object := b.Objects[b.index]
	return object
}

// BatchUploadObject contains all necessary information to run a batch operation once.
type BatchUploadObject struct {
	Object *s3.PutObjectInput
	// After will run after each iteration during the batch process. This function will
	// be executed whether or not the request was successful.
	After func() error
}
