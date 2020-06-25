package retry_test

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/awserr"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
)

type mockTemporaryError struct{ b bool }

func (m mockTemporaryError) Temporary() bool { return m.b }
func (m mockTemporaryError) Error() string {
	return fmt.Sprintf("mock temporary %t", m.b)
}

type mockTimeoutError struct{ b bool }

func (m mockTimeoutError) Timeout() bool { return m.b }
func (m mockTimeoutError) Error() string {
	return fmt.Sprintf("mock timeout %t", m.b)
}

type mockRetryableError struct{ b bool }

func (m mockRetryableError) RetryableError() bool { return m.b }
func (m mockRetryableError) Error() string {
	return fmt.Sprintf("mock retryable %t", m.b)
}

type mockCanceledError struct{ b bool }

func (m mockCanceledError) CanceledError() bool { return m.b }
func (m mockCanceledError) Error() string {
	return fmt.Sprintf("mock canceled %t", m.b)
}

type mockStatusCodeError struct{ code int }

func (m mockStatusCodeError) StatusCode() int { return m.code }
func (m mockStatusCodeError) Error() string {
	return fmt.Sprintf("status code error, %v", m.code)
}

func TestRetryConnectionErrors(t *testing.T) {
	cases := map[string]struct {
		Err       error
		Retryable aws.Ternary
	}{
		"nested connection reset": {
			Retryable: aws.TrueTernary,
			Err: fmt.Errorf("serialization error, %w",
				fmt.Errorf("connection reset")),
		},
		"top level connection reset": {
			Retryable: aws.TrueTernary,
			Err:       fmt.Errorf("connection reset"),
		},
		"awserr connection reset": {
			Retryable: aws.TrueTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				fmt.Errorf("connection reset")),
		},
		"url.Error connection refused": {
			Retryable: aws.TrueTernary,
			Err: fmt.Errorf("some error, %w", &url.Error{
				Err: fmt.Errorf("connection refused"),
			}),
		},
		"other connection refused": {
			Retryable: aws.UnknownTernary,
			Err:       fmt.Errorf("connection refused"),
		},
		"nil error connection reset": {
			Retryable: aws.UnknownTernary,
		},
		"some other error": {
			Retryable: aws.UnknownTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				fmt.Errorf("something bad")),
		},
		"request send error": {
			Retryable: aws.TrueTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				&aws.RequestSendError{Err: &url.Error{
					Err: fmt.Errorf("another error"),
				}}),
		},
		"temporary error": {
			Retryable: aws.TrueTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				mockTemporaryError{b: true},
			),
		},
		"timeout error": {
			Retryable: aws.TrueTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				mockTimeoutError{b: true},
			),
		},
		"timeout false error": {
			Retryable: aws.UnknownTernary,
			Err: awserr.New(aws.ErrCodeSerialization, "some error",
				mockTimeoutError{b: false},
			),
		},
		"net.OpError dial": {
			Retryable: aws.TrueTernary,
			Err: &net.OpError{
				Op:  "dial",
				Err: mockTimeoutError{b: false},
			},
		},
		"net.OpError nested": {
			Retryable: aws.TrueTernary,
			Err: &net.OpError{
				Op:  "read",
				Err: fmt.Errorf("some error %w", mockTimeoutError{b: true}),
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var r retry.RetryableConnectionError

			retryable := r.IsErrorRetryable(c.Err)
			if e, a := c.Retryable, retryable; e != a {
				t.Errorf("expect %v retryable, got %v", e, a)
			}
		})
	}
}

func TestRetryHTTPStatusCodes(t *testing.T) {
	cases := map[string]struct {
		Err    error
		Expect aws.Ternary
	}{
		"top level": {
			Err:    &mockStatusCodeError{code: 500},
			Expect: aws.TrueTernary,
		},
		"nested": {
			Err:    fmt.Errorf("some error, %w", &mockStatusCodeError{code: 500}),
			Expect: aws.TrueTernary,
		},
		"response error": {
			Err: fmt.Errorf("some error, %w", &aws.HTTPResponseError{
				Response: &http.Response{StatusCode: 502},
			}),
			Expect: aws.TrueTernary,
		},
	}

	r := retry.RetryableHTTPStatusCode{Codes: map[int]struct{}{
		500: struct{}{},
		502: struct{}{},
	}}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestRetryErrorCodes(t *testing.T) {
	cases := map[string]struct {
		Err    error
		Expect aws.Ternary
	}{
		"retryable code": {
			Err: &aws.MaxAttemptsError{
				Err: awserr.New("ErrorCode1", "some error", nil),
			},
			Expect: aws.TrueTernary,
		},
		"not retryable code": {
			Err: &aws.MaxAttemptsError{
				Err: awserr.New("SomeErroCode", "some error", nil),
			},
			Expect: aws.UnknownTernary,
		},
		"other error": {
			Err:    fmt.Errorf("some other error"),
			Expect: aws.UnknownTernary,
		},
	}

	r := retry.RetryableErrorCode{Codes: map[string]struct{}{
		"ErrorCode1": struct{}{},
		"ErrorCode2": struct{}{},
	}}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}

func TestCanceledError(t *testing.T) {
	cases := map[string]struct {
		Err    error
		Expect aws.Ternary
	}{
		"canceled error": {
			Err: fmt.Errorf("some error, %w", &aws.RequestCanceledError{
				Err: fmt.Errorf(":("),
			}),
			Expect: aws.FalseTernary,
		},
		"canceled retryable error": {
			Err: fmt.Errorf("some error, %w", &aws.RequestCanceledError{
				Err: mockRetryableError{b: true},
			}),
			Expect: aws.FalseTernary,
		},
		"not canceled error": {
			Err:    fmt.Errorf("some error, %w", mockCanceledError{b: false}),
			Expect: aws.UnknownTernary,
		},
		"retryable error": {
			Err:    fmt.Errorf("some error, %w", mockRetryableError{b: true}),
			Expect: aws.TrueTernary,
		},
	}

	r := retry.IsErrorRetryables{
		retry.NoRetryCanceledError{},
		retry.RetryableError{},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.Expect, r.IsErrorRetryable(c.Err); e != a {
				t.Errorf("Expect %v retryable, got %v", e, a)
			}
		})
	}
}
