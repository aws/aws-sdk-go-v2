package rdsutils_test

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/rdsutils"
)

func TestConnectionStringBuilder(t *testing.T) {
	cases := []struct {
		user         string
		endpoint     string
		region       string
		dbName       string
		values       url.Values
		format       rdsutils.ConnectionFormat
		credProvider aws.CredentialsProvider

		expectedErr          error
		expectedConnectRegex string
	}{
		{
			user:                 "foo",
			endpoint:             "foo.bar",
			region:               "region",
			dbName:               "name",
			format:               rdsutils.NoConnectionFormat,
			credProvider:         aws.NewStaticCredentialsProvider("AKID", "SECRET", "SESSION"),
			expectedErr:          rdsutils.ErrNoConnectionFormat,
			expectedConnectRegex: "",
		},
		{
			user:                 "foo",
			endpoint:             "foo.bar",
			region:               "region",
			dbName:               "name",
			format:               rdsutils.TCPFormat,
			credProvider:         aws.NewStaticCredentialsProvider("AKID", "SECRET", "SESSION"),
			expectedConnectRegex: `^foo:foo.bar\?Action=connect\&DBUser=foo.*\@tcp\(foo.bar\)/name`,
		},
	}

	for _, c := range cases {
		b := rdsutils.NewConnectionStringBuilder(c.endpoint, c.region, c.user, c.dbName, c.credProvider)
		connectStr, err := b.WithFormat(c.format).Build()

		if e, a := c.expectedErr, err; e != a {
			t.Errorf("expected %v error, but received %v", e, a)
		}

		if re, a := regexp.MustCompile(c.expectedConnectRegex), connectStr; !re.MatchString(a) {
			t.Errorf("expect %s to match %s", re, a)
		}
	}
}
