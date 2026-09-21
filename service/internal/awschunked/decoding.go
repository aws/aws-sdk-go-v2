package awschunked

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// UnknownDecodedLength signals that the caller does not know the decoded
// payload length ahead of time, disabling the decoded-length verification.
const UnknownDecodedLength int64 = -1

const (
	// crlf is the carriage-return line-feed sequence that frames every chunk
	// size line, chunk data, and trailer field in the aws-chunked encoding.
	crlf = "\r\n"

	// trailerKeyValueSeparator separates a trailer field name from its value.
	trailerKeyValueSeparator = ":"
)

// Reader decodes an inbound aws-chunked encoded stream,
// yielding only the decoded payload bytes to the caller. Chunk framing
// (sizes, chunk extensions, and CRLF separators) is consumed internally and
// never surfaced.
//
// The reader supports the unsigned aws-chunked framing used on the S3 response
// path:
//
//	<hex-size>[;chunk-ext]\r\n
//	<chunk-data>\r\n
//	...
//	0[;chunk-ext]\r\n
//	[trailer-name:trailer-value\r\n]...
//	\r\n
//
// Chunk extensions (the optional ";"-delimited parameters following the chunk
// size) are parsed but ignored: the decoder makes no assumption about their
// content for forward compatibility. This reader does not verify per-chunk
// signatures; any integrity metadata is delivered in the trailer and left for
// a later stage to consume.
//
// Trailer values are only available after the terminating zero-length chunk
// has been decoded. They are surfaced through Trailers, which is populated
// once the reader reaches io.EOF. Callers MUST NOT read trailers before the
// stream has been fully consumed.
type Reader struct {
	source io.ReadCloser
	reader *bufio.Reader

	// expectedDecodedLength is the value of x-amz-decoded-content-length, used
	// to guard against truncated responses. UnknownDecodedLength disables the
	// check.
	expectedDecodedLength int64
	decodedLength         int64

	// remaining bytes left to read out of the current data chunk, excluding
	// the trailing CRLF that frames it.
	remaining int64

	trailers     http.Header
	sawLastChunk bool
	done         bool
	err          error
}

// Options configures a Reader. Use the functional option functions passed to
// NewReader to set fields.
type Options struct {
	// DecodedLength is the expected length of the decoded payload, from the
	// x-amz-decoded-content-length header. Set to UnknownDecodedLength to
	// disable verification.
	DecodedLength int64
}

// NewReader returns a reader that decodes the outermost
// aws-chunked framing of source, surfacing only decoded payload bytes.
func NewReader(
	source io.ReadCloser,
	optFns ...func(*Options),
) *Reader {
	options := Options{
		DecodedLength: UnknownDecodedLength,
	}
	for _, fn := range optFns {
		fn(&options)
	}

	return &Reader{
		source:                source,
		reader:                bufio.NewReader(source),
		expectedDecodedLength: options.DecodedLength,
		trailers:              http.Header{},
	}
}

// Trailers returns the trailing headers decoded from the terminating chunk.
// The returned header is only fully populated once the reader has reached
// io.EOF; before that it is empty.
func (r *Reader) Trailers() http.Header {
	return r.trailers
}

// Read yields decoded payload bytes. Framing is consumed internally. When the
// terminating chunk and trailers have been consumed, Read returns io.EOF and
// the trailers become available.
func (r *Reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.done {
		return 0, io.EOF
	}

	// If we're between chunks, advance to the next data chunk (or the
	// terminator + trailers).
	if r.remaining == 0 {
		if err := r.beginNextChunk(); err != nil {
			r.err = err
			return 0, err
		}
		if r.done {
			return 0, io.EOF
		}
	}

	// Read out of the current data chunk, capped by what remains in it.
	toRead := int64(len(p))
	if toRead > r.remaining {
		toRead = r.remaining
	}
	n, err := r.reader.Read(p[:toRead])
	r.remaining -= int64(n)
	r.decodedLength += int64(n)

	if err == io.EOF {
		// Underlying stream ended mid-chunk: framing is truncated.
		r.err = &decodeError{Msg: "unexpected EOF within chunk data"}
		return n, r.err
	}
	if err != nil {
		r.err = err
		return n, err
	}

	// When a data chunk is exhausted, consume its trailing CRLF now so the
	// next Read starts cleanly at the next chunk-size line.
	if r.remaining == 0 {
		if err := r.consumeCRLF(); err != nil {
			r.err = err
			return n, err
		}
	}

	return n, nil
}

