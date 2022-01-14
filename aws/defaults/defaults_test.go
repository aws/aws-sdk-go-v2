package defaults

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestConfigV1(t *testing.T) {
	cases := []struct {
		Mode     aws.DefaultsMode
		Expected Configuration
	}{
		{
			Mode: aws.DefaultsModeStandard,
			Expected: Configuration{
				ConnectTimeout:        aws.Duration(2000 * time.Millisecond),
				TLSNegotiationTimeout: aws.Duration(2000 * time.Millisecond),
			},
		},
		{
			Mode: aws.DefaultsModeInRegion,
			Expected: Configuration{
				ConnectTimeout:        aws.Duration(1000 * time.Millisecond),
				TLSNegotiationTimeout: aws.Duration(1000 * time.Millisecond),
			},
		},
		{
			Mode: aws.DefaultsModeCrossRegion,
			Expected: Configuration{
				ConnectTimeout:        aws.Duration(2800 * time.Millisecond),
				TLSNegotiationTimeout: aws.Duration(2800 * time.Millisecond),
			},
		},
		{
			Mode: aws.DefaultsModeMobile,
			Expected: Configuration{
				ConnectTimeout:        aws.Duration(10000 * time.Millisecond),
				TLSNegotiationTimeout: aws.Duration(11000 * time.Millisecond),
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := v1TestResolver(tt.Mode)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if diff := cmp.Diff(tt.Expected, got); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}
