package sdk

import (
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/aws/aws-sdk-go-v2/internal/rand"
)

type byteReader byte

func (b byteReader) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		p[i] = byte(b)
	}
	return len(p), nil
}

type errorReader struct{ err error }

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestUUIDVersion4(t *testing.T) {
	origReader := rand.Reader
	defer func() { rand.Reader = origReader }()

	cases := map[string]struct {
		Expect string
		Reader io.Reader
		Err    string
	}{
		"0x00": {
			Expect: `00000000-0000-4000-8000-000000000000`,
			Reader: byteReader(0),
		},
		"0x01": {
			Expect: `01010101-0101-4101-8101-010101010101`,
			Reader: byteReader(1),
		},
		"partial": {
			Expect: `01010101-0101-4101-8101-010101010101`,
			Reader: iotest.HalfReader(byteReader(1)),
		},
		"error": {
			Reader: errorReader{err: io.ErrUnexpectedEOF},
			Err:    io.ErrUnexpectedEOF.Error(),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			rand.Reader = c.Reader

			uuid, err := UUIDVersion4()
			if len(c.Err) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.Err, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %q in error, %q", e, a)
				}
			} else if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := c.Expect, uuid; e != a {
				t.Errorf("expect %v uuid, got %v", e, a)
			}
		})
	}
}
