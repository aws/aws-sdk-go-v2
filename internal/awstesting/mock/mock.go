package mock

import (
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

func init() {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	config = defaults.Config()
	config.Region = "mock-region"
	config.EndpointResolver = aws.ResolveWithEndpoint(aws.Endpoint{
		URL:           server.URL,
		SigningRegion: config.Region,
	})
}

// Config is a mock configuration for a localhost mock server returning 200 status.
var config aws.Config

// Config returns a copy of the mock configuration for tests.
func Config() aws.Config { return config.Copy() }

// NewMockClient creates and initializes a client that will connect to the
// mock server
func NewMockClient(cfg aws.Config) *aws.Client {
	return aws.NewClient(
		cfg,
		aws.Metadata{
			ServiceName:   "Mock",
			SigningRegion: cfg.Region,
			APIVersion:    "2015-12-08",
			JSONVersion:   "1.1",
			TargetPrefix:  "MockServer",
		},
	)
}
