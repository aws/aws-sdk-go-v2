package manager

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/window"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

// DownloadStream downloads an object in S3 and writes the payload into w
// using concurrent GET requests. The n int64 returned is the size of the object downloaded
// in bytes.
//
// Additional functional options can be provided to configure the individual
// download. These options are copies of the Downloader instance Download is
// called from. Modifying the options will not impact the original Downloader
// instance. Use the WithDownloaderClientOptions helper function to pass in request
// options that will be applied to all API operations made with this downloader.
//
// The w io.Writer can be satisfied by an any io.Writer for in memory downloads
// it does not require pre-allocation
//
// Example:
//	// pre-allocate is not required for this any writer works or hash.Writer
//	w := &bytes.Buffer{}
//	// download file into the memory
//	numBytesDownloaded, err := downloader.DownloadStream(ctx, w, &s3.GetObjectInput{
//		Bucket: aws.String(bucket),
//		Key:    aws.String(item),
//	})
//
// Specifying a Downloader.Concurrency of 1 will cause the Downloader to
// download the parts from S3 sequentially.
//
// It is safe to call this method concurrently across goroutines.
//
// If the GetObjectInput's Range value is provided that will cause the downloader
// to perform a single GetObjectInput request for that object's range. This will
// caused the part size, and concurrency configurations to be ignored.
func (d *Downloader) DownloadStream(ctx context.Context, w io.Writer, input *s3.GetObjectInput) (int64, error) {
	inner, cancel := context.WithCancel(ctx)
	defer cancel()

	var reportedSize int64 = -1 // total size of the file in bytes as reported by S3
	outChan, errChan := window.SlidingWindow(inner, d.Concurrency, func(location int) (interface{}, error) {
		start := int64(location) * d.PartSize
		old := atomic.LoadInt64(&reportedSize)
		// test start is less than the total
		if old > 0 && start > old {
			// This is after the end of the file. Indicate nil nil to close the window
			return nil, nil
		}
		// Copy the input and change the range
		in := &s3.GetObjectInput{}
		awsutil.Copy(in, input)

		in.Range = aws.String(d.byteRange(start))

		options := append(d.ClientOptions, d.retryOption)
		out, err := d.S3.GetObject(inner, in, options...)
		// Check if this is the end of the file err if it is don't loop
		if err != nil {
			if isEOF(err) {
				return nil, nil
			}
			return nil, err
		}

		// re-load incase it has changed
		old = atomic.LoadInt64(&reportedSize)
		if old < 0 {
			new, err := getTotalBytes(out)
			if err != nil {
				return nil, err
			}
			if !atomic.CompareAndSwapInt64(&reportedSize, old, new) {
				d.Logger.Logf(logging.Debug, "Failed to set size due to race condition. Continue uninterrupted")
			}
		}
		return out, nil
	})

	// total size of the file in bytes that are copied to the output writer
	var totalBytes int64 = 0

	// outChan is guaranteed to be in order and it needs to be
	// this is enforced by sliding window with a single consumer
	for {
		select {
		case <-inner.Done():
			// we were cancelled so just return
			return totalBytes, inner.Err()
		case err := <-errChan:
			return totalBytes, err
		case out, ok := <-outChan:
			if !ok {
				// channel was closed we have read everything so return
				// and stop the routines
				return totalBytes, nil
			}
			get, ok := out.(*s3.GetObjectOutput)
			if !ok {
				return totalBytes, fmt.Errorf("was the wrong object type got %T", out)
			}

			defer get.Body.Close()
			written, err := io.Copy(w, get.Body)

			if err != nil {
				return totalBytes, err
			}

			totalBytes += written

			if written != get.ContentLength {
				return totalBytes, io.ErrShortWrite
			}
		}
	}
}

func (d *Downloader) retryOption(options *s3.Options) {
	options.Retryer = retry.AddWithMaxAttempts(options.Retryer, d.PartBodyMaxRetries)
}
