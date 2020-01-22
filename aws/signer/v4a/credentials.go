package v4a

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
)

// Credentials is Context, ECDSA, and Optional Session Token that can be used
// to sign requests using SigV4a
type Credentials struct {
	Context      string
	PrivateKey   *ecdsa.PrivateKey
	SessionToken string

	// Time the credentials will expire.
	CanExpire bool
	Expires   time.Time
}

// Expired returns if the credentials have expired.
func (v Credentials) Expired() bool {
	if v.CanExpire {
		return !v.Expires.After(sdk.NowTime())
	}

	return false
}

// HasKeys returns if the credentials keys are set.
func (v Credentials) HasKeys() bool {
	return len(v.Context) > 0 && v.PrivateKey != nil
}

// SymmetricCredentialAdaptor wraps a SigV4 AccessKey/SecretKey provider and adapts the credentials
// to a ECDSA PrivateKey for signing with SiV4a
type SymmetricCredentialAdaptor struct {
	SymmetricProvider aws.CredentialsProvider

	creds atomic.Value
	m     sync.Mutex
}

// RetrievePrivateKey returns credentials suitable for SigV4a signing
func (s *SymmetricCredentialAdaptor) RetrievePrivateKey(ctx context.Context) (Credentials, error) {
	if creds := s.getCreds(); creds != nil {
		return *creds, nil
	}

	s.m.Lock()
	defer s.m.Unlock()

	if creds := s.getCreds(); creds != nil {
		return *creds, nil
	}

	symmetricCreds, err := s.retrieveFromSymmetricProvider(ctx)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to retrieve symmetric credentials: %v", err)
	}

	privateKey, err := deriveKeyFromAccessKeyPair(symmetricCreds.AccessKeyID, symmetricCreds.SecretAccessKey)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to derive assymetric key from credentials")
	}

	creds := Credentials{
		Context:      symmetricCreds.AccessKeyID,
		PrivateKey:   privateKey,
		SessionToken: symmetricCreds.SessionToken,
		CanExpire:    symmetricCreds.CanExpire,
		Expires:      symmetricCreds.Expires,
	}

	s.creds.Store(&creds)

	return creds, nil
}

func (s *SymmetricCredentialAdaptor) getCreds() *Credentials {
	v := s.creds.Load()

	if v == nil {
		return nil
	}

	c := v.(*Credentials)
	if c != nil && c.HasKeys() && !c.Expired() {
		return c
	}

	return nil
}

func (s *SymmetricCredentialAdaptor) retrieveFromSymmetricProvider(ctx context.Context) (aws.Credentials, error) {
	credentials, err := s.SymmetricProvider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	return credentials, nil
}

type CredentialsProvider interface {
	RetrievePrivateKey(context.Context) (Credentials, error)
}
