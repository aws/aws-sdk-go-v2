package kinesis

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/defaults"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/internal/sdk"
)

type testReader struct {
	duration time.Duration
}

func (r *testReader) Read(b []byte) (int, error) {
	time.Sleep(r.duration)
	return 0, io.EOF
}

func (r *testReader) Close() error {
	return nil
}

// GetRecords will hang unexpectedly during reads.
// See https://github.com/aws/aws-sdk-go-v2/issues/1141
func TestKinesisGetRecordsCustomization(t *testing.T) {
	restoreSleep := sdk.TestingUseNoOpSleep()
	defer restoreSleep()

	readDuration = time.Millisecond
	attempts := 0

	cfg := unit.Config()
	svc := New(cfg)
	req := svc.GetRecordsRequest(&GetRecordsInput{
		ShardIterator: aws.String("foo"),
	})
	req.Handlers.Send.Clear()
	req.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"X-Amz-Request-Id": []string{"abc123"},
			},
			Body:          &testReader{duration: 10 * time.Second},
			ContentLength: -1,
		}
		r.HTTPResponse.Status = http.StatusText(r.HTTPResponse.StatusCode)
		attempts++
	})
	req.ApplyOptions(aws.WithResponseReadTimeout(time.Second))
	_, err := req.Send(context.Background())
	if err == nil {
		t.Fatalf("expect error, got none")
	}

	if e, a := (*aws.ResponseTimeoutError)(nil), err; !errors.As(a, &e) {
		t.Fatalf("expect %T, error in %v", e, a)
	}

	if e, a := svc.Retryer.MaxAttempts(), attempts; e != a {
		t.Errorf("Expected %v attempts, but received %d", e, a)
	}
}

func TestKinesisGetRecordsNoTimeout(t *testing.T) {
	readDuration = time.Second
	svc := New(unit.Config())
	req := svc.GetRecordsRequest(&GetRecordsInput{
		ShardIterator: aws.String("foo"),
	})
	req.Handlers.Send.Clear()
	req.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"X-Amz-Request-Id": []string{"abc123"},
			},
			Body:          &testReader{duration: time.Duration(0)},
			ContentLength: -1,
		}
		r.HTTPResponse.Status = http.StatusText(r.HTTPResponse.StatusCode)
	})
	req.ApplyOptions(aws.WithResponseReadTimeout(time.Second))
	_, err := req.Send(context.Background())
	if err != nil {
		t.Errorf("Expected no error, but received %v", err)
	}
}

func TestKinesisCustomRetryErrorCodes(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = retry.NewStandard(func(s *retry.StandardOptions) {
		s.MaxAttempts = 2
	})

	svc := New(cfg)
	svc.Handlers.Validate.Clear()

	const jsonErr = `{"__type":%q, "message":"some error message"}`
	var attempts int
	resps := []*http.Response{
		{
			StatusCode: 400,
			Header:     http.Header{},
			Body: ioutil.NopCloser(bytes.NewReader(
				[]byte(fmt.Sprintf(jsonErr, ErrCodeLimitExceededException)),
			)),
		},
		{
			StatusCode: 200,
			Header:     http.Header{},
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		},
	}

	req := svc.GetRecordsRequest(&GetRecordsInput{})
	req.Handlers.Send.Swap(defaults.SendHandler.Name, aws.NamedHandler{
		Name: "custom send handler",
		Fn: func(r *aws.Request) {
			r.HTTPResponse = resps[attempts]
			attempts++
		},
	})

	if _, err := req.Send(context.Background()); err != nil {
		t.Fatalf("expect no error, got %T, %v", err, err)
	}

	if e, a := 2, attempts; e != a {
		t.Errorf("expect %v requests, got %v", e, a)
	}
}
