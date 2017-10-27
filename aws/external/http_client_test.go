package external

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

func TestAddHTTPClientCABundle(t *testing.T) {
	certFile, keyFile, caFile, err := awstesting.CreateTLSBundleFiles()
	if err != nil {
		t.Fatalf("failed to create cert temp files, %v", err)
	}
	defer func() {
		awstesting.CleanupTLSBundleFiles(certFile, keyFile, caFile)
	}()

	serverAddr, err := awstesting.CreateTLSServer(certFile, keyFile, nil)
	if err != nil {
		t.Fatalf("failed to start TLS server, %v", err)
	}

	caPEM, err := ioutil.ReadFile(caFile)
	if err != nil {
		t.Fatalf("failed to read CA file, %v", err)
	}

	client := &http.Client{}
	err = addHTTPClientCABundle(client, caPEM)
	if err != nil {
		t.Fatalf("failed to add CA PEM to HTTP client, %v", err)
	}

	req, _ := http.NewRequest("GET", serverAddr, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request to TLS server, %v", err)
	}
	resp.Body.Close()

	if e, a := http.StatusOK, resp.StatusCode; e != a {
		t.Errorf("expect %v status, got %v", e, a)
	}
}
