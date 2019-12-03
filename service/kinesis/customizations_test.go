package kinesis

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
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
	readDuration = time.Millisecond
	retryCount := 0

	cfg := unit.Config()
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 4
	})

	svc := New(cfg)
	req := svc.GetRecordsRequest(&types.GetRecordsInput{
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
		retryCount++
	})
	req.ApplyOptions(aws.WithResponseReadTimeout(time.Second))
	_, err := req.Send(context.Background())
	if err == nil {
		t.Errorf("Expected error, but received nil")
	} else if v, ok := err.(awserr.Error); !ok {
		t.Errorf("Expected awserr.Error but received %v", err)
	} else if v.Code() != aws.ErrCodeResponseTimeout {
		t.Errorf("Expected 'RequestTimeout' error, but received %s instead", v.Code())
	}
	if retryCount != 5 {
		t.Errorf("Expected '5' retries, but received %d", retryCount)
	}
}

func TestKinesisGetRecordsNoTimeout(t *testing.T) {
	readDuration = time.Second
	svc := New(unit.Config())
	req := svc.GetRecordsRequest(&types.GetRecordsInput{
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
	cfg.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
	cfg.Retryer = aws.NewDefaultRetryer(func(d *aws.DefaultRetryer) {
		d.NumMaxRetries = 1
	})

	svc := New(cfg)
	svc.Handlers.Validate.Clear()

	const jsonErr = `{"__type":%q, "message":"some error message"}`
	var reqCount int
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

	req := svc.GetRecordsRequest(&types.GetRecordsInput{})
	req.Handlers.Send.Swap(defaults.SendHandler.Name, aws.NamedHandler{
		Name: "custom send handler",
		Fn: func(r *aws.Request) {
			r.HTTPResponse = resps[reqCount]
			reqCount++
		},
	})

	if _, err := req.Send(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 2, reqCount; e != a {
		t.Errorf("expect %v requests, got %v", e, a)
	}
}