// beginNextChunk parses the next chunk-size line. For a non-zero size it sets
// remaining so subsequent Reads yield that chunk's data. For a zero size it
// consumes the trailers and finalizes the stream.
func (r *Reader) beginNextChunk() error {
	line, err := r.readLine()
	if err != nil {
		return err
	}

	size, err := parseChunkSize(line)
	if err != nil {
		return err
	}

	if size == 0 {
		if r.sawLastChunk {
			return &decodeError{Msg: "unexpected additional zero-length chunk"}
		}
		r.sawLastChunk = true
		return r.finish()
	}

	r.remaining = size
	return nil
}

// finish is called after the terminating zero-length chunk. It parses the
// trailer section, verifies the decoded length, and marks the stream done.
func (r *Reader) finish() error {
	if err := r.readTrailers(); err != nil {
		return err
	}
	if err := r.verifyNoTrailingData(); err != nil {
		return err
	}
	if err := r.verifyDecodedLength(); err != nil {
		return err
	}

	r.done = true
	return nil
}

// readTrailers reads zero or more "name:value" trailer lines terminated by a
// blank line.
func (r *Reader) readTrailers() error {
	for {
		line, err := r.readLine()
		if err != nil {
			return err
		}
		if line == "" {
			return nil
		}

		idx := strings.Index(line, trailerKeyValueSeparator)
		if idx < 0 {
			return &decodeError{Msg: fmt.Sprintf("malformed trailer, missing separator: %q", line)}
		}
		name := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if name == "" {
			return &decodeError{Msg: fmt.Sprintf("malformed trailer, empty name: %q", line)}
		}
		r.trailers.Add(name, value)
	}
}

// verifyNoTrailingData ensures the stream ends immediately after the trailer
// section's terminating blank line.
func (r *Reader) verifyNoTrailingData() error {
	if _, err := r.reader.ReadByte(); err != io.EOF {
		if err == nil {
			return &decodeError{Msg: "unexpected data after aws-chunked terminator"}
		}
		return err
	}
	return nil
}

func (r *Reader) verifyDecodedLength() error {
	if r.expectedDecodedLength == UnknownDecodedLength {
		return nil
	}
	if r.decodedLength != r.expectedDecodedLength {
		return &decodeError{Msg: fmt.Sprintf(
			"decoded content length mismatch, expected %d, got %d",
			r.expectedDecodedLength, r.decodedLength)}
	}
	return nil
}

// readLine reads a single CRLF-terminated line, returning it without the
// trailing CRLF. A bare LF or a line missing its CRLF is treated as a decode
// error.
func (r *Reader) readLine() (string, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", &decodeError{Msg: "unexpected EOF reading aws-chunked line"}
		}
		return "", err
	}
	if !strings.HasSuffix(line, crlf) {
		return "", &decodeError{Msg: "malformed aws-chunked line, missing CRLF"}
	}
	return line[:len(line)-len(crlf)], nil
}

// consumeCRLF reads the CRLF that frames a data chunk.
func (r *Reader) consumeCRLF() error {
	b := make([]byte, 2)
	if _, err := io.ReadFull(r.reader, b); err != nil {
		return &decodeError{Msg: "malformed aws-chunked chunk, missing trailing CRLF"}
	}
	if string(b) != crlf {
		return &decodeError{Msg: "malformed aws-chunked chunk, invalid trailing CRLF"}
	}
	return nil
}

// Close closes the underlying reader. It does NOT perform validation; an
// incomplete read (Close before EOF) never triggers checksum comparison.
func (r *Reader) Close() error {
	return r.source.Close()
}

// parseChunkSize parses a chunk-size line of the form "<hex>[;ext[;ext...]]".
// Chunk extensions after the first ";" are parsed off but ignored.
func parseChunkSize(line string) (int64, error) {
	sizeStr := line
	if idx := strings.IndexByte(line, ';'); idx >= 0 {
		sizeStr = line[:idx]
		// remainder (line[idx+1:]) is the chunk extension(s); ignored
	}
	sizeStr = strings.TrimSpace(sizeStr)
	if sizeStr == "" {
		return 0, &decodeError{Msg: "malformed aws-chunked chunk, empty size"}
	}
	size, err := strconv.ParseInt(sizeStr, 16, 64)
	if err != nil {
		return 0, &decodeError{Msg: fmt.Sprintf("invalid aws-chunked chunk size %q", sizeStr)}
	}
	if size < 0 {
		return 0, &decodeError{Msg: fmt.Sprintf("negative aws-chunked chunk size %q", sizeStr)}
	}
	return size, nil
}

// decodeError represents a malformed aws-chunked stream. It is surfaced from
// Read so callers (and, downstream, deserialization middleware) can treat it
// as a deserialization error.
type decodeError struct {
	Msg string
}

func (e *decodeError) Error() string {
	return "aws-chunked decode error: " + e.Msg
}
