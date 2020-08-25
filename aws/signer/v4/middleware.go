package v4

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	v4Internal "github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

// HashComputationError indicates an error occurred while computing the signing hash
type HashComputationError struct {
	Err error
}

// Error is the error message
func (e *HashComputationError) Error() string {
	return fmt.Sprintf("failed to compute payload hash: %v", e.Err)
}

// Unwrap returns the underlying error if one is set
func (e *HashComputationError) Unwrap() error {
	return e.Err
}

// SigningError indicates an error condition occurred while performing SigV4 signing
type SigningError struct {
	Err error
}

func (e *SigningError) Error() string {
	return fmt.Sprintf("failed to sign request: %v", e.Err)
}

// Unwrap returns the underlying error cause
func (e *SigningError) Unwrap() error {
	return e.Err
}

// unsignedPayloadMiddleware sets the SigV4 request payload hash to unsigned.
//
// Will not set the Unsigned Payload magic SHA value, if a SHA has already been
// stored in the context. (e.g. application pre-computed SHA256 before making
// API call).
//
// This middleware does not check the X-Amz-Content-Sha256 header, if that
// header is serialized a middleware must translate it into the context.
type unsignedPayloadMiddleware struct{}

// AddUnsignedPayloadMiddleware adds unsignedPayloadMiddleware to the operation
// middleware stack
func AddUnsignedPayloadMiddleware(stack *middleware.Stack) {
	stack.Build.Add(&unsignedPayloadMiddleware{}, middleware.After)
}

// ID returns the unsignedPayloadMiddleware identifier
func (m *unsignedPayloadMiddleware) ID() string {
	return "SigV4UnsignedPayloadMiddleware"
}

// HandleBuild sets the payload hash to be an unsigned payload
func (m *unsignedPayloadMiddleware) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	// This should not compute the content SHA256 if the value is already
	// known. (e.g. application pre-computed SHA256 before making API call).
	// Does not have any tight coupling to the X-Amz-Content-Sha256 header, if
	// that header is provided a middleware must translate it into the context.
	contentSHA := GetPayloadHash(ctx)
	if len(contentSHA) == 0 {
		contentSHA = v4Internal.UnsignedPayload
	}

	ctx = SetPayloadHash(ctx, contentSHA)
	return next.HandleBuild(ctx, in)
}

// computePayloadSHA256Middleware computes SHA256 payload hash to sign.
//
// Will not set the Unsigned Payload magic SHA value, if a SHA has already been
// stored in the context. (e.g. application pre-computed SHA256 before making
// API call).
//
// This middleware does not check the X-Amz-Content-Sha256 header, if that
// header is serialized a middleware must translate it into the context.
type computePayloadSHA256Middleware struct{}

// AddComputePayloadSHA256Middleware adds computePayloadSHA256Middleware to the
// operation middleware stack
func AddComputePayloadSHA256Middleware(stack *middleware.Stack) {
	stack.Build.Add(&computePayloadSHA256Middleware{}, middleware.After)
}

// ID is the middleware name
func (m *computePayloadSHA256Middleware) ID() string {
	return "ComputePayloadSHA256Middleware"
}

// HandleBuild compute the payload hash for the request payload
func (m *computePayloadSHA256Middleware) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &HashComputationError{
			Err: fmt.Errorf("unexpected request middleware type %T", in.Request),
		}
	}

	// This should not compute the content SHA256 if the value is already
	// known. (e.g. application pre-computed SHA256 before making API call)
	// Does not have any tight coupling to the X-Amz-Content-Sha256 header, if
	// that header is provided a middleware must translate it into the context.
	if contentSHA := GetPayloadHash(ctx); len(contentSHA) != 0 {
		return next.HandleBuild(ctx, in)
	}

	hash := sha256.New()
	if stream := req.GetStream(); stream != nil {
		_, err = io.Copy(hash, stream)
		if err != nil {
			return out, metadata, &HashComputationError{
				Err: fmt.Errorf("failed to compute payload hash, %w", err),
			}
		}

		if err := req.RewindStream(); err != nil {
			return out, metadata, &HashComputationError{
				Err: fmt.Errorf("failed to seek body to start, %w", err),
			}
		}
	}

	ctx = SetPayloadHash(ctx, hex.EncodeToString(hash.Sum(nil)))

	return next.HandleBuild(ctx, in)
}

