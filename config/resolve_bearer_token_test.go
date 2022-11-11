package config

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ssocreds"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	smithybearer "github.com/aws/smithy-go/auth/bearer"
	"github.com/google/go-cmp/cmp"
)

func TestResolveBearerAuthToken(t *testing.T) {
	cases := map[string]struct {
		setNowTime     func() func()
		configs        configs
		expectProvider bool
		expectToken    smithybearer.Token
	}{
		"no provider": {},
		"config source provider": {
			configs: configs{
				LoadOptions{
					BearerAuthTokenProvider: smithybearer.StaticTokenProvider{
						Token: smithybearer.Token{Value: "abc123"},
					},
				},
			},
			expectProvider: true,
			expectToken:    smithybearer.Token{Value: "abc123"},
		},
		"sso session legacy": {
			setNowTime: func() func() {
				return sdk.TestingUseReferenceTime(time.Date(2044, 4, 4, 0, 0, 0, 0, time.UTC))
			},
			configs: configs{
				LoadOptions{
					SSOTokenProviderOptions: func(o *ssocreds.SSOTokenProviderOptions) {
						expectPath, _ := ssocreds.StandardCachedTokenFilepath("https://example.aws/start")
						if e, a := expectPath, o.CachedTokenFilepath; e != a {
							t.Errorf("expect %v cache file path, got %v", e, a)
						}

						o.CachedTokenFilepath = filepath.Join("testdata", "cached_sso_token.json")
					},
				},
				SharedConfig{
					SSORegion:   "us-west-2",
					SSOStartURL: "https://example.aws/start",
				},
			},
			expectProvider: false,
			expectToken: smithybearer.Token{
				Value:     "access token",
				CanExpire: true,
				Expires:   time.Date(2044, 4, 4, 7, 0, 1, 0, time.UTC),
			},
		},
		"sso session named": {
			setNowTime: func() func() {
				return sdk.TestingUseReferenceTime(time.Date(2044, 4, 4, 0, 0, 0, 0, time.UTC))
			},
			configs: configs{
				LoadOptions{
					SSOTokenProviderOptions: func(o *ssocreds.SSOTokenProviderOptions) {
						expectPath, _ := ssocreds.StandardCachedTokenFilepath("test-session")
						if e, a := expectPath, o.CachedTokenFilepath; e != a {
							t.Errorf("expect %v cache file path, got %v", e, a)
						}

						o.CachedTokenFilepath = filepath.Join("testdata", "cached_sso_token.json")
					},
				},
				SharedConfig{
					SSOSessionName: "test-session",
					SSOSession: &SSOSession{
						Name:        "test-session",
						SSORegion:   "us-west-2",
						SSOStartURL: "https://example.aws/start",
					},
				},
			},
			expectProvider: true,
			expectToken: smithybearer.Token{
				Value:     "access token",
				CanExpire: true,
				Expires:   time.Date(2044, 4, 4, 7, 0, 1, 0, time.UTC),
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if c.setNowTime != nil {
				restoreTime := c.setNowTime()
				defer restoreTime()
			}

			var cfg aws.Config
			err := resolveBearerAuthToken(context.Background(), &cfg, c.configs)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if !c.expectProvider {
				if v := cfg.BearerAuthTokenProvider; v != nil {
					t.Errorf("expect no provider, got %T, %v", v, v)
				}
				return
			}

			token, err := cfg.BearerAuthTokenProvider.RetrieveBearerToken(context.Background())
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(c.expectToken, token); diff != "" {
				t.Errorf("expect token match\n%s", diff)
			}
		})
	}
}

func TestWrapWithBearerAuthTokenProvider(t *testing.T) {
	cases := map[string]struct {
		configs         configs
		provider        smithybearer.TokenProvider
		optFns          []func(*smithybearer.TokenCacheOptions)
		compareProvider bool
		expectToken     smithybearer.Token
	}{
		"already wrapped": {
			provider: smithybearer.NewTokenCache(smithybearer.StaticTokenProvider{
				Token: smithybearer.Token{Value: "abc123"},
			}),
			compareProvider: true,
			expectToken:     smithybearer.Token{Value: "abc123"},
		},
		"to be wrapped": {
			provider: smithybearer.StaticTokenProvider{
				Token: smithybearer.Token{Value: "abc123"},
			},
			expectToken: smithybearer.Token{Value: "abc123"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			provider, err := wrapWithBearerAuthTokenCache(context.Background(),
				c.configs, c.provider, c.optFns...)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if p, ok := provider.(*smithybearer.TokenCache); !ok {
				t.Fatalf("expect provider wrapped in %T, got %T", p, provider)
			}

			if c.compareProvider && provider != c.provider {
				t.Errorf("expect same provider, was not")
			}

			token, err := provider.RetrieveBearerToken(context.Background())
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmp.Diff(c.expectToken, token); diff != "" {
				t.Errorf("expect token match\n%s", diff)
			}
		})
	}
}
