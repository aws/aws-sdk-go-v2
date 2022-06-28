package manager

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/smithy-go/logging"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (d Downloader) Download2(ctx context.Context, w io.Writer, input *s3.GetObjectInput,
	options ...func(*Downloader)) (n int64, err error) {
	if err := validateSupportedARNType(aws.ToString(input.Bucket)); err != nil {
		return 0, err
	}

	impl := downloader2{w: w, in: input, cfg: d, ctx: ctx}

	// Copy ClientOptions
	clientOptions := make([]func(*s3.Options), 0, len(impl.cfg.ClientOptions)+1)
	clientOptions = append(clientOptions, func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions, middleware.AddSDKAgentKey(middleware.FeatureMetadata, "s3-transfer"))
	})
	clientOptions = append(clientOptions, impl.cfg.ClientOptions...)
	impl.cfg.ClientOptions = clientOptions

	for _, option := range options {
		option(&impl.cfg)
	}

	// Ensures we don't need nil checks later on
	impl.cfg.Logger = logging.WithContext(ctx, impl.cfg.Logger)

	impl.partBodyMaxRetries = d.PartBodyMaxRetries

	impl.totalBytes = -1
	if impl.cfg.Concurrency == 0 {
		impl.cfg.Concurrency = DefaultDownloadConcurrency
	}

	if impl.cfg.PartSize == 0 {
		impl.cfg.PartSize = DefaultDownloadPartSize
	}

	return impl.download()
}

// downloader2 is the implementation structure used internally by Downloader.
type downloader2 struct {
	ctx context.Context
	cfg Downloader

	in *s3.GetObjectInput
	w  io.Writer

	wg sync.WaitGroup
	m  sync.Mutex

	pos        int64
	totalBytes int64
	written    int64
	err        error

	partBodyMaxRetries int
}

// download performs the implementation of the object download across ranged
// GETs.
//nolint:funlen
func (d *downloader2) download() (n int64, err error) {
	// If PartSize was changed or partPool was never setup then we need to allocated a new pool
	// so that we return []byte slices of the correct size
	poolCap := d.cfg.Concurrency + 1
	if d.cfg.partPool == nil || d.cfg.partPool.SliceSize() != d.cfg.PartSize {
		d.cfg.partPool = newByteSlicePool(d.cfg.PartSize)
		d.cfg.partPool.ModifyCapacity(poolCap)
	} else {
		d.cfg.partPool = &returnCapacityPoolCloser{byteSlicePool: d.cfg.partPool}
		d.cfg.partPool.ModifyCapacity(poolCap)
	}

	seq := 0

	var g sync.WaitGroup
	completedCh := make(chan *dlchunk2, d.cfg.Concurrency)
	g.Add(1)
	go func() {
		defer g.Done()
		if err := receiveChunks(d.w, poolCap, completedCh); err != nil {
			d.setErr(err)
		}
	}()

	// If range is specified fall back to single download of that range
	// this enables the functionality of ranged gets with the downloader but
	// at the cost of no multipart downloads.
	if rng := aws.ToString(d.in.Range); len(rng) > 0 {
		d.downloadRange(seq, completedCh, rng)

		close(completedCh)
		g.Wait()

		return d.written, d.err
	}

	// Spin off first worker to check additional header information
	d.getChunk(seq, completedCh)
	seq++

	if total := d.getTotalBytes(); total >= 0 {
		// Spin up workers
		ch := make(chan *dlchunk2, d.cfg.Concurrency)

		for i := 0; i < d.cfg.Concurrency; i++ {
			d.wg.Add(1)
			go d.downloadPart(ch, completedCh)
		}

		// Assign work
		for d.getErr() == nil {
			if d.pos >= total {
				break // We're finished queuing chunks
			}

			// Queue the next range of bytes to read.
			chunk, err := d.nextPart(seq, d.pos)
			seq++
			if err != nil {
				break
			}

			ch <- chunk
			d.pos += d.cfg.PartSize
		}

		// Wait for completion
		close(ch)
		d.wg.Wait()

		close(completedCh)
		g.Wait()
	} else {
		// Checking if we read anything new
		for d.err == nil {
			d.getChunk(0, completedCh)
		}

		close(completedCh)
		g.Wait()

		// We expect a 416 error letting us know we are done downloading the
		// total bytes. Since we do not know the content's length, this will
		// keep grabbing chunks of data until the range of bytes specified in
		// the request is out of range of the content. Once, this happens, a
		// 416 should occur.
		var responseError interface {
			HTTPStatusCode() int
		}
		if errors.As(d.err, &responseError) {
			if responseError.HTTPStatusCode() == http.StatusRequestedRangeNotSatisfiable {
				d.err = nil
			}
		}
	}

	// Return error
	return d.written, d.err
}

