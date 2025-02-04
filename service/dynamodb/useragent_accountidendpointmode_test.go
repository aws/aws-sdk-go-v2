package dynamodb

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
)

type captureUserAgentFeatures struct {
	ua       string
	features []string
}

func (u *captureUserAgentFeatures) Do(r *http.Request) (*http.Response, error) {
	u.ua = r.Header.Get("User-Agent")
	parts := strings.Split(u.ua, " ")

	var features string
	for _, part := range parts {
		if strings.HasPrefix(part, "m/") {
			features = part
			break
		}
	}

	if len(features) > 2 {
		u.features = strings.Split(features[2:], ",")
	}
	return &http.Response{StatusCode: 403, Body: http.NoBody}, nil
}

type mockCredentials struct{}

// bypass the "do you have an account ID if necessary" checker
func (mockCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccountID: "1234567890",
	}, nil
}

func TestUserAgent_AccountIDEndpointMode(t *testing.T) {
	for name, tt := range map[string]struct {
		Config aws.AccountIDEndpointMode
		Expect awsmiddleware.UserAgentFeature
	}{
		"preferred": {aws.AccountIDEndpointModePreferred, awsmiddleware.UserAgentFeatureAccountIDModePreferred},
		"required":  {aws.AccountIDEndpointModeRequired, awsmiddleware.UserAgentFeatureAccountIDModeRequired},
		"disabled":  {aws.AccountIDEndpointModeDisabled, awsmiddleware.UserAgentFeatureAccountIDModeDisabled},
	} {
		t.Run(name, func(t *testing.T) {
			var ua captureUserAgentFeatures
			client := New(Options{
				Region:                "us-east-1",
				AccountIDEndpointMode: tt.Config,
				HTTPClient:            &ua,
				Credentials:           mockCredentials{},
			})

			client.Scan(context.Background(), &ScanInput{
				TableName: aws.String("foo"),
			})
			expectContains(t, ua.features, string(tt.Expect))
		})
	}
}

func expectContains(t *testing.T, have []string, want string) {
	t.Helper()
	for _, s := range have {
		if s == want {
			return
		}
	}
	t.Errorf("[]string %v did not contain %s", have, want)
}
