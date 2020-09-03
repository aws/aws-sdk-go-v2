package customizations

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCRC32ChecksumValidate(t *testing.T) {
	cases := map[string]struct {
		Reader      io.ReadCloser
		ExpectCRC32 uint32
		ExpectErr   string
	}{
		"empty reader": {
			Reader:      ioutil.NopCloser(&bytes.Buffer{}),
			ExpectCRC32: 0,
		},
		"wrong checksum": {
			Reader:      ioutil.NopCloser(bytes.NewBuffer([]byte("abc123"))),
			ExpectCRC32: 123456,
			ExpectErr:   "did not match",
		},
		"with closer": {
			Reader: &wasClosedReadCloser{
				Reader: bytes.NewBuffer([]byte("abc123")),
			},
			ExpectCRC32: 0xcf02bb5c,
		},
		"without closer": {
			Reader:      ioutil.NopCloser(bytes.NewBuffer([]byte("abc123"))),
			ExpectCRC32: 0xcf02bb5c,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			reader := wrapCRC32ChecksumValidate(c.ExpectCRC32, c.Reader)
			// Asserts
			io.Copy(ioutil.Discard, reader)

			err := reader.Close()
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if c, ok := c.Reader.(interface{ WasClosed() bool }); ok {
				if !c.WasClosed() {
					t.Errorf("expect original reader closed, but was not")
				}
			}
		})
	}
}

type wasClosedReadCloser struct {
	io.Reader
	closed bool
}

func (c *wasClosedReadCloser) WasClosed() bool {
	return c.closed
}

func (c *wasClosedReadCloser) Close() error {
	c.closed = true
	if v, ok := c.Reader.(io.Closer); ok {
		return v.Close()
	}
	return nil
}
