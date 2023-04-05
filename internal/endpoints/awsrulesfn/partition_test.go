package awsrulesfn

import (
	"testing"
	"github.com/aws/smithy-go/ptr"
)

var mockPartitions = []Partition{
	{
		ID:          "aws",
		RegionRegex: "^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$",
		DefaultConfig: PartitionConfig{
			Name:               ptr.String("aws"),
			DNSSuffix:          ptr.String("amazonaws.com"),
			DualStackDNSSuffix: ptr.String("api.aws"),
			SupportsFIPS:       ptr.Bool(true),
			SupportsDualStack:  ptr.Bool(true),
		},
		Regions: map[string]PartitionConfig{
			"af-south-1": {
				Name:               nil,
				DNSSuffix:          nil,
				DualStackDNSSuffix: nil,
				SupportsFIPS:       nil,
				SupportsDualStack:  nil,
			},
			"us-west-2": {
				Name:               nil,
				DNSSuffix:          nil,
				DualStackDNSSuffix: nil,
				SupportsFIPS:       nil,
				SupportsDualStack:  nil,
			},
		},
	},
	{
		ID:          "aws-cn",
		RegionRegex: "^cn\\-\\w+\\-\\d+$",
		DefaultConfig: PartitionConfig{
			Name:               ptr.String("aws-cn"),
			DNSSuffix:          ptr.String("amazonaws.com.cn"),
			DualStackDNSSuffix: ptr.String("api.amazonwebservices.com.cn"),
			SupportsFIPS:       ptr.Bool(true),
			SupportsDualStack:  ptr.Bool(true),
		},
		Regions: map[string]PartitionConfig{
			"aws-cn-global": {
				Name:               nil,
				DNSSuffix:          nil,
				DualStackDNSSuffix: nil,
				SupportsFIPS:       nil,
				SupportsDualStack:  nil,
			},
			"cn-north-1": {
				Name:               nil,
				DNSSuffix:          nil,
				DualStackDNSSuffix: nil,
				SupportsFIPS:       nil,
				SupportsDualStack:  nil,
			},
			"cn-northwest-1": {
				Name:               nil,
				DNSSuffix:          nil,
				DualStackDNSSuffix: nil,
				SupportsFIPS:       nil,
				SupportsDualStack:  nil,
			},
		},
	},
}

func TestGetPartition(t *testing.T) {
	cases := map[string]struct {
		Region        string
		PartitionName string
	}{
		"test region match aws": {
			Region: "us-west-2", PartitionName: "aws",
		},
		"test region match aws-cn": {
			Region: "aws-cn-global", PartitionName: "aws-cn",
		},
		"test invalid region; default partition": {
			Region: "foo", PartitionName: "aws",
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {

			// monkey patch the partitions data structure
			// thats used by the GetPartition func
			partitions = mockPartitions

			p := GetPartition(c.Region)
			expected := c.PartitionName
			actual := p.Name
			if expected != *actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}
