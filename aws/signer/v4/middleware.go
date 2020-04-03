package v4

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

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
	}
	return fmt.Sprintf("%s: %s", msg, e.cause)
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
	}
	return fmt.Sprintf("%s: %s", msg, e.cause)
}

// Unwrap returns the underlying error cause
func (e *SigningError) Unwrap() error {
	return e.err
}

// UnsignedPayloadMiddleware sets the SigV4 request payload hash to unsigned
type UnsignedPayloadMiddleware struct{}

// ID returns the UnsignedPayloadMiddleware identifier
func (m *UnsignedPayloadMiddleware) ID() string {
	return "SigV4 unsigned payload middleware"
}

// HandleFinalize sets the payload hash to be an unsigned payload
func (m *UnsignedPayloadMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
	ctx = SetPayloadHash(ctx, v4Internal.UnsignedPayload)
	return next.HandleFinalize(ctx, in)
}

// ComputePayloadHashMiddleware computes sha256 payload hash to sign
type ComputePayloadHashMiddleware struct{}

// ID is the middleware name
func (m *ComputePayloadHashMiddleware) ID() string {
	return "SignV4 Payload Hash Middleware"
}

// HandleFinalize compute the payload hash for the request payload
func (m *ComputePayloadHashMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
	req, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &HashComputationError{cause: fmt.Sprintf("unexpected request middleware type %T", in.Request)}
	}

	body, ok := req.Stream.(io.ReadSeeker)
	if !ok {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &HashComputationError{cause: errPayloadNotSeekable}
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &HashComputationError{err: err}
	}

	hash := sha256.New()
	_, err = io.Copy(hash, body)
	if err != nil {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &HashComputationError{err: err}
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &HashComputationError{err: err}
	}

	ctx = SetPayloadHash(ctx, hex.EncodeToString(hash.Sum(nil)))

	return next.HandleFinalize(ctx, in)
}

// SignHTTPRequestMiddleware is a `FinalizeMiddleware` implementation for SigV4 HTTP Signing
type SignHTTPRequestMiddleware struct {
	Signer HTTPSigner
}

// NewSignHTTPRequestMiddleware constructs a SignHTTPRequestMiddleware using the given Signer for signing requests
func NewSignHTTPRequestMiddleware(signer HTTPSigner) *SignHTTPRequestMiddleware {
	return &SignHTTPRequestMiddleware{Signer: signer}
}

// ID is the SignHTTPRequestMiddleware identifier
func (s *SignHTTPRequestMiddleware) ID() string {
	return "SigV4 HTTP Signer"
}

// HandleFinalize will take the provided input and sign the request using the SigV4 authentication scheme
func (s *SignHTTPRequestMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, metadata middleware.Metadata, err error) {
	req, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &SigningError{cause: fmt.Sprintf("unexpected request middleware type %T", in.Request)}
	}

	signingMetadata := GetSigningMetadata(ctx)
	payloadHash := GetPayloadHash(ctx)
	if len(payloadHash) == 0 {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &SigningError{cause: "computed payload hash missing from context"}
	}

	err = s.Signer.SignHTTP(ctx, req.Request, payloadHash, signingMetadata.SigningName, signingMetadata.SigningRegion, sdk.NowTime())
	if err != nil {
		return middleware.FinalizeOutput{}, middleware.NewMetadata(), &SigningError{err: err}
	}

	return next.HandleFinalize(ctx, in)
}

// SigningMetadata contains the signing name and signing region to be used when signing
// with SigV4 authentication scheme.
type SigningMetadata struct {
	SigningName   string
	SigningRegion string
}

type signingMetadataKey struct{}

// GetSigningMetadata retrieves the SigningMetadata from context. If there is no SigningMetadata attached to the context
// an zero-value SigningMetadata will be returned.
func GetSigningMetadata(ctx context.Context) SigningMetadata {
	value := ctx.Value(signingMetadataKey{})

	sm, ok := value.(SigningMetadata)
	if !ok {
		return SigningMetadata{}
	}

	return sm
}

// SetSigningMetadata adds the provided metadata to the context
func SetSigningMetadata(ctx context.Context, metadata SigningMetadata) context.Context {
	ctx = context.WithValue(ctx, signingMetadataKey{}, metadata)
	return ctx
}

type payloadHashKey struct{}

// GetPayloadHash retrieves the payload hash to use for signing
func GetPayloadHash(ctx context.Context) string {
	payloadHash, ok := ctx.Value(payloadHashKey{}).(string)
	if !ok {
		return ""
	}
	return payloadHash
}

// SetPayloadHash sets the payload hash to be used for signing the request
func SetPayloadHash(ctx context.Context, hash string) context.Context {
	ctx = context.WithValue(ctx, payloadHashKey{}, hash)
	return ctx
}
