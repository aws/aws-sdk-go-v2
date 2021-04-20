package auth_test

import (
	"context"
	"regexp"
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
	}{
		{
			"https://prod-instance.us-east-1.rds.amazonaws.com:3306",
			"us-west-2",
			"mysqlUser",
			`^prod-instance\.us-east-1\.rds\.amazonaws\.com:3306\?Action=connect.*?DBUser=mysqlUser.*`,
		},
		{
			"prod-instance.us-east-1.rds.amazonaws.com:3306",
			"us-west-2",
			"mysqlUser",
			`^prod-instance\.us-east-1\.rds\.amazonaws\.com:3306\?Action=connect.*?DBUser=mysqlUser.*`,
		},
	}

	for _, c := range cases {
		creds := &staticCredentials{AccessKey: "AKID", SecretKey: "SECRET", Session: "SESSION"}
		url, err := auth.BuildAuthToken(context.Background(), c.endpoint, c.region, c.user, creds)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
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
