package v4a

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"testing"
)

type rotatingCredsProvider struct {
	count int
}

func (r *rotatingCredsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	credentials := aws.Credentials{
		AccessKeyID:     fmt.Sprintf("ACCESS_KEY_ID_%d", r.count),
		SecretAccessKey: fmt.Sprintf("SECRET_ACCESS_KEY_%d", r.count),
		SessionToken:    fmt.Sprintf("SESSION_TOKEN_%d", r.count),
	}
	return credentials, nil
}

func TestSymmetricCredentialAdaptor(t *testing.T) {
	provider := &rotatingCredsProvider{}

	adaptor := &SymmetricCredentialAdaptor{SymmetricProvider: provider}

	if symCreds, err := adaptor.Retrieve(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	} else if !symCreds.HasKeys() {
		t.Fatalf("expect symmetric credentials to have keys")
	}

	if load := adaptor.asymmetric.Load(); load != nil {
		t.Errorf("expect asymmetric credentials to be nil")
	}

	if asymCreds, err := adaptor.RetrievePrivateKey(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	} else if !asymCreds.HasKeys() {
		t.Fatalf("expect asymmetric credentials to have keys")
	}

	if _, err := adaptor.Retrieve(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if load := adaptor.asymmetric.Load(); load.(*Credentials) == nil {
		t.Errorf("expect asymmetric credentials to be not nil")
	}

	provider.count++

	if _, err := adaptor.Retrieve(context.Background()); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if load := adaptor.asymmetric.Load(); load.(*Credentials) != nil {
		t.Errorf("expect asymmetric credentials to be nil")
	}
}
