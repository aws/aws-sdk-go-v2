package rdsutils_test

import (
	"regexp"
	"testing"

	credentials "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/rdsutils"
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
		provider := credentials.NewStaticCredentialsProvider("AKID", "SECRET", "SESSION")
		url, err := rdsutils.BuildAuthToken(c.endpoint, c.region, c.user, provider)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
		if re, a := regexp.MustCompile(c.expectedRegex), url; !re.MatchString(a) {
			t.Errorf("expect %s to match %s", re, a)
		}
	}
}
