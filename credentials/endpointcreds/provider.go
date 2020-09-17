// Package endpointcreds provides support for retrieving credentials from an
// arbitrary HTTP endpoint.
//
// The credentials endpoint Provider can receive both static and refreshable
// credentials that will expire. Credentials are static when an "Expiration"
// value is not provided in the endpoint's response.
//
// Static credentials will never expire once they have been retrieved. The format
// of the static credentials response:
//    {
//        "AccessKeyId" : "MUA...",
//        "SecretAccessKey" : "/7PC5om....",
//    }
//
// Refreshable credentials will expire within the "ExpiryWindow" of the Expiration
// value in the response. The format of the refreshable credentials response:
//    {
//        "AccessKeyId" : "MUA...",
//        "SecretAccessKey" : "/7PC5om....",
//        "Token" : "AQoDY....=",
//        "Expiration" : "2016-02-25T06:03:31Z"
//    }
//
// Errors should be returned in the following format and only returned with 400
// or 500 HTTP status codes.
//    {
//        "code": "ErrorCode",
//        "message": "Helpful error message."
//    }
package endpointcreds

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds/internal/client"
)

// ProviderName is the name of the credentials provider.
const ProviderName = `CredentialsEndpointProvider`

type getCredentials interface {
	GetCredentials(ctx context.Context, optFns ...func(*client.Options)) (*client.GetCredentialsOutput, error)
}

// Provider satisfies the aws.CredentialsProvider interface, and is a client to
// retrieve credentials from an arbitrary endpoint.
type Provider struct {
	// The AWS Client to make HTTP requests to the endpoint with. The endpoint
	// the request will be made to is provided by the aws.Config's
	// EndpointResolver.
	client getCredentials

	options ProviderOptions
}

// ProviderOptions is structure of configurable options for Provider
type ProviderOptions struct {
	// ExpiryWindow will allow the credentials to trigger refreshing prior to
	// the credentials actually expiring. This is beneficial so race conditions
	// with expiring credentials do not cause request to fail unexpectedly
	// due to ExpiredTokenException exceptions.
	//
	// So a ExpiryWindow of 10s would cause calls to IsExpired() to return true
	// 10 seconds before the credentials are actually expired.
	//
	// If ExpiryWindow is 0 or less it will be ignored.
	ExpiryWindow time.Duration

	// Endpoint to retrieve credentials from
	Endpoint string

	// Optional authorization token value if set will be used as the value of
	// the Authorization header of the endpoint credential request.
	AuthorizationToken string
}

// New returns a credentials Provider for retrieving AWS credentials
// from arbitrary endpoint.
func New(cfg aws.Config, endpoint string, options ...func(*ProviderOptions)) *Provider {
	p := &Provider{
		client: client.NewFromConfig(cfg),
		options: ProviderOptions{
			Endpoint: endpoint,
		},
	}

	for _, option := range options {
		option(&p.options)
	}

	return p
}

// Retrieve will attempt to request the credentials from the endpoint the Provider
// was configured for. And error will be returned if the retrieval fails.
func (p *Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	resp, err := p.getCredentials(ctx)
	if err != nil {
		return aws.Credentials{}, fmt.Errorf("failed to load credentials, %w", err)
	}

	creds := aws.Credentials{
		AccessKeyID:     resp.AccessKeyID,
		SecretAccessKey: resp.SecretAccessKey,
		SessionToken:    resp.Token,
		Source:          ProviderName,
	}

	if resp.Expiration != nil {
		creds.CanExpire = true
		creds.Expires = resp.Expiration.Add(-p.options.ExpiryWindow)
	}

	return creds, nil
}

func (p *Provider) getCredentials(ctx context.Context) (*client.GetCredentialsOutput, error) {
	return p.client.GetCredentials(ctx, func(options *client.Options) {
		options.Endpoint = p.options.Endpoint
		options.AuthorizationToken = p.options.AuthorizationToken
	})
}
