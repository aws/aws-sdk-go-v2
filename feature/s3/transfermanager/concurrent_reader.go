package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var errReadClosed = errors.New("read on closed GetObject body")

// concurrentReader downloads object parts with a pool of workers and returns
// them to the caller in order. The workers start on the first Read and live
// until the download ends. At most window parts are in flight or buffered, and
// the next part starts as soon as the caller finishes reading one.
type concurrentReader struct {
	options Options
	in      *GetObjectInput

	partSize   int64
	partsCount int32
	window     int32
	totalBytes int64 // end offset of the requested range, exclusive
	pos        int64 // start offset of the requested range
	etag       *string

	ch   chan outChunk // downloaded parts, from workers to Read
	work chan getChunk // part requests, from dispatch to workers
	sem  chan struct{} // one slot for each part in flight or buffered
	done chan struct{} // closed once when the download ends

	ctx    context.Context
	cancel context.CancelFunc

	buf       map[int32]*outChunk // received parts not yet fully read
	nextIndex int32
	consumed  int64 // bytes returned to the caller

	startOnce sync.Once
	closeOnce sync.Once
	err       atomic.Pointer[error]
}

// getWindow returns how many parts can be in flight at once: at most
// concurrency and bufferSize/partSize, and at least 1.
func getWindow(concurrency int, bufferSize, partSize int64) int32 {
	return int32(max(1, min(int64(concurrency), bufferSize/max(partSize, 1))))
}

// Read implements io.Reader, returning the object's bytes in order. It waits
// only when it has no bytes to return.
func (r *concurrentReader) Read(p []byte) (int, error) {
	r.start()

	if len(p) == 0 {
		return 0, nil
	}

	var written int
	for written < len(p) {
		c, ok := r.buf[r.nextIndex]
		if !ok {
			if r.nextIndex >= r.partsCount {
				r.teardown()
				return written, io.EOF
			}
			if written > 0 {
				return written, nil
			}
			select {
			case <-r.done:
				return written, r.terminalErr()
			case oc := <-r.ch:
				r.buf[oc.index] = &oc
			}
			continue
		}

		if c.cur == 0 {
			if err := r.checkReassembly(c); err != nil {
				r.fail(err)
				return written, err
			}
		}

		n, _ := c.body.Read(p[written:])
		c.cur += int64(n)
		r.consumed += int64(n)
		written += n

		if c.cur >= c.length {
			delete(r.buf, r.nextIndex)
			r.nextIndex++
			<-r.sem
			if r.nextIndex >= r.partsCount {
				r.teardown()
				return written, io.EOF
			}
		}
	}

	return written, nil
}

// Close stops the download. Read then fails, possibly after returning bytes
// that were already downloaded.
func (r *concurrentReader) Close() error {
	r.fail(errReadClosed)
	return nil
}

// initReader prepares r so that Close is safe before the first Read.
func (r *concurrentReader) initReader(ctx context.Context) {
	r.ctx, r.cancel = context.WithCancel(ctx)
	r.done = make(chan struct{})
	r.buf = make(map[int32]*outChunk)
}

// start launches the workers and the dispatcher on the first Read, so a
// caller that never reads launches nothing.
func (r *concurrentReader) start() {
	r.startOnce.Do(func() {
		select {
		case <-r.done:
			return
		default:
		}

		r.ch = make(chan outChunk, r.window)
		r.work = make(chan getChunk)
		r.sem = make(chan struct{}, r.window)

		clientOptions := []func(*s3.Options){
			func(o *s3.Options) {
				o.APIOptions = append(o.APIOptions,
					middleware.AddSDKAgentKeyValue(middleware.FeatureMetadata, userAgentKey, goModuleVersion),
					addFeatureUserAgent,
				)
			}}

		for range r.window {
			go r.worker(clientOptions...)
		}
		go r.dispatch()
		go r.watch()
	})
}

// watch ends the download when ctx is canceled, since the workers and the
// dispatcher only wait on done.
func (r *concurrentReader) watch() {
	select {
	case <-r.ctx.Done():
		r.fail(r.ctx.Err())
	case <-r.done:
	}
}

// dispatch takes a window slot for each part, then hands the part to a worker.
func (r *concurrentReader) dispatch() {
	pos := r.pos
	for i := int32(0); i < r.partsCount; i++ {
		select {
		case <-r.done:
			return
		case r.sem <- struct{}{}:
		}

		var gc getChunk
		if r.options.GetObjectType == types.GetObjectParts {
			gc = getChunk{part: i + 1, index: i}
		} else {
			gc = getChunk{withRange: r.byteRange(pos), index: i}
		}
		pos += r.partSize

		select {
		case <-r.done:
			return
		case r.work <- gc:
		}
	}
	close(r.work)
}

