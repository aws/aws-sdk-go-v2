package v4

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/client"
	v4Internal "github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

const errPayloadNotSeekable = "request payload is not seekable"

// HashComputationError indicates an error occurred while computing the signing hash
type HashComputationError struct {
	cause string
	err   error
}

// Error is the error message
func (e *HashComputationError) Error() string {
	msg := "failed to compute payload hash"
	if len(e.cause) == 0 {
		return msg
	} else {
		return fmt.Sprintf("%s: %s", msg, e.cause)
	}
}

// Unwrap returns the underlying error if one is set
func (e *HashComputationError) Unwrap() error {
	return e.err
}

// SigningError indicates an error condition occurred while performing SigV4 signing
type SigningError struct {
	cause string
	err   error
}

func (e *SigningError) Error() string {
	msg := "failed to sign request"
	if len(e.cause) == 0 {
		return msg
	} else {
		return fmt.Sprintf("%s: %s", msg, e.cause)
	}
}

// Unwrap returns the underlying error cause
func (e *SigningError) Unwrap() error {
	return e.err
}

type payloadHashKey struct{}

// UnsignedPayloadMiddleware sets the SigV4 request payload hash to unsigned
type UnsignedPayloadMiddleware struct{}

func (m *UnsignedPayloadMiddleware) Name() string {
	return "SigV4 unsigned payload middleware"
}

func (m *UnsignedPayloadMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	ctx = context.WithValue(ctx, payloadHashKey{}, v4Internal.UnsignedPayload)

	return next.HandleFinalize(ctx, in)
}

// ComputePayloadHashMiddleware computes sha256 payload hash to sign
type ComputePayloadHashMiddleware struct{}

// Name is the middleware name
func (m *ComputePayloadHashMiddleware) Name() string {
	return "SignV4 Payload Hash Middleware"
}

// HandleFinalize compute the payload has fof the request payload for SigV4 authentication
func (m *ComputePayloadHashMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	req := in.Request.(*smithyHTTP.Request)

	body, ok := req.Stream.(io.ReadSeeker)
	if !ok {
		return middleware.FinalizeOutput{}, &HashComputationError{cause: errPayloadNotSeekable}
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return middleware.FinalizeOutput{}, &HashComputationError{err: err}
	}

	hash := sha256.New()
	_, err = io.Copy(hash, body)
	if err != nil {
		return middleware.FinalizeOutput{}, &HashComputationError{err: err}
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return middleware.FinalizeOutput{}, &HashComputationError{err: err}
	}

	ctx = context.WithValue(ctx, payloadHashKey{}, hex.EncodeToString(hash.Sum(nil)))

	return next.HandleFinalize(ctx, in)
}

// SignHTTPRequestMiddleware is a `FinalizeMiddleware` implementation for SigV4 HTTP Signing
type SignHTTPRequestMiddleware struct {
	Signer interface {
		SignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time)
	}
}

func NewSignHTTPRequestMiddleware(signer interface {
	SignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time)
}) *SignHTTPRequestMiddleware {
	return &SignHTTPRequestMiddleware{Signer: signer}
}

// Name is the middlware name
func (s *SignHTTPRequestMiddleware) Name() string {
	return "SigV4 HTTP Signer"
}

// HandleFinalize will take the provided input and sign the request using the SigV4 authentication scheme
func (s *SignHTTPRequestMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	req := in.Request.(*smithyHTTP.Request)

	signingMetadata := client.GetSigningMetadata(ctx)
	payloadHash, ok := ctx.Value(payloadHashKey{}).(string)
	if !ok {
		return middleware.FinalizeOutput{}, &SigningError{cause: "computed payload hash missing from context"}
	}

	s.Signer.SignHTTP(ctx, req.HTTPRequest, payloadHash, signingMetadata.SigningName, signingMetadata.SigningRegion, sdk.NowTime())

	return next.HandleFinalize(ctx, in)
}
