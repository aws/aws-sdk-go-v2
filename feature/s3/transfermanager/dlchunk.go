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

func (c *dlChunk) ReadFrom(r io.Reader) (int64, error) {
	var total int64
	for {
		buf := c.sink.Buffer()
		n, err := readChunk(r, buf)
		off := c.start + total
		if n > 0 {
			c.sink.WriteAt(buf, n, off)
		} else {
			c.sink.Release(buf)
		}
		total += int64(n)

		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return total, err
		}
	}
}
