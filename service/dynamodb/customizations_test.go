package dynamodb_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/retry"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/internal/sdk"
	"github.com/jviney/aws-sdk-go-v2/service/dynamodb"
)

var db *dynamodb.Client

func TestMain(m *testing.M) {
	cfg := unit.Config()
	cfg.Retryer = retry.NewStandard(func(s *retry.StandardOptions) {
		s.MaxAttempts = 3
	})

	db = dynamodb.New(cfg)
	db.Handlers.Send.Clear() // mock sending

	os.Exit(m.Run())
}

func mockCRCResponse(svc *dynamodb.Client, status int, body, crc string) (req *aws.Request) {
	header := http.Header{}
	header.Set("x-amz-crc32", crc)

	listReq := svc.ListTablesRequest(nil)
	req = listReq.Request
	req.Handlers.Build.RemoveByName("crr.endpointdiscovery")
	req.Handlers.Send.PushBack(func(*aws.Request) {
		req.HTTPResponse = &http.Response{
			ContentLength: int64(len(body)),
			StatusCode:    status,
			Body:          ioutil.NopCloser(bytes.NewReader([]byte(body))),
			Header:        header,
		}
	})
	req.Send()
	return
}

func TestDefaultRetryRules(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = nil

	svc := dynamodb.New(cfg)
	if e, a := 10, svc.Retryer.MaxAttempts(); e != a {
		t.Errorf("expect %d max retries, got %d", e, a)
	}
}

func TestCustomRetryRules(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = retry.NewStandard(func(s *retry.StandardOptions) {
		s.MaxAttempts = 1
	})

	svc := dynamodb.New(cfg)
	if e, a := 1, svc.Retryer.MaxAttempts(); e != a {
		t.Errorf("expect %d max retries, got %d", e, a)
	}
}

func TestCustomRetry_FromConfig(t *testing.T) {
	cfg := unit.Config()
	cfg.Retryer = retry.NewStandard(func(s *retry.StandardOptions) {
		s.MaxAttempts = 9
	})

	svc := dynamodb.New(cfg)

	if e, a := 9, svc.Retryer.MaxAttempts(); e != a {
		t.Errorf("expect %d max retries from custom retryer, got %d", e, a)
	}
}

func TestValidateCRC32NoHeaderSkip(t *testing.T) {
	req := mockCRCResponse(db, 200, "{}", "")
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}
}

func TestValidateCRC32InvalidHeaderSkip(t *testing.T) {
	req := mockCRCResponse(db, 200, "{}", "ABC")
	if req.Error != nil {
		t.Errorf("expect no error, got %v", req.Error)
	}
}

func TestValidateCRC32AlreadyErrorSkip(t *testing.T) {
	req := mockCRCResponse(db, 400, "{}", "1234")
	if req.Error == nil {
		t.Fatalf("expect error, but got none")
	}

	if e, a := (&dynamodb.CRC32CheckFailedError{}), req.Error; errors.Is(a, e) {
		t.Fatalf("expect error not to be %T", e)
	}
}

func TestValidateCRC32IsValid(t *testing.T) {
	req := mockCRCResponse(db, 200, `{"TableNames":["A"]}`, "3090163698")
	if req.Error != nil {
		t.Fatalf("expect no error, got %v", req.Error)
	}

	// CRC check does not affect output parsing
	out := req.Data.(*dynamodb.ListTablesOutput)
	if e, a := "A", out.TableNames[0]; e != a {
		t.Errorf("expect %q table name, got %q", e, a)
	}
}

func TestValidateCRC32DoesNotMatch(t *testing.T) {
	cleanup := sdk.TestingUseNoOpSleep()
	defer cleanup()

	req := mockCRCResponse(db, 200, "{}", "1234")
	if req.Error == nil {
		t.Fatalf("expect error, but got none")
	}
	req.Handlers.Build.RemoveByName("crr.endpointdiscovery")

	if e, a := 2, req.RetryCount; e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}
	if e, a := (&dynamodb.CRC32CheckFailedError{}), req.Error; !errors.Is(a, e) {
		t.Fatalf("expect %T error, got %T", e, a)
	}
}

func TestValidateCRC32DoesNotMatchNoComputeChecksum(t *testing.T) {
	svc := new(dynamodb.Client)
	*svc = *db
	svc.DisableComputeChecksums = true

	req := mockCRCResponse(svc, 200, `{"TableNames":["A"]}`, "1234")
	if req.Error != nil {
		t.Fatalf("expect no error, got %v", req.Error)
	}

	if e, a := 0, req.RetryCount; e != a {
		t.Errorf("expect %d retry count, got %d", e, a)
	}

	// CRC check disabled. Does not affect output parsing
	out := req.Data.(*dynamodb.ListTablesOutput)
	if e, a := "A", out.TableNames[0]; e != a {
		t.Errorf("expect %q table name, got %q", e, a)
	}
}
