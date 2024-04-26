package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// ObjectReader is an io.Reader implementation for an S3 Object
type ObjectReader struct {
	ctx           context.Context
	s3            DownloadAPIClient
	rd            *io.PipeReader
	chunkSize     int64
	concurrency   int
	clientOptions []func(*s3.Options)
	input         *s3.GetObjectInput
}

// NewReader returns a new ObjectReader
func (d *Downloader) NewReader(ctx context.Context, input *s3.GetObjectInput) io.Reader {
	rd := &ObjectReader{
		ctx:           ctx,
		s3:            d.S3,
		input:         input,
		chunkSize:     d.PartSize,
		concurrency:   d.Concurrency,
		clientOptions: d.ClientOptions,
	}

	return rd
}

// Read is the io.Reader implementation for the ObjectReader.
//
// It returns an fs.ErrNotExists if the object doesn't exist in the given bucket.
// And returns an io.EOF when all bytes are read.
func (r *ObjectReader) Read(p []byte) (int, error) {
	if r.rd == nil {
		if err := r.preRead(); err != nil {
			return 0, err
		}
	}

	c, err := r.rd.Read(p)
	if err != nil && err == io.ErrClosedPipe {
		err = fs.ErrClosed
	}

	return c, err
}

func (r *ObjectReader) preRead() error {
	ctx := r.ctx
	rd, wr := io.Pipe()

	r.rd = rd

	// Get total Content length
	res, err := r.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: r.input.Bucket,
		Key:    r.input.Key,
		Range:  aws.String("bytes=0-0"),
	}, r.clientOptions...)
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				rd.CloseWithError(io.EOF)

				return fs.ErrNotExist
			default:
				return apiError
			}
		}
	}
	defer res.Body.Close()

	var contentLen int64
	if res.ContentRange == nil {
		if l := aws.ToInt64(res.ContentLength); l > 0 {
			contentLen = l
		}
	} else {
		parts := strings.Split(*res.ContentRange, "/")

		total := int64(-1)
		var err error
		// Checking for whether or not a numbered total exists
		// If one does not exist, we will assume the total to be -1, undefined,
		// and sequentially download each chunk until hitting a 416 error
		totalStr := parts[len(parts)-1]
		if totalStr != "*" {
			total, err = strconv.ParseInt(totalStr, 10, 64)
			if err != nil {
				return err
			}
		}

		contentLen = total
	}

	cl := newConcurrencyLock(r.concurrency)

	nextLock := make(chan struct{}, 1)

	go r.getChunk(ctx, wr, cl, nextLock, 0, contentLen)
	defer close(nextLock)

	return nil
}

func (r *ObjectReader) getChunk(
	ctx context.Context,
	wr *io.PipeWriter,
	cl *concurrencyLock,
	sequenceLock chan struct{},
	start, contentLen int64,
) {
	if start == contentLen+1 { // EOF
		defer cl.Close()

		select {
		case <-ctx.Done():
		case <-sequenceLock:
			wr.CloseWithError(io.EOF)
		}

		return
	}

	cl.Lock()
	defer cl.Unlock()

	end := start + int64(r.chunkSize)
	if end > contentLen {
		end = contentLen
	}

	nextLock := make(chan struct{}, 1)
	defer close(nextLock)

	go r.getChunk(ctx, wr, cl, nextLock, end+1, contentLen)

	res, err := r.getObject(ctx, start, end)
	if err != nil {
		wr.CloseWithError(err)
		return
	}

	defer res.Body.Close()

	select {
	case <-ctx.Done():
		return
	case <-sequenceLock:
		if _, err := io.Copy(wr, res.Body); err != nil && err != io.EOF {
			wr.CloseWithError(err)
		}
	}
}

func (r *ObjectReader) getObject(ctx context.Context, start, end int64) (*s3.GetObjectOutput, error) {
	byteRange := fmt.Sprintf("bytes=%d-%d", start, end)

	copyInput := *r.input
	copyInput.Range = &byteRange

	return r.s3.GetObject(ctx, &copyInput, r.clientOptions...)
}
