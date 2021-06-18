package endpoints

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
)

func TestEndpointResolve(t *testing.T) {
	defs := Endpoint{
		Hostname:          "service.{region}.amazonaws.com",
		SignatureVersions: []string{"v4"},
	}

	e := Endpoint{
		Protocols:         []string{"http", "https"},
		SignatureVersions: []string{"v4"},
		CredentialScope: CredentialScope{
			Region:  "us-west-2",
			Service: "service",
		},
	}

	resolved, err := e.resolve("aws", "us-west-2", defs, Options{})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if e, a := "https://service.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "aws", resolved.PartitionID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "v4", resolved.SigningMethod; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEndpointMergeIn(t *testing.T) {
	expected := Endpoint{
		Hostname:          "other hostname",
		Protocols:         []string{"http"},
		SignatureVersions: []string{"v4"},
		CredentialScope: CredentialScope{
			Region:  "region",
			Service: "service",
		},
	}

	actual := Endpoint{}
	actual.mergeIn(Endpoint{
		Hostname:          "other hostname",
		Protocols:         []string{"http"},
		SignatureVersions: []string{"v4"},
		CredentialScope: CredentialScope{
			Region:  "region",
			Service: "service",
		},
	})

	if e, a := expected, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

type testCase struct {
	Region   string
	Options  Options
	Expected aws.Endpoint
	WantErr  bool
}

type serviceTest struct {
	Partitions Partitions
	Cases      []testCase
}

var testCases = map[string]serviceTest{
	"s3": {
		Partitions: Partitions{
			{
				ID:          "aws",
				RegionRegex: regexp.MustCompile("^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$"),
				Defaults: map[EndpointVariant]Endpoint{
					0: {
						Hostname:  "s3.{region}.amazonaws.com",
						Protocols: []string{"https"},
						CredentialScope: CredentialScope{
							Service: "s3",
						},
						SignatureVersions: []string{"s3v4"},
					},
					DualStackVariant: {
						Hostname:  "s3.dualstack.{region}.amazonaws.com",
						Protocols: []string{"https"},
						CredentialScope: CredentialScope{
							Service: "s3",
						},
						SignatureVersions: []string{"s3v4"},
					},
					FIPSVariant: {
						Hostname:  "s3-fips.{region}.amazonaws.com",
						Protocols: []string{"https"},
						CredentialScope: CredentialScope{
							Service: "s3",
						},
						SignatureVersions: []string{"s3v4"},
					},
					DualStackVariant | FIPSVariant: {
						Hostname:  "s3-fips.{region}.api.aws",
						Protocols: []string{"https"},
						CredentialScope: CredentialScope{
							Service: "s3",
						},
						SignatureVersions: []string{"s3v4"},
					},
				},
				IsRegionalized: true,
				Endpoints: map[EndpointKey]Endpoint{
					{
						Region: "us-west-2",
					}: {
						Hostname: "s3.api.us-west-2.amazonaws.com",
					},
					{
						Region:  "us-west-2",
						Variant: DualStackVariant,
					}: {
						Hostname: "s3.api.dualstack.us-west-2.amazonaws.com",
					},
					{
						Region:  "us-west-2",
						Variant: FIPSVariant,
					}: {
						Hostname: "s3-fips.api.us-west-2.amazonaws.com",
					},
				},
			},
		},
		Cases: []testCase{
			// General Test Cases
			{
				Region: "us-west-2",
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3.api.us-west-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-west-2",
					SigningMethod: "s3v4",
				},
			},
			{
				Region: "us-east-2",
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3.us-east-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-east-2",
					SigningMethod: "s3v4",
				},
			},
			{
				Region: "us-east-2",
				Options: Options{
					DisableHTTPS: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://s3.us-east-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-east-2",
					SigningMethod: "s3v4",
				},
			},
			// Dual-Stack cases
			{
				Region: "us-west-2",
				Options: Options{
					UseDualStack: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3.api.dualstack.us-west-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-west-2",
					SigningMethod: "s3v4",
				},
			},
			{
				Region: "us-east-2",
				Options: Options{
					UseDualStack: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3.dualstack.us-east-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-east-2",
					SigningMethod: "s3v4",
				},
			},
			// FIPS Test Cases
			{
				Region: "us-west-2",
				Options: Options{
					UseFIPS: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3-fips.api.us-west-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-west-2",
					SigningMethod: "s3v4",
				},
			},
			{
				Region: "us-east-2",
				Options: Options{
					UseFIPS: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://s3-fips.us-east-2.amazonaws.com",
					SigningName:   "s3",
					SigningRegion: "us-east-2",
					SigningMethod: "s3v4",
				},
			},
		},
	},
	"dynamodb": {
		Partitions: Partitions{
			{
				ID:          "aws",
				RegionRegex: regexp.MustCompile("^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$"),
				Defaults: map[EndpointVariant]Endpoint{
					FIPSVariant: {
						Hostname:          "dynamodb-fips.{region}.amazonaws.com",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
					DualStackVariant: {
						Hostname:          "dynamodb.{region}.api.aws",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
					DualStackVariant | FIPSVariant: {
						Hostname:          "dynamodb-fips.{region}.api.aws",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
					0: {
						Hostname:          "dynamodb.{region}.amazonaws.com",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
				},
				IsRegionalized: true,
				Endpoints: map[EndpointKey]Endpoint{
					{
						Region:  "us-west-2",
						Variant: FIPSVariant,
					}: {
						Hostname: "dynamodb-fips.us-west-2.amazonaws.com",
					},
					{
						Region:  "us-west-2",
						Variant: FIPSVariant | DualStackVariant,
					}: {
						Hostname: "fips.dynamodb.us-west-2.api.aws",
					},
					{
						Region:  "us-west-2",
						Variant: DualStackVariant,
					}: {},
				},
			},
		},
		Cases: []testCase{
			{
				Region: "us-west-2",
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb.us-west-2.amazonaws.com",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-east-2",
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb.us-east-2.amazonaws.com",
					SigningRegion: "us-east-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-west-2",
				Options: Options{
					UseDualStack: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb.us-west-2.api.aws",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-east-2",
				Options: Options{
					UseDualStack: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb.us-east-2.api.aws",
					SigningRegion: "us-east-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-west-2",
				Options: Options{
					UseFIPS: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb-fips.us-west-2.amazonaws.com",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-east-2",
				Options: Options{
					UseFIPS: true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://dynamodb-fips.us-east-2.amazonaws.com",
					SigningRegion: "us-east-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-west-2",
				Options: Options{
					UseDualStack: true,
					UseFIPS:      true,
				},
				Expected: aws.Endpoint{
					PartitionID:   "aws",
					URL:           "https://fips.dynamodb.us-west-2.api.aws",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
		},
	},
	"ec2": {
		Partitions: Partitions{
			{
				ID:          "aws",
				RegionRegex: regexp.MustCompile("^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$"),
				Defaults: map[EndpointVariant]Endpoint{
					0: {
						Hostname:          "api.ec2.{region}.amazonaws.com",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
					DualStackVariant: {
						Hostname:          "api.ec2.{region}.api.aws",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
					FIPSVariant: {
						Hostname:          "api.ec2-fips.{region}.amazonaws.com",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
				},
				IsRegionalized: true,
				Endpoints: map[EndpointKey]Endpoint{
					{
						Region: "us-west-2",
					}: {
						CredentialScope: CredentialScope{
							Region: "us-west-2",
						},
						Hostname: "ec2.us-west-2.amazonaws.com",
					},
					{
						Region:  "us-west-2",
						Variant: DualStackVariant,
					}: {
						CredentialScope: CredentialScope{
							Region: "us-west-2",
						},
						Hostname: "ec2.us-west-2.api.aws",
					},
					{
						Region:  "us-west-2",
						Variant: FIPSVariant,
					}: {},
					{
						Region: "fips-us-west-2",
					}: {
						Hostname: "ec2-fips.us-west-2.amazonaws.com",
						CredentialScope: CredentialScope{
							Region: "us-west-2",
						},
						Deprecated: true,
					},
				},
			},
			{
				ID:          "aws-iso",
				RegionRegex: regexp.MustCompile("^us\\-iso\\-\\w+\\-\\d+$"),
				Defaults: map[EndpointVariant]Endpoint{
					0: {
						Hostname:          "ec2.{region}.c2s.ic.gov",
						Protocols:         []string{"http", "https"},
						SignatureVersions: []string{"v4"},
					},
				},
				IsRegionalized: true,
				Endpoints: Endpoints{
					EndpointKey{
						Region: "us-iso-east-1",
					}: {},
				},
			},
		},
		Cases: []testCase{
			{
				Region: "us-west-2",
				Expected: aws.Endpoint{
					URL:           "https://ec2.us-west-2.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-east-2",
				Expected: aws.Endpoint{
					URL:           "https://api.ec2.us-east-2.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-2",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-west-2",
				Options: Options{UseDualStack: true},
				Expected: aws.Endpoint{
					URL:           "https://ec2.us-west-2.api.aws",
					PartitionID:   "aws",
					SigningRegion: "us-west-2",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-east-2",
				Options: Options{UseDualStack: true},
				Expected: aws.Endpoint{
					URL:           "https://api.ec2.us-east-2.api.aws",
					PartitionID:   "aws",
					SigningRegion: "us-east-2",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-iso-east-1",
				Expected: aws.Endpoint{
					URL:           "https://ec2.us-iso-east-1.c2s.ic.gov",
					PartitionID:   "aws-iso",
					SigningRegion: "us-iso-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-iso-east-1",
				Options: Options{
					UseDualStack: true,
				},
				WantErr: true,
			},
		},
	},
	"route53": {
		Partitions: Partitions{
			{
				ID: "aws",
				Defaults: map[EndpointVariant]Endpoint{
					DualStackVariant: {
						Hostname:          "route53.{region}.api.aws",
						Protocols:         []string{"https"},
						SignatureVersions: []string{"v4"},
					},
					FIPSVariant: {
						Hostname:          "route53-fips.{region}.amazonaws.com",
						Protocols:         []string{"https"},
						SignatureVersions: []string{"v4"},
					},
					FIPSVariant | DualStackVariant: {
						Hostname:          "route53-fips.{region}.api.aws",
						Protocols:         []string{"https"},
						SignatureVersions: []string{"v4"},
					},
					0: {
						Hostname:          "route53.{region}.amazonaws.com",
						Protocols:         []string{"https"},
						SignatureVersions: []string{"v4"},
					},
				},
				RegionRegex:       regexp.MustCompile("^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$"),
				IsRegionalized:    false,
				PartitionEndpoint: "aws-global",
				Endpoints: map[EndpointKey]Endpoint{
					{
						Region: "aws-global",
					}: {
						Hostname: "route53.amazonaws.com",
						CredentialScope: CredentialScope{
							Region: "us-east-1",
						},
					},
					{
						Region:  "aws-global",
						Variant: DualStackVariant,
					}: {
						Hostname: "route53.global.api.aws",
						CredentialScope: CredentialScope{
							Region: "us-east-1",
						},
					},
					{
						Region:  "aws-global",
						Variant: FIPSVariant,
					}: {
						Hostname: "route53-fips.amazonaws.com",
						CredentialScope: CredentialScope{
							Region: "us-east-1",
						},
					},
					{
						Region: "other-thing",
					}: {
						Hostname: "other-thing.route53.amazonaws.com",
						CredentialScope: CredentialScope{
							Region: "us-east-1",
						},
					},
				},
			},
		},
		Cases: []testCase{
			{
				Region: "us-west-2",
				Expected: aws.Endpoint{
					URL:           "https://route53.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region: "us-east-2",
				Expected: aws.Endpoint{
					URL:           "https://route53.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-west-2",
				Options: Options{UseDualStack: true},
				Expected: aws.Endpoint{
					URL:           "https://route53.global.api.aws",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-east-2",
				Options: Options{UseDualStack: true},
				Expected: aws.Endpoint{
					URL:           "https://route53.global.api.aws",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-west-2",
				Options: Options{UseFIPS: true},
				Expected: aws.Endpoint{
					URL:           "https://route53-fips.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region:  "us-east-2",
				Options: Options{UseFIPS: true},
				Expected: aws.Endpoint{
					URL:           "https://route53-fips.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
			{
				Region: "other-thing",
				Expected: aws.Endpoint{
					URL:           "https://other-thing.route53.amazonaws.com",
					PartitionID:   "aws",
					SigningRegion: "us-east-1",
					SigningMethod: "v4",
				},
			},
		},
	},
}

func TestResolveEndpoint(t *testing.T) {
	for service := range testCases {
		t.Run(service, func(t *testing.T) {
			partitions := testCases[service].Partitions
			testCases := testCases[service].Cases

			for i, tt := range testCases {
				t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
					endpoint, err := partitions.ResolveEndpoint(tt.Region, tt.Options)
					if (err != nil) != (tt.WantErr) {
						t.Errorf("WantErr=%v, got error: %v", tt.WantErr, err)
					}
					if diff := cmp.Diff(tt.Expected, endpoint); len(diff) > 0 {
						t.Error(diff)
					}
				})
			}
		})
	}
}