// contentSHA256HeaderMiddleware sets the X-Amz-Content-Sha256 header value to
// the Payload hash stored in the context.
type contentSHA256HeaderMiddleware struct{}

// AddContentSHA256HeaderMiddleware adds ContentSHA256HeaderMiddleware to the
// operation middleware stack
func AddContentSHA256HeaderMiddleware(stack *middleware.Stack) {
	stack.Build.Add(&contentSHA256HeaderMiddleware{}, middleware.After)
}

// ID returns the ContentSHA256HeaderMiddleware identifier
func (m *contentSHA256HeaderMiddleware) ID() string {
	return "SigV4ContentSHA256HeaderMiddleware"
}

// HandleBuild sets the X-Amz-Content-Sha256 header value to the Payload hash
// stored in the context.
func (m *contentSHA256HeaderMiddleware) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &HashComputationError{Err: fmt.Errorf("unexpected request middleware type %T", in.Request)}
	}

	req.Header.Set(v4Internal.ContentSHAKey, GetPayloadHash(ctx))

	return next.HandleBuild(ctx, in)
}

// SignHTTPRequestMiddleware is a `FinalizeMiddleware` implementation for SigV4 HTTP Signing
type SignHTTPRequestMiddleware struct {
	credentialsProvider aws.CredentialsProvider
	signer              HTTPSigner
}

// NewSignHTTPRequestMiddleware constructs a SignHTTPRequestMiddleware using the given Signer for signing requests
func NewSignHTTPRequestMiddleware(credentialsProvider aws.CredentialsProvider, signer HTTPSigner) *SignHTTPRequestMiddleware {
	return &SignHTTPRequestMiddleware{credentialsProvider: credentialsProvider, signer: signer}
}

// ID is the SignHTTPRequestMiddleware identifier
func (s *SignHTTPRequestMiddleware) ID() string {
	return "SigV4SignHTTPRequestMiddleware"
}

// HandleFinalize will take the provided input and sign the request using the SigV4 authentication scheme
func (s *SignHTTPRequestMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &SigningError{Err: fmt.Errorf("unexpected request middleware type %T", in.Request)}
	}

	signingName, signingRegion := awsmiddleware.GetSigningName(ctx), awsmiddleware.GetSigningRegion(ctx)
	payloadHash := GetPayloadHash(ctx)
	if len(payloadHash) == 0 {
		return out, metadata, &SigningError{Err: fmt.Errorf("computed payload hash missing from context")}
	}

	credentials, err := s.credentialsProvider.Retrieve(ctx)
	if err != nil {
		return out, metadata, &SigningError{Err: fmt.Errorf("failed to retrieve credentials: %w", err)}
	}

	err = s.signer.SignHTTP(ctx, credentials, req.Request, payloadHash, signingName, signingRegion, sdk.NowTime())
	if err != nil {
		return out, metadata, &SigningError{Err: fmt.Errorf("failed to sign http request, %w", err)}
	}

	return next.HandleFinalize(ctx, in)
}

type payloadHashKey struct{}

// GetPayloadHash retrieves the payload hash to use for signing
func GetPayloadHash(ctx context.Context) (v string) {
	v, _ = ctx.Value(payloadHashKey{}).(string)
	return v
}

// SetPayloadHash sets the payload hash to be used for signing the request
func SetPayloadHash(ctx context.Context, hash string) context.Context {
	ctx = context.WithValue(ctx, payloadHashKey{}, hash)
	return ctx
}
