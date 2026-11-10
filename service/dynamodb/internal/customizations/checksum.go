package customizations

import (
	"context"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"strconv"

	smithy "github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// AddValidateResponseChecksumOptions provides the options for the
// AddValidateResponseChecksum middleware setup.
type AddValidateResponseChecksumOptions struct {
	Disable bool
}

// AddValidateResponseChecksum adds the Checksum to the middleware
// stack if checksum is not disabled.
func AddValidateResponseChecksum(stack *middleware.Stack, options AddValidateResponseChecksumOptions) error {
	if options.Disable {
		return nil
	}

	return stack.Deserialize.Add(&Checksum{}, middleware.After)
}

// Checksum provides a middleware to validate the DynamoDB response
// body's integrity by comparing the computed CRC32 checksum with the value
// provided in the HTTP response header.
type Checksum struct{}

// ID returns the middleware ID.
func (*Checksum) ID() string { return "DynamoDB:ResponseChecksumValidation" }

// HandleDeserialize implements the Deserialize middleware handle method.
func (m *Checksum) HandleDeserialize(
	ctx context.Context, input middleware.DeserializeInput, next middleware.DeserializeHandler,
) (
	output middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	output, metadata, err = next.HandleDeserialize(ctx, input)
	if err != nil {
		return output, metadata, err
	}

	resp, ok := output.RawResponse.(*smithyhttp.Response)
	if !ok {
		return output, metadata, &smithy.DeserializationError{
			Err: fmt.Errorf("unknown response type %T", output.RawResponse),
		}
	}

	expectChecksum, ok, err := getCRC32Checksum(resp.Header)
	if err != nil {
		return output, metadata, &smithy.DeserializationError{Err: err}
	}
	if !ok {
		return output, metadata, nil
	}

	resp.Body = wrapCRC32ChecksumValidate(expectChecksum, resp.Body)

	return output, metadata, nil
}

const crc32ChecksumHeader = "X-Amz-Crc32"

func getCRC32Checksum(header http.Header) (uint32, bool, error) {
	v := header.Get(crc32ChecksumHeader)
	if len(v) == 0 {
		return 0, false, nil
	}

	c, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return 0, false, fmt.Errorf("unable to parse checksum header %v, %w", v, err)
	}

	return uint32(c), true, nil
}

// crc32ChecksumValidate wraps a response body and fails the read that reaches
// EOF if the CRC32 checksum of the bytes read does not match the expected
// checksum.
type crc32ChecksumValidate struct {
	reader io.Reader
	closer io.Closer
	expect uint32
	hash   hash.Hash32
}

// wrapCRC32ChecksumValidate constructs a new crc32ChecksumValidate that will
// compute a running CRC32 checksum of the bytes read.
func wrapCRC32ChecksumValidate(checksum uint32, reader io.ReadCloser) *crc32ChecksumValidate {
	hash := crc32.NewIEEE()
	return &crc32ChecksumValidate{
		expect: checksum,
		reader: io.TeeReader(reader, hash),
		closer: reader,
		hash:   hash,
	}
}

// Read reads from the wrapped body. At EOF it returns a
// checksumMismatchError instead of io.EOF if the checksum does not match.
func (c *crc32ChecksumValidate) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	if err != io.EOF {
		return n, err
	}

	if actual := c.hash.Sum32(); actual != c.expect {
		return n, &checksumMismatchError{expect: c.expect, actual: actual}
	}

	return n, io.EOF
}

// Close closes the wrapped body.
func (c *crc32ChecksumValidate) Close() error {
	return c.closer.Close()
}

// checksumMismatchError is returned when the CRC32 checksum of a response
// body does not match the X-Amz-Crc32 header.
type checksumMismatchError struct {
	expect, actual uint32
}

func (e *checksumMismatchError) Error() string {
	return fmt.Sprintf("response did not match expected checksum, %d, %d", e.expect, e.actual)
}

// RetryableError marks the error as retryable. The body was likely corrupted
// in transit.
func (*checksumMismatchError) RetryableError() bool { return true }
