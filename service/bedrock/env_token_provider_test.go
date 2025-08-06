package bedrock

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/auth/bearer"
)

// Valid service-specific token configured
func TestEnvTokenProviderSEP0(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "bedrock-token")

	svc := New(Options{})
	token, err := svc.Options().BearerAuthTokenProvider.RetrieveBearerToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token.Value != "bedrock-token" {
		t.Errorf("expect token \"bedrock-token\" != %q", token.Value)
	}

	expectPref := []string{"httpBearerAuth"}
	if actual := svc.Options().AuthSchemePreference; !reflect.DeepEqual(expectPref, actual) {
		t.Errorf("expect auth scheme preference %#v != %#v", expectPref, actual)
	}
}

// Token configured for a different service
func TestEnvTokenProviderSEP1(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_FOO", "foo-token")

	svc := New(Options{})
	if actual := svc.Options().BearerAuthTokenProvider; actual != nil {
		t.Errorf("BearerAuthTokenProvider should be nil but it's %#v", actual)
	}

	expectPref := []string(nil)
	if actual := svc.Options().AuthSchemePreference; !reflect.DeepEqual(expectPref, actual) {
		t.Errorf("expect auth scheme preference %#v != %#v", expectPref, actual)
	}
}

// Token configured with auth scheme preference also set in env. The Bedrock
// token is more specific which is why it takes precedence
func TestEnvTokenProviderSEP2(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "bedrock-token")

	cfg := aws.Config{
		// simulate loaded from env/shared config
		AuthSchemePreference: []string{"sigv4", "httpBearerAuth"},
	}
	svc := NewFromConfig(cfg)

	token, err := svc.Options().BearerAuthTokenProvider.RetrieveBearerToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if token.Value != "bedrock-token" {
		t.Errorf("expect token \"bedrock-token\" != %q", token.Value)
	}

	expectPref := []string{"httpBearerAuth"}
	if actual := svc.Options().AuthSchemePreference; !reflect.DeepEqual(expectPref, actual) {
		t.Errorf("expect auth scheme preference %#v != %#v", expectPref, actual)
	}
}

// Explicit service config takes precedence
func TestEnvTokenProviderSEP3(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "bedrock-token")

	expectToken := "explicit-code-token"
	svc := New(Options{}, func(o *Options) {
		o.BearerAuthTokenProvider = bearer.TokenProviderFunc(func(ctx context.Context) (bearer.Token, error) {
			return bearer.Token{Value: expectToken}, nil
		})
		// also test explicit code preference
		o.AuthSchemePreference = []string{"httpBasicAuth", "httpBearerAuth"}
	})

	token, err := svc.Options().BearerAuthTokenProvider.RetrieveBearerToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if expectToken != token.Value {
		t.Errorf("expect token %q != %q", expectToken, token.Value)
	}

	expectPref := []string{"httpBasicAuth", "httpBearerAuth"}
	if actual := svc.Options().AuthSchemePreference; !reflect.DeepEqual(expectPref, actual) {
		t.Errorf("expect auth scheme preference %#v != %#v", expectPref, actual)
	}
}
