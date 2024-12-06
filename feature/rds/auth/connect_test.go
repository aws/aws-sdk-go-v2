package auth_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

func TestBuildAuthToken(t *testing.T) {
	cases := []struct {
		endpoint      string
		region        string
		user          string
		expectedRegex string
		expectedError string
	}{
		{
			endpoint:      "https://prod-instance.us-east-1.rds.amazonaws.com:3306",
			region:        "us-west-2",
			user:          "mysqlUser",
			expectedRegex: `^prod-instance\.us-east-1\.rds\.amazonaws\.com:3306\?Action=connect.*?DBUser=mysqlUser.*`,
		},
		{
			endpoint:      "prod-instance.us-east-1.rds.amazonaws.com:3306",
			region:        "us-west-2",
			user:          "mysqlUser",
			expectedRegex: `^prod-instance\.us-east-1\.rds\.amazonaws\.com:3306\?Action=connect.*?DBUser=mysqlUser.*`,
		},
		{
			endpoint:      "prod-instance.us-east-1.rds.amazonaws.com",
			region:        "us-west-2",
			user:          "mysqlUser",
			expectedError: "port",
		},
		{
			endpoint:      "prod-instance.us-east-1.rds.amazonaws.com:kakasdkasd",
			region:        "us-west-2",
			user:          "mysqlUser",
			expectedError: "port",
		},
	}

	for _, c := range cases {
		creds := &staticCredentials{AccessKey: "AKID", SecretKey: "SECRET", Session: "SESSION"}
		url, err := auth.BuildAuthToken(context.Background(), c.endpoint, c.region, c.user, creds)
		if len(c.expectedError) > 0 {
			if err != nil {
				if !strings.Contains(err.Error(), c.expectedError) {
					t.Fatalf("expect err: %v, actual err: %v", c.expectedError, err)
				} else {
					continue
				}
			} else {
				t.Fatalf("expect err: %v, actual err: %v", c.expectedError, err)
			}
		} else if err != nil {
			t.Fatalf("expect no err, got: %v", err)
		}

		if re, a := regexp.MustCompile(c.expectedRegex), url; !re.MatchString(a) {
			t.Errorf("expect %s to match %s", re, a)
		}
	}
}

type staticCredentials struct {
	AccessKey, SecretKey, Session string
}

func (s *staticCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     s.AccessKey,
		SecretAccessKey: s.SecretKey,
		SessionToken:    s.Session,
	}, nil
}
