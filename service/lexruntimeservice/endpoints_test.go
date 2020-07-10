/*
# Considerations
* Could we eliminate pseudo-regions from the SDK?
    * Do we need to support backwards compatibility with pseudo-regions and "clean" them if read in from ~/.aws/config
    * Backwards compatibility would effectively duplicate logic between Java endpoint codegen and Go
* Still need to support -global variants
* Could the resolver be configurable with grab bag attributes to enable subset of features
	* FIPS
	* DualStack
	* S3 Accelerate
* What is the desired behavior when a region is not modeled but endpoint attributes are specified?
    * Should these fail? Should we attempt to define a non-ratified service URL format?
* Java V2 just has logic in it's region validation to allow matching of regions that either have fips- or -fips prefix
or suffix when validating with a regex by effectively stream those components off
*/

package lexruntimeservice

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

var serviceName string

type EndpointResolver interface {
	ResolveEndpoint(region string, options ...func(*ResolverOptions)) (aws.Endpoint, error)
}

func NewDefaultResolver(options ...func(options *ResolverOptions)) AWSEndpointProvider {
	opts := &ResolverOptions{}

	for _, fn := range options {
		fn(opts)
	}

	return AWSEndpointProvider{
		options:    *opts,
		partitions: defaultPartitions,
	}
}

type ResolverOptions struct {
	Constraints  *Capabilities
	DisableHTTPS bool
}

type AWSEndpointProvider struct {
	options    ResolverOptions
	partitions partitions
}

func (d AWSEndpointProvider) ResolveEndpoint(region string, options ...func(o *ResolverOptions)) (aws.Endpoint, error) {
	ro := d.options
	for _, fn := range options {
		fn(&ro)
	}
	return d.partitions.ResolveEndpoint(region, ro)
}

var defaultPartitions partitions

func init() {
	defaultPartitions = partitions{
		partition{
			ID: "aws",
			RegionRegex: func() *regexp.Regexp {
				reg, _ := regexp.Compile("^(us|eu|ap|sa|ca|me)\\-\\w+\\-\\d+$")
				return reg
			}(),
			Defaults: []endpoint{
				// Default from endpoints.json
				{
					Hostname:    "runtime.lex.{region}.amazonaws.com",
					Protocols:   []string{"https"},
					SigningName: "lex",
				},
				// Optional inferred template with capabilities
				{
					Hostname:  "runtime-fips.lex.{region}.amazonaws.com",
					Protocols: []string{"https"},
					InferredCapabilities: &Capabilities{
						FIPS: true,
					},
					SigningName: "lex",
				},
			},
			Endpoints: map[string]endpoint{
				"us-west-2": {
					Hostname:  "runtime.lex.us-west-2.amazonaws.com",
					Protocols: []string{"https"},
				},
				"us-west-2-fips": {
					Hostname:  "runtime-fips.lex.us-west-2.amazonaws.com",
					Protocols: []string{"https"},
					InferredCapabilities: &Capabilities{
						FIPS: true,
					},
					SigningRegion: "us-west-2",
				},
				"eu-west-1": {},
				"eu-west-2": {
					Unresolvable: aws.TrueTernary,
				},
			},
		},
		partition{
			ID: "aws-cn",
			RegionRegex: func() *regexp.Regexp {
				reg, _ := regexp.Compile("^(cn)\\-\\w+\\-\\d+$")
				return reg
			}(),
			isRegionalized:    aws.FalseTernary,
			PartitionEndpoint: "aws-cn-global",
			Defaults: []endpoint{
				{
					SigningName: "lex",
				},
			},
			Endpoints: map[string]endpoint{
				"aws-cn-global": {
					Hostname:      "runtime.lex.amazonaws.cn",
					Protocols:     []string{"https"},
					SigningRegion: "cn-central-1",
				},
			},
		},
	}
}

type partitions []partition

