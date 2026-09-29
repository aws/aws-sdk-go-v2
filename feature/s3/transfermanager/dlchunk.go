package transfermanager

import (
	"io"

	internalio "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/io"
)

type dlChunk struct {
	start int64

	part      int32
	withRange string

	sink *internalio.AsyncWriterAt
}

func readChunk(r io.Reader, buf []byte) (int, error) {
	var n int
	for n < len(buf) {
		nr, err := r.Read(buf[n:])
		n += nr
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// ReadFrom reads the body into pooled buffers and queues them for writing.
//
// On a read error other than io.EOF the partially filled buffer is discarded
// rather than written: the chunk is retried from start anyway, and a short
// write at an arbitrary length is not valid on a file opened with O_DIRECT.
// The returned count only includes bytes that were queued for writing.
func (c *dlChunk) ReadFrom(r io.Reader) (int64, error) {
	var total int64
	for {
		buf := c.sink.Buffer()
		n, err := readChunk(r, buf)
		if err != nil && err != io.EOF {
			c.sink.Release(buf)
			return total, err
		}

		if n > 0 {
			c.sink.WriteAt(buf, n, c.start+total)
		} else {
			c.sink.Release(buf)
		}
		total += int64(n)

		if err == io.EOF {
			return total, nil
		}
	}
}
