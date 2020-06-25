package aws_test

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/external"
)

func ExampleBuildableHTTPClient_WithTransportOptions_connectionPool() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	// The SDK's HTTPClient implementation implements a WithTransportOptions
	// method for getting an HTTP client with custom transport options.
	type httpClientTransportOptions interface {
		WithTransportOptions(...func(*http.Transport)) aws.HTTPClient
	}

	// Unless cfg.HTTPClient is set to another custom implementation by the
	// application the SDK will use aws.BuildableHTTPClient as the
	// implementation for cfg.HTTPClient.
	client, ok := cfg.HTTPClient.(httpClientTransportOptions)
	if !ok {
		log.Fatalf("expected http client to be SDK's default but wasn't, %T", cfg.HTTPClient)
	}

	// Get a client with custom connection pooling options. The client's
	// options are immutable, and return copies of the client when options are
	// applied.
	cfg.HTTPClient = client.WithTransportOptions(func(tr *http.Transport) {
		tr.MaxIdleConnsPerHost = 150
		// Experiment with 2 * MaxIdleConnsPerHost * number of services, and
		// regions used. Need to balance burst concurrency, with max open
		// connections. May need to adjust how long idle connections are keep
		// around for with IdleConnTimeout.
		tr.MaxIdleConns = tr.MaxIdleConnsPerHost * 2 * 1
	})

	fmt.Printf("Have client: %T, MaxIdleConnsPerHost: %v\n",
		cfg.HTTPClient,
		cfg.HTTPClient.(*aws.BuildableHTTPClient).GetTransport().MaxIdleConnsPerHost,
	)

	// Create service API clients with cfg value to use the custom Transport options.

	// Output: Have client: *aws.BuildableHTTPClient, MaxIdleConnsPerHost: 150
}

func ExampleBuildableHTTPClient_WithTransportOptions_responseTimeouts() {
	// Create a new client by calling the constructor.
	// Add custom *http.Transport configuration to the client. The
	// modifications will be on the new client returned.
	client := aws.NewBuildableHTTPClient().
		WithTransportOptions(func(tr *http.Transport) {
			// Only wait 10 seconds for the full response headers to be read.
			tr.ResponseHeaderTimeout = 10 * time.Second
		})

	// Set the configured HTTP client to the SDK's Config, and use the
	// Config to create API clients with.
	cfg.HTTPClient = client
}

func ExampleBuildableHTTPClient_WithDialerOptions() {
	// Create a new client by calling the constructor.
	client := aws.NewBuildableHTTPClient().
		WithDialerOptions(func(d *net.Dialer) {
			// Set the network (e.g. TCP) dial timeout to 10 seconds.
			d.Timeout = 10 * time.Second
		})

	// Set the configured HTTP client to the SDK's Config, and use the
	// Config to create API clients with.
	cfg.HTTPClient = client
}

var cfg aws.Config
