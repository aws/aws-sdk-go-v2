package eventstream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"testing"
)

// testHeadersLen is an arbitrary non-zero headers length. A message of
// minMsgLen+testHeadersLen bytes has room for exactly testHeadersLen bytes of
// headers and no payload.
const testHeadersLen = 4

func TestValidateLens(t *testing.T) {
	cases := map[string]struct {
		length, headersLen uint32
		wantErr            bool
	}{
		"zero length":                  {0, 0, true},
		"length below minimum":         {1, 0, true},
		"length one below minimum":     {minMsgLen - 1, 0, true},
		"minimum length":               {minMsgLen, 0, false},
		"headers fill message":         {minMsgLen + testHeadersLen, testHeadersLen, false},
		"headers exceed message":       {minMsgLen + testHeadersLen, testHeadersLen + 1, true},
		"headers equal total length":   {minMsgLen, minMsgLen, true},
		"max length, no headers":       {maxMsgLen, 0, false},
		"max length, max headers":      {maxMsgLen, maxMsgLen - minMsgLen, false},
		"max length, headers too long": {maxMsgLen, maxMsgLen - minMsgLen + 1, true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p := messagePrelude{Length: c.length, HeadersLen: c.headersLen}
			err := p.ValidateLens()
			if c.wantErr {
				if !errors.As(err, &LengthError{}) {
					t.Fatalf("expect LengthError, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if want, got := c.length-c.headersLen-minMsgLen, p.PayloadLen(); want != got {
				t.Errorf("expect payload length %d, got %d", want, got)
			}
		})
	}
}

// TestDecode_InconsistentPreludeLengths ensures a prelude whose lengths would
// underflow the payload length is rejected before any headers or payload are
// read, even when its prelude checksum is valid.
func TestDecode_InconsistentPreludeLengths(t *testing.T) {
	cases := map[string]struct{ length, headersLen uint32 }{
		"length below minimum":   {1, 0},
		"headers exceed message": {minMsgLen + testHeadersLen, testHeadersLen + 1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			prelude := make([]byte, preludeLen+preludeCRCLen)
			binary.BigEndian.PutUint32(prelude[0:], c.length)
			binary.BigEndian.PutUint32(prelude[4:], c.headersLen)
			binary.BigEndian.PutUint32(prelude[8:], crc32.ChecksumIEEE(prelude[:preludeLen]))

			// Follow the prelude with more data than a valid message of this
			// prelude could hold, so reading any of it would show the decoder
			// trusted the lengths. The exact size does not matter.
			const trailingStreamLen = 1 << 20 // 1 MiB
			r := &countingReader{r: io.MultiReader(
				bytes.NewReader(prelude),
				bytes.NewReader(make([]byte, trailingStreamLen)),
			)}

			_, err := NewDecoder().Decode(r, nil)
			if !errors.As(err, &LengthError{}) {
				t.Fatalf("expect LengthError, got %v", err)
			}
			if r.n > int64(len(prelude)) {
				t.Errorf("expect at most %d bytes read, got %d", len(prelude), r.n)
			}
		})
	}
}

type countingReader struct {
	r io.Reader
	n int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	return n, err
}