func (r *concurrentReader) worker(clientOptions ...func(*s3.Options)) {
	for {
		select {
		case <-r.done:
			return
		case chunk, ok := <-r.work:
			if !ok {
				return
			}
			r.downloadChunk(chunk, clientOptions...)
		}
	}
}

// downloadChunk downloads one part and sends it to Read. Like DownloadObject,
// it makes up to PartBodyMaxRetries attempts when reading the body fails.
func (r *concurrentReader) downloadChunk(chunk getChunk, clientOptions ...func(*s3.Options)) {
	params := r.in.mapGetObjectInput(!r.options.DisableChecksumValidation)
	if chunk.part != 0 {
		params.PartNumber = aws.Int32(chunk.part)
	}
	if chunk.withRange != "" {
		params.Range = aws.String(chunk.withRange)
	}
	if params.VersionId == nil {
		params.IfMatch = r.etag
	}

	var buf []byte
	var respStart, respEnd int64
	var err error
	for range max(r.options.PartBodyMaxRetries, 1) {
		buf, respStart, respEnd, err = r.tryDownloadChunk(params, clientOptions...)
		if _, ok := err.(*errReadingBody); !ok {
			break
		}
	}
	if err != nil {
		r.fail(err)
		return
	}

	select {
	case <-r.done:
	case r.ch <- outChunk{body: bytes.NewReader(buf), index: chunk.index, length: int64(len(buf)), rangeStart: respStart, rangeEnd: respEnd}:
	}
}

// tryDownloadChunk makes one GetObject call and reads the part body. It also
// returns the response's Content-Range bounds, or -1 if there is none.
func (r *concurrentReader) tryDownloadChunk(params *s3.GetObjectInput, clientOptions ...func(*s3.Options)) ([]byte, int64, int64, error) {
	out, err := r.options.S3.GetObject(r.ctx, params, clientOptions...)
	if err != nil {
		return nil, -1, -1, err
	}
	defer out.Body.Close()

	respStart, respEnd := int64(-1), int64(-1)
	if out.ContentRange != nil {
		respStart, respEnd, err = getRespRange(aws.ToString(out.ContentRange))
		if err != nil {
			return nil, -1, -1, err
		}
		if params.Range != nil {
			reqStart, reqEnd, err := getReqRange(aws.ToString(params.Range))
			if err != nil {
				return nil, -1, -1, err
			}
			// don't validate first chunk since object size is unknown when getting that
			if reqStart != 0 && (reqStart != respStart || reqEnd != respEnd) {
				return nil, -1, -1, fmt.Errorf("range mismatch between request %d-%d and response %d-%d", reqStart, reqEnd, respStart, respEnd)
			}
		}
	}

	buf := make([]byte, aws.ToInt64(out.ContentLength))
	_, err = io.ReadFull(out.Body, buf)
	if err == nil {
		// read to io.EOF so the client validates the response checksum
		_, err = io.Copy(io.Discard, out.Body)
	}
	if err != nil {
		return nil, -1, -1, &errReadingBody{err: err}
	}
	return buf, respStart, respEnd, nil
}

// checkReassembly verifies that a part's Content-Range starts where the stream
// is and covers exactly the bytes received. Parts can have unequal sizes, so
// the response range is the only reliable record of where a part belongs.
func (r *concurrentReader) checkReassembly(c *outChunk) error {
	if c.rangeStart < 0 {
		return nil
	}
	if e, a := r.pos+r.consumed, c.rangeStart; e != a {
		return fmt.Errorf("range mismatch between computed start position %d and actual start position %d for part %d", e, a, c.index+1)
	}
	if e, a := c.rangeStart+c.length-1, c.rangeEnd; c.length > 0 && e != a {
		return fmt.Errorf("range mismatch between computed end position %d and actual end position %d for part %d", e, a, c.index+1)
	}
	return nil
}

// byteRange returns an HTTP Byte-Range header value for the part starting at pos.
func (r *concurrentReader) byteRange(pos int64) string {
	return fmt.Sprintf("bytes=%d-%d", pos, min(r.totalBytes-1, pos+r.partSize-1))
}

// fail records the first error and ends the download.
func (r *concurrentReader) fail(err error) {
	r.err.CompareAndSwap(nil, &err)
	r.teardown()
}

func (r *concurrentReader) loadErr() error {
	if p := r.err.Load(); p != nil {
		return *p
	}
	return nil
}

// terminalErr is what a waiting Read returns once the download has ended.
func (r *concurrentReader) terminalErr() error {
	if err := r.loadErr(); err != nil {
		return err
	}
	return io.EOF
}

// teardown cancels in-flight requests and closes done, once.
func (r *concurrentReader) teardown() {
	r.closeOnce.Do(func() {
		r.cancel()
		close(r.done)
	})
}

type getChunk struct {
	part      int32
	withRange string

	index int32
}

type outChunk struct {
	body  io.Reader
	index int32

	rangeStart int64
	rangeEnd   int64
	length     int64
	cur        int64
}