func receiveChunks(w io.Writer, concurrency int, ch chan *dlchunk2) error {
	r := ring.New(concurrency)
	defer func() {
		// cleanup any remaining chunks
		r.Do(func(v interface{}) {
			if v == nil {
				return
			}
			chunk := v.(*dlchunk2)
			chunk.Cleanup()
		})
	}()

	seq := 0
	size := 0
	for chunk := range ch {
		if size == 0 {
			r.Value = chunk
		} else {
			// find any non-nil entry in the ring
			for ; r.Value == nil; r = r.Next() {
			}
			rChunk := r.Value.(*dlchunk2)
			// seek to the right place for our new chunk
			r = r.Move(chunk.seq - rChunk.seq)
			r.Value = chunk
		}
		size++

		// find the minimum entry in the ring
		min := math.MaxInt
		r.Do(func(v interface{}) {
			if v == nil {
				return
			}

			chunk := v.(*dlchunk2)
			if min > chunk.seq {
				min = chunk.seq
			}
		})

		// seek to min entry - note that we know the ring currently points at the new chunk
		r = r.Move(min - chunk.seq)

		// drain any ready chunks from the ring
		for {
			if r.Value == nil {
				r.Prev()
				break
			}
			rChunk := r.Value.(*dlchunk2)
			if rChunk.seq != seq {
				break
			}

			// erase the chunk from the ring
			size--
			seq++
			r.Value = nil

			// only return the chunk to the pool after clearing space here
			_, err := io.Copy(w, rChunk)
			rChunk.Cleanup()
			if err != nil {
				return err
			}

			r = r.Next()
		}
	}

	return nil
}

func (d *downloader2) nextPart(seq int, pos int64) (*dlchunk2, error) {
	part, err := d.cfg.partPool.Get(d.ctx)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		d.cfg.partPool.Put(part)
	}

	return &dlchunk2{
		seq:     seq,
		start:   pos,
		part:    part,
		cleanup: cleanup,
		size:    int64(len(*part)),
	}, nil
}

// downloadPart is an individual goroutine worker reading from the ch channel
// and performing a GetObject request on the data with a given byte range.
//
// If this is the first worker, this operation also resolves the total number
// of bytes to be read so that the worker manager knows when it is finished.
func (d *downloader2) downloadPart(ch chan *dlchunk2, completedCh chan *dlchunk2) {
	defer d.wg.Done()
	for {
		chunk, ok := <-ch
		if !ok {
			break
		}
		if d.getErr() != nil {
			// Drain the channel if there is an error, to prevent deadlocking
			// of download producer.
			continue
		}

		if err := d.downloadChunk(chunk, completedCh); err != nil {
			d.setErr(err)
		}
	}
}

// getChunk grabs a chunk of data from the body.
// Not thread safe. Should only used when grabbing data on a single thread.
func (d *downloader2) getChunk(seq int, completedCh chan *dlchunk2) {
	if d.getErr() != nil {
		return
	}

	chunk, err := d.nextPart(seq, d.pos)
	if err != nil {
		d.setErr(err)
		return
	}
	d.pos += d.cfg.PartSize

	if err := d.downloadChunk(chunk, completedCh); err != nil {
		d.setErr(err)
	}
}

// downloadRange downloads an Object given the passed in Byte-Range value.
// The chunk used down download the range will be configured for that range.
func (d *downloader2) downloadRange(seq int, completedCh chan *dlchunk2, rng string) {
	if d.getErr() != nil {
		return
	}

	chunk, err := d.nextPart(seq, d.pos)
	if err != nil {
		d.setErr(err)
	}
	// Ranges specified will short circuit the multipart download
	chunk.withRange = rng

	if err := d.downloadChunk(chunk, completedCh); err != nil {
		d.setErr(err)
	}

	// Update the position based on the amount of data received.
	d.pos = d.written
}

// downloadChunk downloads the chunk from s3
func (d *downloader2) downloadChunk(chunk *dlchunk2, completedCh chan *dlchunk2) error {
	defer func() {
		completedCh <- chunk
	}()

	var params s3.GetObjectInput
	awsutil.Copy(&params, d.in)

	// Get the next byte range of data
	params.Range = aws.String(chunk.ByteRange())

	var n int64
	var err error
	for retry := 0; retry <= d.partBodyMaxRetries; retry++ {
		n, err = d.tryDownloadChunk(&params, chunk)
		if err == nil {
			break
		}
		// Check if the returned error is an errReadingBody.
		// If err is errReadingBody this indicates that an error
		// occurred while copying the http response body.
		// If this occurs we unwrap the err to set the underlying error
		// and attempt any remaining retries.
		if bodyErr, ok := err.(*errReadingBody); ok {
			err = bodyErr.Unwrap()
		} else {
			return err
		}

		chunk.cur = 0

		d.cfg.Logger.Logf(logging.Debug,
			"object part body download interrupted %s, err, %v, retrying attempt %d",
			aws.ToString(params.Key), err, retry)
	}

	d.incrWritten(n)
	return err
}

