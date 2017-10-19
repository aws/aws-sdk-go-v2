package mock

import (
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
)

func init() {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// TODO getting a default populated config should be in the "defaults" package
	*Config = defaults.Config()
	Config.Region = aws.String("mock-region")
	Config.EndpointResolver = aws.ResolveWithEndpoint(endpoints.ResolvedEndpoint{
		URL:           server.URL,
		SigningRegion: aws.StringValue(Config.Region),
	})
}

// Config is a mock configuration for a localhost mock server returning 200 status.
var Config = &aws.Config{}

// NewMockClient creates and initializes a client that will connect to the
// mock server
func NewMockClient(cfgs ...*aws.Config) *aws.Client {
	c := Config.Copy(cfgs...)

	endpoint, _ := c.EndpointResolver.EndpointFor("Mock", aws.StringValue(c.Region))
	svc := aws.NewClient(
		*c,
		aws.ClientInfo{
			ServiceName:   "Mock",
			SigningRegion: endpoint.SigningRegion,
			Endpoint:      endpoint.URL,
			APIVersion:    "2015-12-08",
			JSONVersion:   "1.1",
			TargetPrefix:  "MockServer",
		},
		c.Handlers,
	)

	return svc
}
