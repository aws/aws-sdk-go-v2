package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io"
)

func AddRequestCompression(stack *middleware.Stack, DisableRequestCompression bool, RequestMinCompressSizeBytes int64) error {
	return stack.Build.Add(&requestCompression{
		disableRequestCompression:   DisableRequestCompression,
		requestMinCompressSizeBytes: RequestMinCompressSizeBytes,
	}, middleware.After)
}

type requestCompression struct {
	disableRequestCompression   bool
	requestMinCompressSizeBytes int64
}

func (m requestCompression) ID() string {
	return "RequestCompression"
}

func (m requestCompression) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	if m.disableRequestCompression {
		return next.HandleBuild(ctx, in)
	}
	// still need to check requestMinCompressSizeBytes in case it is out of range after service client config
	if m.requestMinCompressSizeBytes < 0 || m.requestMinCompressSizeBytes > 10485760 {
		return out, metadata, fmt.Errorf("invalid range for min request compression size bytes %d, must be within 0 and 10485760 inclusively", m.requestMinCompressSizeBytes)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	var isCompressed bool
	if stream := req.GetStream(); stream != nil {
		compressedBytes, err := compress(stream)
		if err != nil {
			return out, metadata, fmt.Errorf("failed to compress request stream, %v", err)
		}

		if newReq, err := req.SetStream(bytes.NewReader(compressedBytes)); err != nil {
			return out, metadata, fmt.Errorf("failed to set request stream, %v", err)
		} else {
			*req = *newReq
		}
		isCompressed = true
	} else if req.ContentLength >= m.requestMinCompressSizeBytes {
		compressedBytes, err := compress(req.Body)
		if err != nil {
			return out, metadata, fmt.Errorf("failed to compress request body, %v", err)
		}

		isCompressed = true
		req.Body = io.NopCloser(bytes.NewReader(compressedBytes))
	}

	if isCompressed {
		// Either append to the header if it already exists, else set it
		if len(req.Header["Content-Encoding"]) != 0 {
			req.Header["Content-Encoding"][0] += ",gzip"
		} else {
			req.Header.Set("Content-Encoding", "gzip")
		}
	}

	return next.HandleBuild(ctx, in)
}

func compress(input io.Reader) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer, %v", err)
	}

	inBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed read payload to compress, %v", err)
	}

	if _, err = w.Write(inBytes); err != nil {
		return nil, fmt.Errorf("failed to write payload to be compressed, %v", err)
	}
	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("failed to flush payload being compressed, %v", err)
	}

	return b.Bytes(), nil
}
