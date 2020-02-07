package stscreds

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/stsiface"
)

const (
	// ErrCodeWebIdentity will be used as an error code when constructing
	// a new error to be returned during session creation or retrieval.
	ErrCodeWebIdentity = "WebIdentityErr"

	// WebIdentityProviderName is the web identity provider name
	WebIdentityProviderName = "WebIdentityCredentials"
)

// WebIdentityRoleProvider is used to retrieve credentials using
// an OIDC token.
type WebIdentityRoleProvider struct {
	aws.SafeCredentialsProvider

	ExpiryWindow time.Duration

	client stsiface.ClientAPI

	tokenFilePath   string
	roleARN         string
	roleSessionName string
}

// NewWebIdentityRoleProvider will return a new WebIdentityRoleProvider with the
// provided stsiface.STSAPI
func NewWebIdentityRoleProvider(svc stsiface.ClientAPI, roleARN, roleSessionName, path string) *WebIdentityRoleProvider {
	p := &WebIdentityRoleProvider{
		client:          svc,
		tokenFilePath:   path,
		roleARN:         roleARN,
		roleSessionName: roleSessionName,
	}

	p.RetrieveFn = p.retrieveFn

	return p
}

// retrieve will attempt to assume a role from a token which is located at
// 'WebIdentityTokenFilePath' specified destination and if that is empty an
// error will be returned.
func (p *WebIdentityRoleProvider) retrieveFn(ctx context.Context) (aws.Credentials, error) {
	b, err := ioutil.ReadFile(p.tokenFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("unable to read file at %s", p.tokenFilePath)
		return aws.Credentials{}, awserr.New(ErrCodeWebIdentity, errMsg, err)
	}

	sessionName := p.roleSessionName
	if len(sessionName) == 0 {
		// session name is used to uniquely identify a session. This simply
		// uses unix time in nanoseconds to uniquely identify sessions.
		sessionName = strconv.FormatInt(sdk.NowTime().UnixNano(), 10)
	}
	req := p.client.AssumeRoleWithWebIdentityRequest(&sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          &p.roleARN,
		RoleSessionName:  &sessionName,
		WebIdentityToken: aws.String(string(b)),
	})

	// InvalidIdentityToken error is a temporary error that can occur
	// when assuming an Role with a JWT web identity token.
	req.RetryErrorCodes = append(req.RetryErrorCodes, sts.ErrCodeInvalidIdentityTokenException)
	resp, err := req.Send(ctx)
	if err != nil {
		return aws.Credentials{}, awserr.New(ErrCodeWebIdentity, "failed to retrieve credentials", err)
	}

	value := aws.Credentials{
		AccessKeyID:     aws.StringValue(resp.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(resp.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(resp.Credentials.SessionToken),
		Source:          WebIdentityProviderName,
		CanExpire:       true,
		Expires:         resp.Credentials.Expiration.Add(-p.ExpiryWindow),
	}
	return value, nil
}
