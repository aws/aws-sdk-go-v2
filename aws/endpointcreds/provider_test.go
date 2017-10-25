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
	"github.com/stretchr/testify/assert"
)

func TestRetrieveRefreshableCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/path/to/endpoint", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "else", r.URL.Query().Get("something"))

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

	svc := endpointcreds.New(cfg)
	creds, err := svc.Retrieve()

	assert.NoError(t, err)

	assert.Equal(t, "AKID", creds.AccessKeyID)
	assert.Equal(t, "SECRET", creds.SecretAccessKey)
	assert.Equal(t, "TOKEN", creds.SessionToken)
	assert.False(t, svc.IsExpired())

	svc.CurrentTime = func() time.Time {
		return time.Now().Add(2 * time.Hour)
	}

	assert.True(t, svc.IsExpired())
}

func TestRetrieveStaticCredentials(t *testing.T) {
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

	client := endpointcreds.New(cfg)
	creds, err := client.Retrieve()

	assert.NoError(t, err)

	assert.Equal(t, "AKID", creds.AccessKeyID)
	assert.Equal(t, "SECRET", creds.SecretAccessKey)
	assert.Empty(t, creds.SessionToken)
	assert.False(t, client.IsExpired())
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

	client := endpointcreds.New(cfg)
	creds, err := client.Retrieve()

	assert.Error(t, err)
	aerr := err.(awserr.Error)

	assert.Equal(t, "CredentialsEndpointError", aerr.Code())
	assert.Equal(t, "failed to load credentials", aerr.Message())

	aerr = aerr.OrigErr().(awserr.Error)
	assert.Equal(t, "Error", aerr.Code())
	assert.Equal(t, "Message", aerr.Message())

	assert.Empty(t, creds.AccessKeyID)
	assert.Empty(t, creds.SecretAccessKey)
	assert.Empty(t, creds.SessionToken)
	assert.True(t, client.IsExpired())
}
