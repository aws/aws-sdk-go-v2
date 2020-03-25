package v4

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/awslabs/smithy-go/middleware"
)

// SignHTTPRequestMiddleware is a `FinalizeMiddleware` implementation for SigV4 HTTP Signing
type SignHTTPRequestMiddleware struct {
	Signer interface {
		Sign(ctx context.Context, r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error)
	}
}

func NewSignHTTPRequestMiddleware(signer interface {
	Sign(ctx context.Context, r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error)
}) *SignHTTPRequestMiddleware {
	return &SignHTTPRequestMiddleware{Signer: signer}
}

func (s *SignHTTPRequestMiddleware) Name() string {
	return "SigV4 HTTP Signer"
}

func (s *SignHTTPRequestMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	req := in.Request.(*http.Request)

	s.Signer.Sign(ctx, req.Body)
}
