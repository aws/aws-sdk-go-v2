package awsrulesfn

import "regexp"

// Partition provides the metadata describing an AWS partition.
type Partition struct {
	ID            string                     `json:"id"`
	Regions       map[string]PartitionConfig `json:"regions"`
	RegionRegex   string                     `json:"regionRegex"`
	DefaultConfig PartitionConfig            `json:"outputs"`
}

// PartitionConfig provides the endpoint metadata for an AWS region or partition.
type PartitionConfig struct {
	Name               *string `json:"name"`
	DNSSuffix          *string `json:"dnsSuffix"`
	DualStackDNSSuffix *string `json:"dualStackDnsSuffix"`
	SupportsFIPS       *bool   `json:"supportsFIPS"`
	SupportsDualStack  *bool   `json:"supportsDualStack"`
}

const defaultPartition = "aws"

// GetPartition returns an AWS [Partition] for the region provided. If the
// partition cannot be determined nil will be returned.
func GetPartition(partitions []Partition, region string) *PartitionConfig {
	for _, partition := range partitions {
		if v, ok := partition.Regions[region]; ok {
			v.mergeWith(partition.DefaultConfig)
			return &v
		}
	}

	for _, partition := range partitions {
		regionRegex := regexp.MustCompile(partition.RegionRegex)
		if regionRegex.MatchString(region) {
			v := partition.DefaultConfig
			return &v
		}
	}

	for _, partition := range partitions {
		if partition.ID == defaultPartition {
			v := partition.DefaultConfig
			return &v
		}
	}

	return nil
}

func (p *PartitionConfig) mergeWith(other PartitionConfig) {
	if p.Name == nil {
		p.Name = other.Name
	}
	if p.DNSSuffix == nil {
		p.DNSSuffix = other.DNSSuffix
	}
	if p.DualStackDNSSuffix == nil {
		p.DualStackDNSSuffix = other.DualStackDNSSuffix
	}
	if p.SupportsFIPS == nil {
		p.SupportsFIPS = other.SupportsFIPS
	}
	if p.SupportsDualStack == nil {
		p.SupportsDualStack = other.SupportsDualStack
	}
}
