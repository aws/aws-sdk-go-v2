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
func GetPartition(region string) *PartitionConfig {
	return getPartition(partitions, region)
}

func getPartition(partitions []Partition, region string) *PartitionConfig {
	for _, partition := range partitions {
		if v, ok := partition.Regions[region]; ok {
			v = mergePartition(v, partition.DefaultConfig)
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

func mergePartition(into PartitionConfig, from PartitionConfig) PartitionConfig {
	if into.Name == nil {
		into.Name = from.Name
	}
	if into.DNSSuffix == nil {
		into.DNSSuffix = from.DNSSuffix
	}
	if into.DualStackDNSSuffix == nil {
		into.DualStackDNSSuffix = from.DualStackDNSSuffix
	}
	if into.SupportsFIPS == nil {
		into.SupportsFIPS = from.SupportsFIPS
	}
	if into.SupportsDualStack == nil {
		into.SupportsDualStack = from.SupportsDualStack
	}
	return into
}
