package awsrulesfn

import (
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"net/netip"
	"strings"
)

// IsVirtualHostableS3Bucket returns if the input is a DNS compatible bucket
// name and can be used with Amazon S3 virtual hosted style addressing. Similar
// to [rulesfn.IsValidHostLabel] with the added restriction that the length of label
// must be [3:63] characters long, all lowercase, and not formatted as an IP
// address.
func IsVirtualHostableS3Bucket(input string, allowSubDomains bool) bool {
	// input should not be formatted as an IP address
	if _, err := netip.ParseAddr(input); err == nil {
		return false
	}

	var labels []string
	if allowSubDomains {
		labels = strings.Split(input, ".")
	} else {
		labels = []string{input}
	}

	for _, label := range labels {
		// validate special length constraints
		if l := len(label); l < 3 || l > 63 {
			return false
		}

		// Validate no capital letters
		for _, r := range label {
			if r >= 'A' && r <= 'Z' {
				return false
			}
		}

		// Validate valid host label
		if !smithyhttp.ValidHostLabel(label) {
			return false
		}
	}

	return true
}
