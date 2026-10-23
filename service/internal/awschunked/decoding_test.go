package awschunked

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"testing/iotest"
)

// buildReader constructs the decoding reader over the given raw framed body,
// optionally wrapping the source so bytes are delivered one at a time to
// exercise every internal state boundary.
func newDecoderForTest(raw string, decodedLen int64, oneByte bool) *Reader {
	var src io.Reader = strings.NewReader(raw)
	if oneByte {
		src = iotest.OneByteReader(src)
	}
	return NewReader(io.NopCloser(src), func(o *Options) {
		o.DecodedLength = decodedLen
	})
}

func TestAWSChunkedDecoding_Success(t *testing.T) {
	cases := map[string]struct {
		raw          string
		decodedLen   int64
		expectBody   string
		expectTrail  map[string]string
		trailerOrder []string
	}{
		"single chunk with trailer": {
			raw:         "5\r\nhello\r\n0\r\ntrailer-key:abc123==\r\n\r\n",
			decodedLen:  5,
			expectBody:  "hello",
			expectTrail: map[string]string{"Trailer-Key": "abc123=="},
		},
		"multiple chunks with trailer": {
			raw:         "5\r\nhello\r\n6\r\n world\r\n0\r\ntrailer-key:AAAA==\r\n\r\n",
			decodedLen:  11,
			expectBody:  "hello world",
			expectTrail: map[string]string{"Trailer-Key": "AAAA=="},
		},
		"chunk data containing CRLF": {
			raw:        "5\r\na\r\nbc\r\n0\r\n\r\n",
			decodedLen: 5,
			expectBody: "a\r\nbc",
		},
		"trailer value surrounded by whitespace": {
			raw:         "3\r\nfoo\r\n0\r\ntrailer-key:  abc123==  \r\n\r\n",
			decodedLen:  3,
			expectBody:  "foo",
			expectTrail: map[string]string{"Trailer-Key": "abc123=="},
		},
		"chunk extensions ignored": {
			raw:         "5;chunk-signature=abc;date=20260101\r\nhello\r\n0;chunk-signature=def\r\ntrailer-key:xyz==\r\n\r\n",
			decodedLen:  5,
			expectBody:  "hello",
			expectTrail: map[string]string{"Trailer-Key": "xyz=="},
		},
		"multiple trailers": {
			raw:        "5\r\nhello\r\n0\r\nfirst-trailer:value1\r\nsecond-trailer:value2\r\n\r\n",
			decodedLen: 5,
			expectBody: "hello",
			expectTrail: map[string]string{
				"First-Trailer":  "value1",
				"Second-Trailer": "value2",
			},
		},
		"empty data zero length content": {
			raw:         "0\r\ntrailer-key:abc==\r\n\r\n",
			decodedLen:  0,
			expectBody:  "",
			expectTrail: map[string]string{"Trailer-Key": "abc=="},
		},
		"no trailer": {
			raw:        "5\r\nhello\r\n0\r\n\r\n",
			decodedLen: 5,
			expectBody: "hello",
		},
		"trailer fixture": {
			raw:         "b\r\nHello world\r\n0\r\ntrailer-key:abc123==\r\n\r\n",
			decodedLen:  11,
			expectBody:  "Hello world",
			expectTrail: map[string]string{"Trailer-Key": "abc123=="},
		},
	}

	for name, c := range cases {
		// Feed each success case one byte at a time and confirm the result is
		// identical to reading it all at once. This exercises every internal
		// state boundary and guards against buffering bugs.
		for _, oneByte := range []bool{false, true} {
			mode := "all-at-once"
			if oneByte {
				mode = "one-byte"
			}
			t.Run(name+"/"+mode, func(t *testing.T) {
				dec := newDecoderForTest(c.raw, c.decodedLen, oneByte)

				body, err := ioutil.ReadAll(dec)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got := string(body); got != c.expectBody {
					t.Fatalf("body mismatch: want %q got %q", c.expectBody, got)
				}
				for k, v := range c.expectTrail {
					if got := dec.Trailers().Get(k); got != v {
						t.Errorf("trailer %q: want %q got %q", k, v, got)
					}
				}
			})
		}
	}
}

func TestAWSChunkedDecoding_TrailerNotAvailableBeforeEOF(t *testing.T) {
	raw := "5\r\nhello\r\n0\r\ntrailer-key:abc==\r\n\r\n"
	dec := newDecoderForTest(raw, 5, false)

	// Read only the first two bytes; trailer must not be populated yet.
	buf := make([]byte, 2)
	if _, err := dec.Read(buf); err != nil {
		t.Fatalf("unexpected error on partial read: %v", err)
	}
	if v := dec.Trailers().Get("trailer-key"); v != "" {
		t.Fatalf("trailer available before EOF: %q", v)
	}

	// Drain to EOF; now it must be present.
	if _, err := io.ReadAll(dec); err != nil {
		t.Fatalf("unexpected error draining: %v", err)
	}
	if v := dec.Trailers().Get("trailer-key"); v != "abc==" {
		t.Fatalf("trailer missing after EOF: %q", v)
	}
}

func TestAWSChunkedDecoding_PartialReadNoError(t *testing.T) {
	// Incomplete read: caller stops early and closes. This must not error.
	raw := "b\r\nHello world\r\n0\r\ntrailer-key:abc123==\r\n\r\n"
	dec := newDecoderForTest(raw, 11, false)

	buf := make([]byte, 5)
	if _, err := dec.Read(buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dec.Close(); err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}
}

func TestAWSChunkedDecoding_Errors(t *testing.T) {
	cases := map[string]struct {
		raw        string
		decodedLen int64
	}{
		"invalid hex in chunk size": {
			raw:        "ZZ\r\n",
			decodedLen: UnknownDecodedLength,
		},
		"missing CRLF after chunk data": {
			raw:        "3\r\nfooX",
			decodedLen: UnknownDecodedLength,
		},
		"malformed trailer no colon": {
			raw:        "0\r\nnot-a-valid-trailer\r\n\r\n",
			decodedLen: UnknownDecodedLength,
		},
		"data after terminal chunk and trailer": {
			raw:        "3\r\nfoo\r\n0\r\ntrailer-key:abc==\r\n\r\nextra garbage",
			decodedLen: UnknownDecodedLength,
		},
		"decoded content length mismatch": {
			raw:        "5\r\nhello\r\n0\r\ntrailer-key:abc==\r\n\r\n",
			decodedLen: 100,
		},
		"multiple zero size chunks": {
			raw:        "0\r\n0\r\n\r\n",
			decodedLen: UnknownDecodedLength,
		},
	}

	for name, c := range cases {
		for _, oneByte := range []bool{false, true} {
			mode := "all-at-once"
			if oneByte {
				mode = "one-byte"
			}
			t.Run(name+"/"+mode, func(t *testing.T) {
				dec := newDecoderForTest(c.raw, c.decodedLen, oneByte)
				_, err := io.ReadAll(dec)
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			})
		}
	}
}