/*
How to to resolve a specified ConfigRegion identifier. As context ConfigRegion == (Region || PseudoRegion)
* If ConfigRegion identifier known
	* If No Endpoint Cap & No Resolver Cap Constraints => Valid
	* If Endpoint Cap & No Resolver Cap Constraints => Valid
	* If Endpoint Cap & Resolver Cap Constraints Match => Valid
	* If Endpoint Cap & Resolver Cap Constraints Don't Match => Invalid, fallback if to a default matching the resolver constraints or error
* Else if ConfigRegion identifier unknown
	* If No Resolver Cap Constraints => fallback to the partition default
	* If Resolver Cap constraints
		* Default endpoint for partition that matches caps => Valid
		* **Could** have a list of aliases of signing regions to pseudo-regions (excluding global) that could be then used to find an pseudo-alias that satisfies
		* Else => Invalid, error
*/
func (p partitions) ResolveEndpoint(region string, ro ResolverOptions) (aws.Endpoint, error) {
	for _, partition := range p {
		if isResolvable := partition.IsResolvable(region); !isResolvable {
			continue
		}
		return partition.ResolveEndpoint(region, ro)
	}

	if len(p) == 0 {
		return aws.Endpoint{}, fmt.Errorf("missing service endpoint metadata")
	}

	return p[0].ResolveEndpoint(region, ro)
}

type partition struct {
	ID                string
	isRegionalized    aws.Ternary
	PartitionEndpoint string
	RegionRegex       *regexp.Regexp
	Defaults          []endpoint
	Endpoints         map[string]endpoint
}

func (p *partition) IsResolvable(region string) bool {
	_, ok := p.Endpoints[region]
	return ok || p.RegionRegex.MatchString(region)
}

func (p *partition) ResolveEndpoint(region string, options ResolverOptions) (aws.Endpoint, error) {
	var merged endpoint

	// First merge any defaults which match our constraints
	for _, def := range p.Defaults {
		if options.Constraints.Match(def.InferredCapabilities) {
			merged.MergeIn(def)
			break
		}
	}

	var e endpoint
	var isKnown bool
	if p.isRegionalized == aws.FalseTernary {
		e, isKnown = p.Endpoints[p.PartitionEndpoint]
	} else {
		e, isKnown = p.Endpoints[region]
	}
	// First merge in matching region identifier information IF we have no capability constraints OR our constraints match exactly
	if isKnown && (options.Constraints == nil || options.Constraints.Match(e.InferredCapabilities)) {
		merged.MergeIn(e)
	}

	if merged.Unresolvable == aws.TrueTernary {
		return aws.Endpoint{}, fmt.Errorf("endpoint must be specified manually")
	}

	// Check that the fully merged endpoint matches constraints
	if options.Constraints != nil && !options.Constraints.Match(merged.InferredCapabilities) {
		return aws.Endpoint{}, fmt.Errorf("failed to resolve endpoint with matching capabilities")
	}

	if len(merged.Hostname) == 0 {
		return aws.Endpoint{}, fmt.Errorf("failed resolve endpoint")
	}

	signingRegion := merged.SigningRegion
	if len(signingRegion) == 0 {
		signingRegion = region
	}

	signingName := merged.SigningName
	var signingNameDerived bool
	if len(signingName) == 0 {
		signingName = serviceName
		signingNameDerived = true
	}

	merged.Hostname = strings.Replace(merged.Hostname, "{region}", signingRegion, 1)

	resolved := aws.Endpoint{
		URL:                fmt.Sprintf("%s://%s", getEndpointScheme(merged.Protocols, options), merged.Hostname),
		SigningName:        signingName,
		SigningRegion:      signingRegion,
		SigningNameDerived: signingNameDerived,
	}

	if merged.InferredCapabilities != nil {
		setInferredCapabilities(&resolved.Metadata, *merged.InferredCapabilities)
	}

	return resolved, nil
}

const defaultProtocol = "https"

var protocolPriority = []string{"https", "http"}

func getByPriority(s []string, p []string, def string) string {
	if len(s) == 0 {
		return def
	}

	for i := 0; i < len(p); i++ {
		for j := 0; j < len(s); j++ {
			if s[j] == p[i] {
				return s[j]
			}
		}
	}

	return s[0]
}

func getEndpointScheme(protocols []string, ro ResolverOptions) string {
	if ro.DisableHTTPS {
		return "http"
	}

	return getByPriority(protocols, protocolPriority, defaultProtocol)
}

type endpoint struct {
	Unresolvable         aws.Ternary
	Hostname             string
	Protocols            []string
	SigningName          string
	SigningRegion        string
	InferredCapabilities *Capabilities
}

func (e *endpoint) MergeIn(other endpoint) {
	if other.Unresolvable != aws.UnknownTernary {
		e.Unresolvable = other.Unresolvable
	}
	if len(other.Hostname) > 0 {
		e.Hostname = other.Hostname
	}
	if len(other.Protocols) > 0 {
		e.Protocols = other.Protocols
	}
	if len(other.SigningName) > 0 {
		e.SigningName = other.SigningName
	}
	if len(other.SigningRegion) > 0 {
		e.SigningRegion = other.SigningRegion
	}
	if other.InferredCapabilities != nil {
		e.InferredCapabilities = other.InferredCapabilities
	}
}

