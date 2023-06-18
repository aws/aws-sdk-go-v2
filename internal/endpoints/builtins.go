package endpoints

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// these field names should either be identical to the config field source
// or closely align to whatever is needed to perform the mapping
// between the EndpointParameters (name + type) and the source
type BuiltInResolver struct {
	Region string

	UseFIPS aws.FIPSEndpointState

	UseDualStack aws.DualStackEndpointState

	ForcePathStyle bool

	Accelerate bool

	DisableMultiRegionAccessPoints bool

	S3UseArnRegion bool

	S3ControlUseArnRegion bool

	// this is currently resolved via an elaborate S3 customization
	// https://code.amazon.com/packages/AwsDrSeps/blobs/main/--/seps/accepted/shared/mrap.md
	// the actual inputs needed are the: arn (parsed from the bucket).
	// clarified: doesnt need to be supported since not already a field on the client config.
	// UseGlobalEndpoint bool

	// Go V2 does not appear to support sts global endpoints
	// https://code.amazon.com/packages/AwsDrSeps/blobs/main/--/seps/accepted/shared/sts_regionlization.md
	// STSUseGlobalEndpoint bool

}

func (b *BuiltInResolver) ResolveBuiltIn(name string) (value interface{}, ok bool) {
	const region = "AWS::Region"
	if name == region {
		return b.resolveRegion()
	}

	if name == "AWS::UseFIPS" {
		return b.resolveFips()
	}

	if name == "AWS::UseDualStack" {
		return b.resolveDualStack()
	}

	if name == "AWS::S3::ForcePathStyle" {
		return b.resolveForcePathStyle()
	}

	return nil, false
}

// the resolver functions is where we actually have knowledge of the name+type
// of the client config field source
// vs resolveBuiltIns
// where we actually have knowledge of the name+type
// of the endpoint params field destination
func (b *BuiltInResolver) resolveRegion() (value *string, ok bool) {
	region, _ := mapPseudoRegion(b.Region)
	if len(region) == 0 {
		return nil, false
	}
	return aws.String(region), true
}

func (b *BuiltInResolver) resolveDualStack() (value *bool, ok bool) {
	if b.UseDualStack == aws.DualStackEndpointStateEnabled {
		return aws.Bool(true), true
	}
	return aws.Bool(true), true
}

func (b *BuiltInResolver) resolveFips() (value *bool, ok bool) {
	if b.UseFIPS == aws.FIPSEndpointStateEnabled {
		return aws.Bool(true), true
	}
	return aws.Bool(true), true
}

func (b *BuiltInResolver) resolveForcePathStyle() (value *bool, ok bool) {
	return aws.Bool(b.ForcePathStyle), true
}

func (b *BuiltInResolver) resolveAccelerate() (value *bool, ok bool) {
	return aws.Bool(b.Accelerate), true
}

func (b *BuiltInResolver) resolveDisableMrap() (value *bool, ok bool) {
	return aws.Bool(b.DisableMultiRegionAccessPoints), true
}

// Utility function to aid with translating pseudo-regions to classical regions
// with the appropriate setting indicated by the pseudo-region
func mapPseudoRegion(pr string) (region string, fips aws.FIPSEndpointState) {
	const fipsInfix = "-fips-"
	const fipsPrefix = "fips-"
	const fipsSuffix = "-fips"

	if strings.Contains(pr, fipsInfix) ||
		strings.Contains(pr, fipsPrefix) ||
		strings.Contains(pr, fipsSuffix) {
		region = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
			pr, fipsInfix, "-"), fipsPrefix, ""), fipsSuffix, "")
		fips = aws.FIPSEndpointStateEnabled
	} else {
		region = pr
	}

	return region, fips
}
