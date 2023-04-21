package middleware

import (
	"context"
	"fmt"
	"os"
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

// addRawResponse middleware adds raw response on to the metadata
type addRawResponse struct{}

// ID the identifier for the ClientRequestID
func (m *addRawResponse) ID() string {
	return "AddRawResponseToMetadata"
}

// HandleDeserialize adds raw response on the middleware metadata
func (m addRawResponse) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	metadata.Set(rawResponseKey{}, out.RawResponse)
	return out, metadata, err
}

// AddRawResponseToMetadata adds middleware to the middleware stack that
// store raw response on to the metadata.
func AddRawResponseToMetadata(stack *middleware.Stack) error {
	return stack.Deserialize.Add(&addRawResponse{}, middleware.Before)
}

// GetRawResponse returns raw response set on metadata
func GetRawResponse(metadata middleware.Metadata) interface{} {
	return metadata.Get(rawResponseKey{})
}

// AddRecursionDetection adds recursionDetection to the middleware stack
func AddRecursionDetection(stack *middleware.Stack) error {
	return stack.Build.Add(&RecursionDetection{}, middleware.After)
}

// RecursionDetection detects Lambda environment and sets its X-Ray trace ID to request header if absent
// to avoid recursion invocation in Lambda
type RecursionDetection struct{}

// ID returns the middleware identifier
func (m *RecursionDetection) ID() string {
	return "RecursionDetection"
}

// HandleBuild detects Lambda environment and adds its trace ID to request header if absent
func (m *RecursionDetection) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	_, hasLambdaEnv := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME")
	xAmznTraceID, hasTraceID := os.LookupEnv("_X_AMZN_TRACE_ID")
	value := req.Header.Get("X-Amzn-Trace-Id")
	// only set the X-Amzn-Trace-Id header when it is not set initially, the
	// current environment is Lambda and the _X_AMZN_TRACE_ID env variable exists
	if value != "" || !hasLambdaEnv || !hasTraceID {
		return next.HandleBuild(ctx, in)
	}

	req.Header.Set("X-Amzn-Trace-Id", percentEncode(xAmznTraceID))
	return next.HandleBuild(ctx, in)
}

func percentEncode(s string) string {
	upperhex := "0123456789ABCDEF"
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEncode(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	required := len(s) + 2*hexCount
	t := make([]byte, required)
	j := 0
	for i := 0; i < len(s); i++ {
		if c := s[i]; shouldEncode(c) {
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		} else {
			t[j] = c
			j++
		}
	}
	return string(t)
}

func shouldEncode(c byte) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '=', ';', ':', '+', '&', '[', ']', '{', '}', '"', '\'', ',':
		return false
	default:
		return true
	}
}