type Capabilities struct {
	FIPS      bool
	DualStack bool
}

func (c *Capabilities) Match(o *Capabilities) bool {
	if c == nil {
		return o == nil
	}
	if o == nil {
		return false
	}
	return *c == *o
}

type inferredCapabilities struct{}

func setInferredCapabilities(metadata *aws.EndpointMetadata, capabilities Capabilities) {
	metadata.Set(inferredCapabilities{}, capabilities)
}

func getInferredCapabilities(metadata aws.EndpointMetadata) (Capabilities, bool) {
	v, ok := metadata.Get(inferredCapabilities{}).(Capabilities)
	return v, ok
}

func TestNewDefaultResolver(t *testing.T) {
	cases := []struct {
		Region      string
		Options     ResolverOptions
		Expected    aws.Endpoint
		ExpectedErr string
	}{
		{
			Region: "us-west-2",
			Expected: aws.Endpoint{
				URL:           "https://runtime.lex.us-west-2.amazonaws.com",
				SigningName:   "lex",
				SigningRegion: "us-west-2",
			},
		},
		{
			Region: "us-west-2",
			Options: ResolverOptions{
				Constraints: &Capabilities{
					FIPS: true,
				},
			},
			Expected: aws.Endpoint{
				URL:           "https://runtime-fips.lex.us-west-2.amazonaws.com",
				SigningName:   "lex",
				SigningRegion: "us-west-2",
				Metadata: func() aws.EndpointMetadata {
					m := aws.EndpointMetadata{}
					setInferredCapabilities(&m, Capabilities{FIPS: true})
					return m
				}(),
			},
		},
		{
			Region: "us-west-2-fips",
			Expected: aws.Endpoint{
				URL:           "https://runtime-fips.lex.us-west-2.amazonaws.com",
				SigningName:   "lex",
				SigningRegion: "us-west-2",
				Metadata: func() aws.EndpointMetadata {
					m := aws.EndpointMetadata{}
					setInferredCapabilities(&m, Capabilities{FIPS: true})
					return m
				}(),
			},
		},
		{
			Region: "eu-west-1",
			Expected: aws.Endpoint{
				URL:           "https://runtime.lex.eu-west-1.amazonaws.com",
				SigningName:   "lex",
				SigningRegion: "eu-west-1",
			},
		},
		{
			Region:  "eu-west-1",
			Options: ResolverOptions{Constraints: &Capabilities{FIPS: true}},
			Expected: aws.Endpoint{
				URL:           "https://runtime-fips.lex.eu-west-1.amazonaws.com",
				SigningName:   "lex",
				SigningRegion: "eu-west-1",
				Metadata: func() aws.EndpointMetadata {
					m := aws.EndpointMetadata{}
					setInferredCapabilities(&m, Capabilities{FIPS: true})
					return m
				}(),
			},
		},
		{
			Region:      "eu-west-1",
			Options:     ResolverOptions{Constraints: &Capabilities{FIPS: true, DualStack: true}},
			ExpectedErr: "failed to resolve endpoint with matching capabilities",
		},
		{
			Region:      "eu-west-2",
			ExpectedErr: "endpoint must be specified manually",
		},
		{
			Region: "cn-north-1",
			Expected: aws.Endpoint{
				URL:           "https://runtime.lex.amazonaws.cn",
				SigningName:   "lex",
				SigningRegion: "cn-central-1",
			},
		},
		{
			Region:      "cn-north-1",
			Options:     ResolverOptions{Constraints: &Capabilities{FIPS: true}},
			ExpectedErr: "failed to resolve endpoint with matching capabilities",
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resolver := NewDefaultResolver(func(options *ResolverOptions) {
				*options = tt.Options
			})
			resolveEndpoint, err := resolver.ResolveEndpoint(tt.Region)
			if err != nil {
				if len(tt.ExpectedErr) == 0 {
					t.Fatalf("expected no error, got %v", err)
				}
				if e, a := tt.ExpectedErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected %v, got %v", e, a)
				}
			} else if len(tt.ExpectedErr) != 0 {
				t.Fatalf("expected error, got none")
			}

			if e, a := tt.Expected, resolveEndpoint; !reflect.DeepEqual(e, a) {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}
