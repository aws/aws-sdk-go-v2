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
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

// ProviderName is the name of the credentials provider.
const ProviderName = `CredentialsEndpointProvider`

// Provider satisfies the aws.CredentialsProvider interface, and is a client to
// retrieve credentials from an arbitrary endpoint.
type Provider struct {
	aws.Expiry

	// The AWS Client to make HTTP requests to the endpoint with. The endpoint
	// the request will be made to is provided by the aws.Config's
	// EndpointResolver.
	Client *aws.Client

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

	staticCreds bool
}

// New returns a credentials Provider for retrieving AWS credentials
// from arbitrary endpoint.
func New(cfg aws.Config) *Provider {
	// TODO don't ignore error, move endpoint resolver to client.
	endpoint, _ := cfg.EndpointResolver.EndpointFor("CredentialsEndpoint", "")

	p := &Provider{
		Client: aws.NewClient(
			cfg,
			aws.ClientInfo{
				ServiceName: "CredentialsEndpoint",
				Endpoint:    endpoint.URL,
			},
			cfg.Handlers,
		),
	}

	p.Client.Handlers.Unmarshal.PushBack(unmarshalHandler)
	p.Client.Handlers.UnmarshalError.PushBack(unmarshalError)
	p.Client.Handlers.Validate.Clear()
	p.Client.Handlers.Validate.PushBack(validateEndpointHandler)

	return p
}

// IsExpired returns true if the credentials retrieved are expired, or not yet
// retrieved.
func (p *Provider) IsExpired() bool {
	if p.staticCreds {
		return false
	}
	return p.Expiry.IsExpired()
}

// Retrieve will attempt to request the credentials from the endpoint the Provider
// was configured for. And error will be returned if the retrieval fails.
func (p *Provider) Retrieve() (aws.Credentials, error) {
	resp, err := p.getCredentials()
	if err != nil {
		return aws.Credentials{Source: ProviderName},
			awserr.New("CredentialsEndpointError", "failed to load credentials", err)
	}

	if resp.Expiration != nil {
		p.SetExpiration(*resp.Expiration, p.ExpiryWindow)
	} else {
		p.staticCreds = true
	}

	return aws.Credentials{
		AccessKeyID:     resp.AccessKeyID,
		SecretAccessKey: resp.SecretAccessKey,
		SessionToken:    resp.Token,
		Source:    ProviderName,
	}, nil
}

type getCredentialsOutput struct {
	Expiration      *time.Time
	AccessKeyID     string
	SecretAccessKey string
	Token           string
}

type errorOutput struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *Provider) getCredentials() (*getCredentialsOutput, error) {
	op := &aws.Operation{
		Name:       "GetCredentials",
		HTTPMethod: "GET",
	}

	out := &getCredentialsOutput{}
	req := p.Client.NewRequest(op, nil, out)
	req.HTTPRequest.Header.Set("Accept", "application/json")

	return out, req.Send()
}

func validateEndpointHandler(r *aws.Request) {
	if len(r.ClientInfo.Endpoint) == 0 {
		r.Error = aws.ErrMissingEndpoint
	}
}

func unmarshalHandler(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()

	out := r.Data.(*getCredentialsOutput)
	if err := json.NewDecoder(r.HTTPResponse.Body).Decode(&out); err != nil {
		r.Error = awserr.New("SerializationError",
			"failed to decode endpoint credentials",
			err,
		)
	}
}

func unmarshalError(r *aws.Request) {
	defer r.HTTPResponse.Body.Close()

	var errOut errorOutput
	if err := json.NewDecoder(r.HTTPResponse.Body).Decode(&errOut); err != nil {
		r.Error = awserr.New("SerializationError",
			"failed to decode endpoint credentials",
			err,
		)
	}

	// Response body format is not consistent between metadata endpoints.
	// Grab the error message as a string and include that as the source error
	r.Error = awserr.New(errOut.Code, errOut.Message, nil)
}
