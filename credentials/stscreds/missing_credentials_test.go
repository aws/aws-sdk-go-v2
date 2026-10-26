package stscreds_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type mockAssumeRoleFunc func(context.Context, *sts.AssumeRoleInput, ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)

func (m mockAssumeRoleFunc) AssumeRole(ctx context.Context, in *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	return m(ctx, in, optFns...)
}

type stubTokenRetriever struct{}

func (stubTokenRetriever) GetIdentityToken() ([]byte, error) { return []byte("token"), nil }

func TestProvidersRejectNilCredentials(t *testing.T) {
	t.Run("AssumeRole", func(t *testing.T) {
		client := mockAssumeRoleFunc(func(context.Context, *sts.AssumeRoleInput, ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
			return &sts.AssumeRoleOutput{}, nil
		})
		_, err := stscreds.NewAssumeRoleProvider(client, roleARN).Retrieve(context.Background())
		if err == nil || !strings.Contains(err.Error(), "did not contain credentials") {
			t.Fatalf("expect missing credentials error, got %v", err)
		}
	})

	t.Run("AssumeRoleWithWebIdentity", func(t *testing.T) {
		client := mockAssumeRoleWithWebIdentity(func(context.Context, *sts.AssumeRoleWithWebIdentityInput, ...func(*sts.Options)) (*sts.AssumeRoleWithWebIdentityOutput, error) {
			return &sts.AssumeRoleWithWebIdentityOutput{}, nil
		})
		_, err := stscreds.NewWebIdentityRoleProvider(client, roleARN, stubTokenRetriever{}).Retrieve(context.Background())
		if err == nil || !strings.Contains(err.Error(), "did not contain credentials") {
			t.Fatalf("expect missing credentials error, got %v", err)
		}
	})
}

// An empty 2xx response from STS deserializes as success with nil members.
func TestProvidersEmptySTSResponseBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := sts.New(sts.Options{
		Region:       "us-east-1",
		BaseEndpoint: aws.String(srv.URL),
		Credentials:  aws.AnonymousCredentials{},
		Retryer:      aws.NopRetryer{},
	})

	t.Run("AssumeRole", func(t *testing.T) {
		cache := aws.NewCredentialsCache(stscreds.NewAssumeRoleProvider(client, roleARN))
		if _, err := cache.Retrieve(context.Background()); err == nil {
			t.Fatal("expect error, got nil")
		}
	})

	t.Run("AssumeRoleWithWebIdentity", func(t *testing.T) {
		cache := aws.NewCredentialsCache(stscreds.NewWebIdentityRoleProvider(client, roleARN, stubTokenRetriever{}))
		if _, err := cache.Retrieve(context.Background()); err == nil {
			t.Fatal("expect error, got nil")
		}
	})
}
