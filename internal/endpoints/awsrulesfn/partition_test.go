package awsrulesfn

import (
	"testing"
)

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
			p := GetPartition(partitions, c.Region)
			expected := c.PartitionName
			actual := p.Name
			if expected != *actual {
				t.Errorf("expected %v, got %v", expected, actual)
			}
		})
	}
}
