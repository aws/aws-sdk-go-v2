package awsrulesfn

import (
	"strings"
	"testing"

	smithyrulesfn "github.com/aws/smithy-go/private/endpoints/rulesfn"
)

func TestIsVirtualHostableS3Bucket(t *testing.T) {
	cases := map[string]struct {
		input           string
		allowSubDomains bool
		expect          bool
		expectErr       string
	}{
		"single label no split": {
			input:  "abc123-",
			expect: true,
		},
		"single label no split too short": {
			input:     "a",
			expectErr: `host label 0 has invalid length, "a", 1`,
		},
		"single label with split": {
			input:           "abc123-",
			allowSubDomains: true,
			expect:          true,
		},
		"multiple labels no split": {
			input:     "abc.123-",
			expectErr: `host label 0 is invalid, "abc.123-"`,
		},
		"multiple labels with split": {
			input:           "abc.123-",
			allowSubDomains: true,
			expect:          true,
		},
		"multiple labels with split invalid label": {
			input:           "abc.123-...",
			allowSubDomains: true,
			expectErr:       `host label 2 has invalid length, "", 0`,
		},
		"max length host label": {
			input:  "012345678901234567890123456789012345678901234567890123456789123",
			expect: true,
		},
		"too large host label": {
			input:     "0123456789012345678901234567890123456789012345678901234567891234",
			expectErr: `host label 0 has invalid length, "0123456789012345678901234567890123456789012345678901234567891234", 64`,
		},
		"too small host label": {
			input:     "",
			expectErr: `host label 0 has invalid length, "", 0`,
		},
		"lower case only": {
			input:     "AbC",
			expectErr: `host label 0 cannot have capital letters, "AbC"`,
		},
		"like IP address": {
			input:     "127.111.222.123",
			expectErr: `host label is formatted like IP address, "127.111.222.123"`,
		},
		"multiple labels like IP address": {
			input:           "127.111.222.123",
			allowSubDomains: true,
			expectErr:       `host label is formatted like IP address, "127.111.222.123"`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ec := smithyrulesfn.NewErrorCollector()
			actual := IsVirtualHostableS3Bucket(c.input, c.allowSubDomains, ec)
			if !c.expect {
				if e, a := c.expectErr, ec.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q error in %q", e, a)
				}
			}

			if e, a := c.expect, actual; e != a {
				t.Fatalf("expect %v hostable bucket, got %v", e, a)
			}
		})
	}
}
