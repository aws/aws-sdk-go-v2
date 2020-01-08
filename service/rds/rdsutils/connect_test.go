package rdsutils_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
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

	provider := aws.NewStaticCredentialsProvider("AKID", "SECRET", "SESSION")
	var i interface{} = v4.NewSigner(provider)
	for _, c := range cases {
		if signer, ok := i.(rdsutils.HTTPV4Signer); ok {
			url, err := rdsutils.BuildAuthToken(context.Background(), c.endpoint, c.region, c.user, signer)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}
			if re, a := regexp.MustCompile(c.expectedRegex), url; !re.MatchString(a) {
				t.Errorf("expect %s to match %s", re, a)
			}
		} else {
			t.Errorf("signer does not satisfy HTTPV4Signer interface")
		}
	}
}
