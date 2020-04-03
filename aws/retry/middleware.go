package retry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

// RequestCloner is a function that can take an input request type and clone the request
// for use in a subsequent retry attempt
type RequestCloner func(context.Context, interface{}) interface{}

type retryMetadata struct {
	AttemptNum       int
	AttemptTime      time.Time
	MaxAttempts      int
	AttemptClockSkew time.Duration
}

type retryMetadataKey struct{}

// RetryMiddleware is a Smithy FinalizeMiddleware that handles retry attempts using the provided
// Retryer implementation
type RetryMiddleware struct {
	retryer       Retryer
	requestCloner RequestCloner
}

// NewRetryMiddleware returns a new RetryMiddleware
func NewRetryMiddleware(retryer Retryer, requestCloner RequestCloner) RetryMiddleware {
	return RetryMiddleware{retryer: retryer, requestCloner: requestCloner}
}

func (r RetryMiddleware) ID() string {
	return "Retry Middleware"
}

func (r RetryMiddleware) Name() string {
	return r.ID()
}

func (r RetryMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	var attemptNum, retryCount int
	var attemptClockSkew time.Duration

	maxAttempts := r.retryer.MaxAttempts()

	relRetryToken := r.retryer.GetInitialToken()

	origReq := r.requestCloner(ctx, in.Request)
	origCtx := ctx

	for {
		attemptNum++

		ctx = setRetryMetadata(origCtx, retryMetadata{
			AttemptNum:       attemptNum,
			AttemptTime:      sdk.NowTime(),
			MaxAttempts:      maxAttempts,
			AttemptClockSkew: attemptClockSkew,
		})

		in.Request = r.requestCloner(ctx, origReq)

		out, reqErr := next.HandleFinalize(ctx, in)

		relRetryToken(reqErr)
		if reqErr == nil {
			return out, nil
		}

		retryable := r.retryer.IsErrorRetryable(reqErr)
		if !retryable {
			return out, err
		}

		if maxAttempts > 0 && attemptNum >= maxAttempts {
			err = &aws.MaxAttemptsError{
				Attempt: attemptNum,
				Err:     err,
			}
			return out, err
		}

		relRetryToken, err = r.retryer.GetRetryToken(ctx, reqErr)
		if err != nil {
			return out, err
		}

		retryDelay, err := r.retryer.RetryDelay(attemptNum, reqErr)
		if err != nil {
			return out, err
		}

		if err = sdk.SleepWithContext(ctx, retryDelay); err != nil {
			err = &aws.RequestCanceledError{Err: err}
			return out, err
		}

		// TODO: Finalize what this interface and types look like, types and returns here are not final and exists to strictly model the required behavior
		type responseMeta interface {
			GetResponseMetadata() interface {
				GetResponseAt() time.Time
				GetServerTime() time.Time
			}
		}

		// TODO: Pull this from future smithy Metadata context type
		respContainer, ok := out.Result.(responseMeta)
		if ok {
			metadata := respContainer.GetResponseMetadata()
			responseAt := metadata.GetResponseAt()
			serverTime := metadata.GetServerTime()

			// TODO: This should probably be computed and bubbled up from a deserializer
			if !(responseAt.IsZero() || serverTime.IsZero()) {
				attemptClockSkew = responseAt.Sub(serverTime)
			}
		}

		retryCount++
	}
}

type RetryMetricsHeaderMiddleware struct{}

func (r RetryMetricsHeaderMiddleware) ID() string {
	return "Retry Metrics Header Middleware"
}

func (r RetryMetricsHeaderMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	metadata, ok := getRetryMetadata(ctx)
	if !ok {
		return out, fmt.Errorf("retry metadata value not found on context")
	}

	const retryMetricHeader = "amz-sdk-request"
	var parts []string

	parts = append(parts, fmt.Sprintf("attempt=%d", metadata.AttemptNum))
	if metadata.MaxAttempts != 0 {
		parts = append(parts, fmt.Sprintf("max=%d", metadata.MaxAttempts))
	}

	var ttl time.Time
	if deadline, ok := ctx.Deadline(); ok {
		ttl = deadline
	}

	// Only append the TTL if it can be determined.
	if !ttl.IsZero() && metadata.AttemptClockSkew > 0 {
		const unixTimeFormat = "20060102T150405Z"
		ttl = ttl.Add(metadata.AttemptClockSkew)
		parts = append(parts, fmt.Sprintf("ttl=%s", ttl.Format(unixTimeFormat)))
	}

	switch req := in.Request.(type) {
	case *http.Request:
		req.Header.Set(retryMetricHeader, strings.Join(parts, "; "))
	default:
		return middleware.FinalizeOutput{}, fmt.Errorf("unknown transport type %T", req)
	}

	return next.HandleFinalize(ctx, in)
}

// getRetryMetadata retrieves retryMetadata from the context and a bool indicating if it was set
func getRetryMetadata(ctx context.Context) (retryMetadata, bool) {
	metadata, ok := ctx.Value(retryMetadataKey{}).(retryMetadata)
	return metadata, ok
}

func setRetryMetadata(ctx context.Context, metadata retryMetadata) context.Context {
	return context.WithValue(ctx, retryMetadataKey{}, metadata)
}
