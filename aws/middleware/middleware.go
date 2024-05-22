package middleware

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/rand"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/logging"
	"github.com/aws/smithy-go/middleware"
	smithyrand "github.com/aws/smithy-go/rand"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// ClientRequestID is a Smithy BuildMiddleware that will generate a unique ID for logical API operation
// invocation.
type ClientRequestID struct{}

// ID the identifier for the ClientRequestID
func (r *ClientRequestID) ID() string {
	return "ClientRequestID"
}

// HandleBuild attaches a unique operation invocation id for the operation to the request
func (r ClientRequestID) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", req)
	}

	invocationID, err := smithyrand.NewUUID(rand.Reader).GetUUID()
	if err != nil {
		return out, metadata, err
	}

	const invocationIDHeader = "Amz-Sdk-Invocation-Id"
	req.Header[invocationIDHeader] = append(req.Header[invocationIDHeader][:0], invocationID)

	return next.HandleBuild(ctx, in)
}

// RecordResponseTiming records the response timing for the SDK client requests.
type RecordResponseTiming struct{}

// ID is the middleware identifier
func (a *RecordResponseTiming) ID() string {
	return "RecordResponseTiming"
}

// HandleDeserialize calculates response metadata and clock skew
func (a RecordResponseTiming) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	responseAt := sdk.NowTime()
	setResponseAt(&metadata, responseAt)

	var serverTime time.Time

	switch resp := out.RawResponse.(type) {
	case *smithyhttp.Response:
		respDateHeader := resp.Header.Get("Date")
		if len(respDateHeader) == 0 {
			break
		}
		var parseErr error
		serverTime, parseErr = smithyhttp.ParseTime(respDateHeader)
		if parseErr != nil {
			logger := middleware.GetLogger(ctx)
			logger.Logf(logging.Warn, "failed to parse response Date header value, got %v",
				parseErr.Error())
			break
		}
		setServerTime(&metadata, serverTime)
	}

	if !serverTime.IsZero() {
		attemptSkew := serverTime.Sub(responseAt)
		setAttemptSkew(&metadata, attemptSkew)
	}

	return out, metadata, err
}

type responseAtKey struct{}

// GetResponseAt returns the time response was received at.
func GetResponseAt(metadata middleware.Metadata) (v time.Time, ok bool) {
	v, ok = metadata.Get(responseAtKey{}).(time.Time)
	return v, ok
}

// setResponseAt sets the response time on the metadata.
func setResponseAt(metadata *middleware.Metadata, v time.Time) {
	metadata.Set(responseAtKey{}, v)
}

type serverTimeKey struct{}

// GetServerTime returns the server time for response.
func GetServerTime(metadata middleware.Metadata) (v time.Time, ok bool) {
	v, ok = metadata.Get(serverTimeKey{}).(time.Time)
	return v, ok
}

// setServerTime sets the server time on the metadata.
func setServerTime(metadata *middleware.Metadata, v time.Time) {
	metadata.Set(serverTimeKey{}, v)
}

type attemptSkewKey struct{}

// GetAttemptSkew returns Attempt clock skew for response from metadata.
func GetAttemptSkew(metadata middleware.Metadata) (v time.Duration, ok bool) {
	v, ok = metadata.Get(attemptSkewKey{}).(time.Duration)
	return v, ok
}

// setAttemptSkew sets the attempt clock skew on the metadata.
func setAttemptSkew(metadata *middleware.Metadata, v time.Duration) {
	metadata.Set(attemptSkewKey{}, v)
}

type clockSkew struct{}

// SetAttemptSkewContext sets the clock skew value on the context
func SetAttemptSkewContext(ctx context.Context, v time.Duration) context.Context {
	return middleware.WithStackValue(ctx, clockSkew{}, v)
}

// GetAttemptSkewContext gets the clock skew value from the context
func GetAttemptSkewContext(ctx context.Context) time.Duration {
	x, _ := middleware.GetStackValue(ctx, clockSkew{}).(time.Duration)
	return x
}

// AddClientRequestIDMiddleware adds ClientRequestID to the middleware stack
func AddClientRequestIDMiddleware(stack *middleware.Stack) error {
	return stack.Build.Add(&ClientRequestID{}, middleware.After)
}

// AddRecordResponseTiming adds RecordResponseTiming middleware to the
// middleware stack.
func AddRecordResponseTiming(stack *middleware.Stack) error {
	return stack.Deserialize.Add(&RecordResponseTiming{}, middleware.After)
}

// rawResponseKey is the accessor key used to store and access the
// raw response within the response metadata.
type rawResponseKey struct{}

// AddRawResponse middleware adds raw response on to the metadata
type AddRawResponse struct{}

// ID the identifier for the ClientRequestID
func (m *AddRawResponse) ID() string {
	return "AddRawResponseToMetadata"
}

// HandleDeserialize adds raw response on the middleware metadata
func (m AddRawResponse) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	metadata.Set(rawResponseKey{}, out.RawResponse)
	return out, metadata, err
}

// AddRawResponseToMetadata adds middleware to the middleware stack that
// store raw response on to the metadata.
func AddRawResponseToMetadata(stack *middleware.Stack) error {
	return stack.Deserialize.Add(&AddRawResponse{}, middleware.Before)
}

// GetRawResponse returns raw response set on metadata
func GetRawResponse(metadata middleware.Metadata) interface{} {
	return metadata.Get(rawResponseKey{})
}

// AddTimeOffsetBuildMiddleware sets a value representing clock skew on the request context.
// This can be read by other operations (such as signing) to correct the date value they send
// on the request
type AddTimeOffsetBuildMiddleware struct {
	Offset *atomic.Int64
}

// ID the identifier for AddTimeOffsetBuildMiddleware
func (m *AddTimeOffsetBuildMiddleware) ID() string { return "AddTimeOffsetMiddleware" }

// HandleBuild sets a value for attemptSkew on the request context if one is set on the client.
func (m AddTimeOffsetBuildMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	if m.Offset != nil {
		offset := time.Duration(m.Offset.Load())
		ctx = SetAttemptSkewContext(ctx, offset)
	}
	return next.HandleBuild(ctx, in)
}

// AddTimeOffsetDeserializeMiddleware sets the clock skew on the client if it's present on the context
// at the end of the request
type AddTimeOffsetDeserializeMiddleware struct {
	Offset *atomic.Int64
}

// ID the identifier for AddTimeOffsetDeserializeMiddleware
func (m *AddTimeOffsetDeserializeMiddleware) ID() string { return "AddTimeOffsetDeserializeMiddleware" }

// HandleDeserialize gets the clock skew context from the context, and if set, sets it on the pointer
// held by AddTimeOffsetDeserializeMiddleware
func (m *AddTimeOffsetDeserializeMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	v := GetAttemptSkewContext(ctx)
	if v != 0 {
		m.Offset.Store(v.Nanoseconds())
	}
	return next.HandleDeserialize(ctx, in)
}
