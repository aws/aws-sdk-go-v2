package customizations

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"

	smithy "github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// AddAcceptEncodingGzipOptions provides the options for the
// AddAcceptEncodingGzip middleware setup.
type AddAcceptEncodingGzipOptions struct {
	Disable bool
}

// AddAcceptEncodingGzip explicitly adds handling for accept-encoding GZIP
// middleware to the operation stack. This allows checksums to be correctly
// computed without disabling GZIP support.
func AddAcceptEncodingGzip(stack *middleware.Stack, options AddAcceptEncodingGzipOptions) {
	if options.Disable {
		stack.Finalize.Add(&DisableAcceptEncodingGzipMiddleware{}, middleware.Before)
		return
	}

	stack.Finalize.Add(&AcceptEncodingGzipMiddleware{}, middleware.Before)
	stack.Deserialize.Insert(&DecompressGzipMiddleware{}, "OperationDeserializer", middleware.After)
}

// DisableAcceptEncodingGzipMiddleware provides the middleware that will
// disable the underlying http client automatically enabling for gzip
// decompress content-encoding support.
type DisableAcceptEncodingGzipMiddleware struct{}

// ID returns the id for the middleware.
func (*DisableAcceptEncodingGzipMiddleware) ID() string {
	return "DynamoDB:DisableAcceptEncodingGzipMiddleware"
}

// HandleFinalize implements the FinalizeMiddlware interface.
func (*DisableAcceptEncodingGzipMiddleware) HandleFinalize(
	ctx context.Context, input middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	output middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := input.Request.(*smithyhttp.Request)
	if !ok {
		return output, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("unknown request type %T", input.Request),
		}
	}

	// Explicitly enable gzip support, this will prevent the http client from
	// auto extracting the zipped content.
	req.Header.Set("Accept-Encoding", "identity")

	return next.HandleFinalize(ctx, input)
}

// AcceptEncodingGzipMiddleware provides a middleware to enable support for
// gzip responses, with manual decompression. This prevents the underlying HTTP
// client from performing the gzip decompression automatically.
type AcceptEncodingGzipMiddleware struct{}

// ID returns the id for the middleware.
func (*AcceptEncodingGzipMiddleware) ID() string { return "DynamoDB:AcceptEncodingGzipMiddleware" }

// HandleFinalize implements the FinalizeMiddlware interface.
func (*AcceptEncodingGzipMiddleware) HandleFinalize(
	ctx context.Context, input middleware.FinalizeInput, next middleware.FinalizeHandler,
) (
	output middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := input.Request.(*smithyhttp.Request)
	if !ok {
		return output, metadata, &smithy.SerializationError{
			Err: fmt.Errorf("unknown request type %T", input.Request),
		}
	}

	// Explicitly enable gzip support, this will prevent the http client from
	// auto extracting the zipped content.
	req.Header.Set("Accept-Encoding", "gzip")

	return next.HandleFinalize(ctx, input)
}

// DecompressGzipMiddleware provides the middleware for decompressing a gzip
// response from the service.
type DecompressGzipMiddleware struct{}

// ID returns the id for the middleware.
func (*DecompressGzipMiddleware) ID() string { return "DynamoDB:DecompressGzipMiddleware" }

// HandleDeserialize implements the DeserializeMiddlware interface.
func (*DecompressGzipMiddleware) HandleDeserialize(
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
	if v := resp.Header.Get("Content-Encoding"); v != "gzip" {
		return output, metadata, err
	}

	// Clear content length since it will no longer be valid once the response
	// body is decompressed.
	resp.Header.Del("Content-Length")
	resp.ContentLength = -1

	resp.Body = wrapGzipReader(resp.Body)

	return output, metadata, err
}

type gzipReader struct {
	reader io.ReadCloser
	gzip   *gzip.Reader
}

func wrapGzipReader(reader io.ReadCloser) *gzipReader {
	return &gzipReader{
		reader: reader,
	}
}

// Read wraps the gzip reader around the underlying io.Reader to extract the
// response bytes on the fly.
func (g *gzipReader) Read(b []byte) (n int, err error) {
	if g.gzip == nil {
		g.gzip, err = gzip.NewReader(g.reader)
		if err != nil {
			g.gzip = nil // ensure uninitialized gzip value isn't used in close.
			return 0, fmt.Errorf("failed to decompress gzip response, %w", err)
		}
	}

	return g.gzip.Read(b)
}

func (g *gzipReader) Close() error {
	if g.gzip == nil {
		return nil
	}

	if err := g.gzip.Close(); err != nil {
		g.reader.Close()
		return fmt.Errorf("failed to decompress gzip response, %w", err)
	}

	return g.reader.Close()
}
