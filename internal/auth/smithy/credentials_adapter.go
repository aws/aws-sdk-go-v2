package smithy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/auth"
)

// CredentialsAdapter adapts aws.Credentials to auth.Identity.
type CredentialsAdapter struct {
	Credentials aws.Credentials
}

var _ auth.Identity = (*CredentialsAdapter)(nil)

// Expiration returns the time of expiration for the credentials.
func (v *CredentialsAdapter) Expiration() time.Time {
	return v.Credentials.Expires
}

// CredentialsProviderAdapter adapts aws.CredentialsProvider to auth.IdentityResolver.
type CredentialsProviderAdapter struct {
	Provider aws.CredentialsProvider
}

var _ (auth.IdentityResolver) = (*CredentialsProviderAdapter)(nil)

// GetIdentity retrieves AWS credentials using the underlying provider.
func (v *CredentialsProviderAdapter) GetIdentity(ctx context.Context, _ smithy.Properties) (
	auth.Identity, error,
) {
	if v.Provider == nil {
		return &CredentialsAdapter{Credentials: aws.Credentials{}}, nil
	}

	creds, err := v.Provider.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}

	return &CredentialsAdapter{Credentials: creds}, nil
}

// InvalidateIdentity marks the cached credentials for refresh if the service rejected them as expired or invalid.
func (v *CredentialsProviderAdapter) InvalidateIdentity(_ context.Context, identity auth.Identity, err error) {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return
	}
	switch apiErr.ErrorCode() {
	case "ExpiredToken", "InvalidToken":
	default:
		return
	}

	cache, ok := v.Provider.(*aws.CredentialsCache)
	if !ok {
		return
	}
	creds, ok := identity.(*CredentialsAdapter)
	if !ok {
		return
	}
	cache.InvalidateCredentials(creds.Credentials)
}
