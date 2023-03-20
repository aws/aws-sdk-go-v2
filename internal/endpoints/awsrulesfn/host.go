package awsrulesfn

import (
	"fmt"
	"net/netip"
	"strings"
	smithyrulesfn "github.com/aws/smithy-go/private/endpoints/rulesfn"
	smithyhttp "github.com/aws/smithy-go/transport/http"

)

// IsVirtualHostableS3Bucket returns if the input is a DNS compatible bucket
// name and can be used with Amazon S3 virtual hosted style addressing. Similar
// to [rulesfn.IsValidHostLabel] with the added restriction that the length of label
// must be [3:63] characters long, all lowercase, and not formatted as an IP
// address.
func IsVirtualHostableS3Bucket(input string, allowSubDomains bool, ec *smithyrulesfn.ErrorCollector) bool {
	// input should not be formatted as an IP address
	if _, err := netip.ParseAddr(input); err == nil {
		ec.AddError(smithyrulesfn.FnError{
			Name: "IsVirtualHostableS3Bucket",
			Err:  fmt.Errorf("host label is formatted like IP address, %q", input),
		})
		return false
	}

	var labels []string
	if allowSubDomains {
		labels = strings.Split(input, ".")
	} else {
		labels = []string{input}
	}

	for i, label := range labels {
		// validate special length constraints
		if l := len(label); l < 3 || l > 63 {
			ec.AddError(smithyrulesfn.FnError{
				Name: "IsVirtualHostableS3Bucket",
				Err:  fmt.Errorf("host label %d has invalid length, %q, %d", i, label, l),
			})
			return false
		}

		// Validate no capital letters
		for _, r := range label {
			if r >= 'A' && r <= 'Z' {
				ec.AddError(smithyrulesfn.FnError{
					Name: "IsVirtualHostableS3Bucket",
					Err:  fmt.Errorf("host label %d cannot have capital letters, %q", i, label),
				})
				return false
			}
		}

		// Validate valid host label
		if !smithyhttp.ValidHostLabel(label) {
			ec.AddError(smithyrulesfn.FnError{
				Name: "IsVirtualHostableS3Bucket",
				Err:  fmt.Errorf("host label %d is invalid, %q", i, label),
			})
			return false
		}
	}

	return true
}
