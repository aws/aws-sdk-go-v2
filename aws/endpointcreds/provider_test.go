package endpointcreds_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

func TestRetrieveRefreshableCredentials(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e, a := "/path/to/endpoint", r.URL.Path; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := "application/json", r.Header.Get("Accept"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
		if e, a := "else", r.URL.Query().Get("something"); e != a {
			t.Errorf("expect %v, got %v", e, a)
		}

		encoder := json.NewEncoder(w)
		err := encoder.Encode(map[string]interface{}{
			"AccessKeyID":     "AKID",
			"SecretAccessKey": "SECRET",
			"Token":           "TOKEN",
			"Expiration":      time.Now().Add(1 * time.Hour),
		})

		if err != nil {
			fmt.Println("failed to write out creds", err)
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL + "/path/to/endpoint?something=else")

	p := endpointcreds.New(cfg)
	creds, err := p.Retrieve()

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "TOKEN", creds.SessionToken; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if creds.Expired() {
		t.Errorf("expect not expired")
	}

	sdk.NowTime = func() time.Time {
		return time.Now().Add(2 * time.Hour)
	}
	if !creds.Expired() {
		t.Errorf("expect to be expired")
	}
}

func TestRetrieveStaticCredentials(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		err := encoder.Encode(map[string]interface{}{
			"AccessKeyID":     "AKID",
			"SecretAccessKey": "SECRET",
		})

		if err != nil {
			fmt.Println("failed to write out creds", err)
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	p := endpointcreds.New(cfg)
	creds, err := p.Retrieve()

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}

	sdk.NowTime = func() time.Time {
		return time.Date(3000, 12, 16, 1, 30, 37, 0, time.UTC)
	}

	if creds.Expired() {
		t.Errorf("expect not to be expired")
	}
}

func TestFailedRetrieveCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(map[string]interface{}{
			"Code":    "Error",
			"Message": "Message",
		})

		if err != nil {
			fmt.Println("failed to write error", err)
		}
	}))

	cfg := unit.Config()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	p := endpointcreds.New(cfg)
	creds, err := p.Retrieve()

	if err == nil {
		t.Fatalf("expect error, got none")
	}

	aerr := err.(awserr.Error)
	if e, a := "CredentialsEndpointError", aerr.Code(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "failed to load credentials", aerr.Message(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	aerr = aerr.OrigErr().(awserr.Error)
	if e, a := "Error", aerr.Code(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "Message", aerr.Message(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if v := creds.AccessKeyID; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SecretAccessKey; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if v := creds.SessionToken; len(v) != 0 {
		t.Errorf("expect empty, got %v", v)
	}
	if creds.Expired() {
		t.Errorf("expect empty creds not to be expired")
	}
}
