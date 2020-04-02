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

type retryMetadata struct {
	AttemptNum       int
	AttemptTime      time.Time
	MaxAttempts      int
	AttemptClockSkew time.Duration
}

type retryMetadataKey struct{}

type RetryMiddleware struct {
	Retryer Retryer
}

func (r RetryMiddleware) ID() string {
	return "Retry Middleware"
}

func (r RetryMiddleware) Name() string {
	return r.ID()
}

func (r RetryMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	const invocationIDHeader = "amz-sdk-invocation-id"

	req := in.Request.(*http.Request)

	// TODO: Don't see a compelling reason to make this a separate middleware that makes this available via context
	invocationID, err := sdk.UUIDVersion4()
	if err != nil {
		return out, err
	}

	req.Header.Set(invocationIDHeader, invocationID)

	var attemptNum, retryCount int
	var attemptClockSkew time.Duration

	maxAttempts := r.Retryer.MaxAttempts()

	relRetryToken := r.Retryer.GetInitialToken()
	for {
		attemptNum++

		ctx = context.WithValue(ctx, retryMetadataKey{}, retryMetadata{
			AttemptNum:       attemptNum,
			AttemptTime:      sdk.NowTime(),
			MaxAttempts:      maxAttempts,
			AttemptClockSkew: attemptClockSkew,
		})

		out, reqErr := next.HandleFinalize(ctx, in)

		relRetryToken(reqErr)
		if reqErr == nil {
			return out, nil
		}

		retryable := r.Retryer.IsErrorRetryable(reqErr)
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

		retryDelay, err := r.Retryer.RetryDelay(attemptNum, reqErr)
		if err != nil {
			return out, err
		}

		relRetryToken, err = r.Retryer.GetRetryToken(ctx, reqErr)
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

		respContainer, ok := out.Result.(responseMeta)
		if ok {
			metadata := respContainer.GetResponseMetadata()
			responseAt := metadata.GetResponseAt()
			serverTime := metadata.GetServerTime()

			if !(responseAt.IsZero() || serverTime.IsZero()) {
				attemptClockSkew = responseAt.Sub(serverTime)
			}
		}

		in.Request = in.Request.(*http.Request).Clone(ctx)

		retryCount++
	}
}

type RetryMetricsHeaderMiddleware struct {
	ClientTimeout time.Duration
}

func (r RetryMetricsHeaderMiddleware) ID() string {
	return "Retry Metrics Header Middleware"
}

func (r RetryMetricsHeaderMiddleware) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (out middleware.FinalizeOutput, err error) {
	metadata, ok := ctx.Value(retryMetadataKey{}).(retryMetadata)
	if !ok {
		return out, fmt.Errorf("retry metadata value not found on context")
	}

	const retryMetricHeader = "amz-sdk-request"
	var parts []string

	parts = append(parts, fmt.Sprintf("attempt=%d", metadata.AttemptNum))
	if metadata.MaxAttempts != 0 {
		parts = append(parts, fmt.Sprintf("max=%d", metadata.MaxAttempts))
	}

	// TODO: Is there a method that we could know the client timeout? Maybe this is passed via middleware constructor?
	//type timeoutGetter interface {
	//	GetTimeout() time.Duration
	//}
	//
	//var ttl time.Time
	//// Attempt extract the TTL from the timeout on the client.
	//if v, ok := r.Config.HTTPClient.(timeoutGetter); ok {
	//	if t := v.GetTimeout(); t > 0 {
	//		ttl = sdk.NowTime().Add(t)
	//	}
	//}

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

	req := in.Request.(*http.Request)

	req.Header.Set(retryMetricHeader, strings.Join(parts, "; "))

	return next.HandleFinalize(ctx, in)
}