func (d *downloader2) tryDownloadChunk(params *s3.GetObjectInput, w io.Writer) (int64, error) {
	cleanup := func() {}
	if d.cfg.BufferProvider != nil {
		w, cleanup = d.cfg.BufferProvider.GetReadFrom(w)
	}
	defer cleanup()

	resp, err := d.cfg.S3.GetObject(d.ctx, params, d.cfg.ClientOptions...)
	if err != nil {
		return 0, err
	}
	d.setTotalBytes(resp) // Set total if not yet set.

	n, err := io.Copy(w, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return n, &errReadingBody{err: err}
	}

	return n, nil
}

// getTotalBytes is a thread-safe getter for retrieving the total byte status.
func (d *downloader2) getTotalBytes() int64 {
	d.m.Lock()
	defer d.m.Unlock()

	return d.totalBytes
}

// setTotalBytes is a thread-safe setter for setting the total byte status.
// Will extract the object's total bytes from the Content-Range if the file
// will be chunked, or Content-Length. Content-Length is used when the response
// does not include a Content-Range. Meaning the object was not chunked. This
// occurs when the full file fits within the PartSize directive.
func (d *downloader2) setTotalBytes(resp *s3.GetObjectOutput) {
	d.m.Lock()
	defer d.m.Unlock()

	if d.totalBytes >= 0 {
		return
	}

	if resp.ContentRange == nil {
		// ContentRange is nil when the full file contents is provided, and
		// is not chunked. Use ContentLength instead.
		if resp.ContentLength > 0 {
			d.totalBytes = resp.ContentLength
			return
		}
	} else {
		parts := strings.Split(*resp.ContentRange, "/")

		total := int64(-1)
		var err error
		// Checking for whether or not a numbered total exists
		// If one does not exist, we will assume the total to be -1, undefined,
		// and sequentially download each chunk until hitting a 416 error
		totalStr := parts[len(parts)-1]
		if totalStr != "*" {
			total, err = strconv.ParseInt(totalStr, 10, 64)
			if err != nil {
				d.err = err
				return
			}
		}

		d.totalBytes = total
	}
}

func (d *downloader2) incrWritten(n int64) {
	d.m.Lock()
	defer d.m.Unlock()

	d.written += n
}

// getErr is a thread-safe getter for the error object
func (d *downloader2) getErr() error {
	d.m.Lock()
	defer d.m.Unlock()

	return d.err
}

// setErr is a thread-safe setter for the error object
func (d *downloader2) setErr(e error) {
	d.m.Lock()
	defer d.m.Unlock()

	d.err = e
}

// dlchunk2 represents a single chunk of data to write by the worker routine.
// This structure also implements an io.SectionReader style interface for
// io.WriterAt, effectively making it an io.SectionWriter (which does not
// exist).
type dlchunk2 struct {
	part    *[]byte
	cleanup func() // to release the part buffer to the pool
	seq     int
	start   int64 // for generating byte ranges
	size    int64 // for generating byte ranges
	cur     int   // for writing
	off     int   // for reading

	// specifies the byte range the chunk should be downloaded with.
	withRange string
}

func (c *dlchunk2) empty() bool {
	return c.cur <= c.off
}

func (c *dlchunk2) Read(p []byte) (n int, err error) {
	if c.empty() {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, (*c.part)[c.off:c.cur])
	c.off += n
	return n, nil
}

// Write wraps io.WriterAt for the dlchunk2, writing from the dlchunk2's start
// position to its end (or EOF).
//
// If a range is specified on the dlchunk2 the size will be ignored when writing.
// as the total size may not of be known ahead of time.
func (c *dlchunk2) Write(p []byte) (n int, err error) {
	if c.cur >= len(*c.part) && len(c.withRange) == 0 {
		return 0, io.EOF
	}

	n = len(p)
	copy((*c.part)[c.cur:], p)
	c.cur += n

	return n, err
}

func (c *dlchunk2) Cleanup() {
	c.cleanup()
}

// ByteRange returns a HTTP Byte-Range header value that should be used by the
// client to request the chunk's range.
func (c *dlchunk2) ByteRange() string {
	if len(c.withRange) != 0 {
		return c.withRange
	}

	return fmt.Sprintf("bytes=%d-%d", c.start, c.start+c.size-1)
}
