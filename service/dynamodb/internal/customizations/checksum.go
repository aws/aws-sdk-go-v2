package customizations

import (
	"context"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"strconv"

	smithy "github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

type checksumClientOptions interface {
	GetDisableResponseChecksumValidation() bool
}

// AddChecksumMiddleware adds the ChecksumMiddleware to the middleware stack if
// checksum is not disabled.
func AddChecksumMiddleware(stack *middleware.Stack, options checksumClientOptions) {
	if options.GetDisableResponseChecksumValidation() {
		return
	}

	stack.Deserialize.Insert(&ChecksumMiddleware{}, "OperationDeserializer", middleware.After)
}

// ChecksumMiddleware provides a middleware to validate the DynamoDB response
// body's integrity by comparing the computed CRC32 checksum with the value
// provided in the HTTP response header.
type ChecksumMiddleware struct{}

// ID returns the middleware ID.
func (*ChecksumMiddleware) ID() string { return "DynamoDBResponseChecksumValidation" }

// HandleDeserialize implements the Deserialize middleware handle method.
func (m *ChecksumMiddleware) HandleDeserialize(
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
		return output, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("unknown response type %T", output.RawResponse),
		}
	}

	expectChecksum, ok, err := getCRC32Checksum(resp.Header)
	if err != nil {
		return output, metadata, &smithy.SerializationError{Err: err}
	}

	resp.Body = wrapCRC32ChecksumValidate(expectChecksum, resp.Body)

	return output, metadata, err
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

// crc32ChecksumValidate provides wrapping of an io.Reader to validate the CRC32
// checksum of the bytes read against the expected checksum.
type crc32ChecksumValidate struct {
	expect uint32
	reader io.Reader
	hash   hash.Hash32
}

// wrapCRC32ChecksumValidate constructs a new crc32ChecksumValidate that will compute a
// running CRC32 checksum of the bytes read.
func wrapCRC32ChecksumValidate(checksum uint32, reader io.Reader) *crc32ChecksumValidate {
	return &crc32ChecksumValidate{
		expect: checksum,
		reader: reader,
		hash:   crc32.NewIEEE(),
	}
}

// Read reads the wrapped reader, and updates the CRC32 checksum validation
// with with the running checksum of the read bytes. Returns number of bytes
// read, and error if wrapped reader returns error.
func (c *crc32ChecksumValidate) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	if n > 0 {
		c.hash.Write(p[:n])
	}

	return n, err
}

// Close validates the wrapped reader's CRC32 checksum. Returns an error if
// the read checksum does not match the expected checksum.
//
// May return an error if the wrapped io.Reader's close returns an error, if it
// implements close.
func (c *crc32ChecksumValidate) Close() error {
	closer, isCloser := c.reader.(io.Closer)

	if actual := c.hash.Sum32(); actual != c.expect {
		if isCloser {
			defer closer.Close()
		}
		return fmt.Errorf("checksum 0x%x does not match expected 0x%x", actual, c.expect)
	}

	if isCloser {
		return closer.Close()
	}
	return nil
}
