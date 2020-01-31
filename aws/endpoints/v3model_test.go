package endpoints

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestUnmarshalRegionRegex(t *testing.T) {
	var input = []byte(`
{
    "regionRegex": "^(us|eu|ap|sa|ca)\\-\\w+\\-\\d+$"
}`)

	p := partition{}
	err := json.Unmarshal(input, &p)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	expectRegexp, err := regexp.Compile(`^(us|eu|ap|sa|ca)\-\w+\-\d+$`)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := expectRegexp.String(), p.RegionRegex.Regexp.String(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestUnmarshalRegion(t *testing.T) {
	var input = []byte(`
{
	"aws-global": {
	  "description": "AWS partition-global endpoint"
	},
	"us-east-1": {
	  "description": "US East (N. Virginia)"
	}
}`)

	rs := regions{}
	err := json.Unmarshal(input, &rs)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 2, len(rs); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	r, ok := rs["aws-global"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := "AWS partition-global endpoint", r.Description; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	r, ok = rs["us-east-1"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := "US East (N. Virginia)", r.Description; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestUnmarshalServices(t *testing.T) {
	var input = []byte(`
{
	"acm": {
	  "endpoints": {
		"us-east-1": {}
	  }
	},
	"apigateway": {
      "isRegionalized": true,
	  "endpoints": {
		"us-east-1": {},
        "us-west-2": {}
	  }
	},
	"notRegionalized": {
      "isRegionalized": false,
	  "endpoints": {
		"us-east-1": {},
        "us-west-2": {}
	  }
	}
}`)

	ss := services{}
	err := json.Unmarshal(input, &ss)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 3, len(ss); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	s, ok := ss["acm"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := 1, len(s.Endpoints); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	if e, a := boxedBoolUnset, s.IsRegionalized; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	s, ok = ss["apigateway"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := 2, len(s.Endpoints); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	if e, a := boxedTrue, s.IsRegionalized; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	s, ok = ss["notRegionalized"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := 2, len(s.Endpoints); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	if e, a := boxedFalse, s.IsRegionalized; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestUnmarshalEndpoints(t *testing.T) {
	var inputs = []byte(`
{
	"aws-global": {
	  "hostname": "cloudfront.amazonaws.com",
	  "protocols": [
		"http",
		"https"
	  ],
	  "signatureVersions": [ "v4" ],
	  "credentialScope": {
		"region": "us-east-1",
		"service": "serviceName"
	  },
	  "sslCommonName": "commonName"
	},
	"us-east-1": {}
}`)

	es := endpoints{}
	err := json.Unmarshal(inputs, &es)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 2, len(es); e != a {
		t.Errorf("expect %v len, got %v", e, a)
	}
	s, ok := es["aws-global"]
	if !ok {
		t.Errorf("expect true")
	}
	if e, a := "cloudfront.amazonaws.com", s.Hostname; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := []string{"http", "https"}, s.Protocols; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := []string{"v4"}, s.SignatureVersions; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := (credentialScope{"us-east-1", "serviceName"}), s.CredentialScope; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "commonName", s.SSLCommonName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEndpointResolve(t *testing.T) {
	defs := []endpoint{
		{
			Hostname:          "{service}.{region}.{dnsSuffix}",
			SignatureVersions: []string{"v2"},
			SSLCommonName:     "sslCommonName",
		},
		{
			Hostname:  "other-hostname",
			Protocols: []string{"http"},
			CredentialScope: credentialScope{
				Region:  "signing_region",
				Service: "signing_service",
			},
		},
	}

	e := endpoint{
		Hostname:          "{service}.{region}.{dnsSuffix}",
		Protocols:         []string{"http", "https"},
		SignatureVersions: []string{"v4"},
		SSLCommonName:     "new sslCommonName",
	}

	resolved := e.resolve("service", "partition", "region", "dnsSuffix", defs, ResolveOptions{})

	if e, a := "https://service.region.dnsSuffix", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "partition", resolved.PartitionID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "signing_service", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "signing_region", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "v4", resolved.SigningMethod; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestEndpointMergeIn(t *testing.T) {
	expected := endpoint{
		Hostname:          "other hostname",
		Protocols:         []string{"http"},
		SignatureVersions: []string{"v4"},
		SSLCommonName:     "ssl common name",
		CredentialScope: credentialScope{
			Region:  "region",
			Service: "service",
		},
	}

	actual := endpoint{}
	actual.mergeIn(endpoint{
		Hostname:          "other hostname",
		Protocols:         []string{"http"},
		SignatureVersions: []string{"v4"},
		SSLCommonName:     "ssl common name",
		CredentialScope: credentialScope{
			Region:  "region",
			Service: "service",
		},
	})

	if e, a := expected, actual; !reflect.DeepEqual(e, a) {
		t.Errorf("expect %v, got %v", e, a)
	}
}

var testPartitions = partitions{
	partition{
		ID:        "part-id",
		Name:      "partitionName",
		DNSSuffix: "amazonaws.com",
		RegionRegex: regionRegex{
			Regexp: func() *regexp.Regexp {
				reg, _ := regexp.Compile("^(us|eu|ap|sa|ca)\\-\\w+\\-\\d+$")
				return reg
			}(),
		},
		Defaults: endpoint{
			Hostname:          "{service}.{region}.{dnsSuffix}",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		Regions: regions{
			"us-east-1": region{
				Description: "region description",
			},
			"us-west-2": region{},
		},
		Services: services{
			"s3": service{},
			"service1": service{
				Defaults: endpoint{
					CredentialScope: credentialScope{
						Service: "service1",
					},
				},
				Endpoints: endpoints{
					"us-east-1": {},
					"us-west-2": {
						HasDualStack:      boxedTrue,
						DualStackHostname: "{service}.dualstack.{region}.{dnsSuffix}",
					},
				},
			},
			"service2": service{
				Defaults: endpoint{
					CredentialScope: credentialScope{
						Service: "service2",
					},
				},
			},
			"httpService": service{
				Defaults: endpoint{
					Protocols: []string{"http"},
				},
			},
			"globalService": service{
				IsRegionalized:    boxedFalse,
				PartitionEndpoint: "aws-global",
				Endpoints: endpoints{
					"aws-global": endpoint{
						CredentialScope: credentialScope{
							Region: "us-east-1",
						},
						Hostname: "globalService.amazonaws.com",
					},
				},
			},
		},
	},
}

func TestResolveEndpoint(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("service2", "us-west-2", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://service2.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service2", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if resolved.SigningNameDerived {
		t.Errorf("expect the signing name not to be derived, but was")
	}
}

func TestResolveEndpoint_DisableSSL(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("service2", "us-west-2", ResolveOptions{DisableSSL: true})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "http://service2.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service2", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if resolved.SigningNameDerived {
		t.Errorf("expect the signing name not to be derived, but was")
	}
}

func TestResolveEndpoint_UseDualStack(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("service1", "us-west-2", ResolveOptions{UseDualStack: true})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://service1.dualstack.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service1", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if resolved.SigningNameDerived {
		t.Errorf("expect the signing name not to be derived, but was")
	}
}

func TestResolveEndpoint_HTTPProtocol(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("httpService", "us-west-2", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "http://httpService.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "httpService", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if !resolved.SigningNameDerived {
		t.Errorf("expect the signing name to be derived")
	}
}

func TestResolveEndpoint_UnknownService(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("unknownservice", "us-west-2", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://unknownservice.us-west-2.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-west-2", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "unknownservice", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if !resolved.SigningNameDerived {
		t.Errorf("expect the signing name to be derived")
	}
}

func TestResolveEndpoint_UnknownMatchedRegion(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("service2", "us-region-1", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://service2.us-region-1.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-region-1", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service2", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if resolved.SigningNameDerived {
		t.Errorf("expect the signing name not to be derived, but was")
	}
}

func TestResolveEndpoint_UnknownRegion(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("service2", "unknownregion", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://service2.unknownregion.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "unknownregion", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "service2", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if resolved.SigningNameDerived {
		t.Errorf("expect the signing name not to be derived, but was")
	}
}

func TestResolveEndpoint_StrictPartitionUnknownEndpoint(t *testing.T) {
	_, err := testPartitions[0].EndpointFor("service2", "unknownregion", ResolveOptions{StrictMatching: true})
	if err == nil {
		t.Fatalf("expect error")
	}

	_, ok := err.(UnknownEndpointError)
	if !ok {
		t.Errorf("expect error to be UnknownEndpointError")
	}
}

func TestResolveEndpoint_StrictPartitionsUnknownEndpoint(t *testing.T) {
	_, err := testPartitions.EndpointFor("service2", "us-region-1", ResolveOptions{StrictMatching: true})
	if err == nil {
		t.Fatalf("expect error")
	}

	_, ok := err.(UnknownEndpointError)
	if !ok {
		t.Errorf("expect error to be UnknownEndpointError")
	}
}

func TestResolveEndpoint_NotRegionalized(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("globalService", "us-west-2", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://globalService.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-east-1", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "globalService", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if !resolved.SigningNameDerived {
		t.Errorf("expect the signing name to be derived")
	}
}

func TestResolveEndpoint_AwsGlobal(t *testing.T) {
	resolved, err := testPartitions.EndpointFor("globalService", "aws-global", ResolveOptions{})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "https://globalService.amazonaws.com", resolved.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-east-1", resolved.SigningRegion; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "globalService", resolved.SigningName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if !resolved.SigningNameDerived {
		t.Errorf("expect the signing name to be derived")
	}
}

func TestEndpointFor_EmptyRegion(t *testing.T) {
	cases := map[string]struct {
		Service    string
		Region     string
		RealRegion string
		ExpectErr  string
	}{
		// Legacy services that previous accepted empty region
		"ec2metadata": {Service: "ec2metadata", RealRegion: "aws-global"},

		// Other services
		"s3":           {Service: "s3", Region: "us-east-1", RealRegion: "us-east-1"},
		"s3 no region": {Service: "s3", ExpectErr: "could not resolve endpoint"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual, err := NewDefaultResolver().ResolveEndpoint(c.Service, c.Region)
			if len(c.ExpectErr) != 0 {
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Errorf("expect %q error in %q", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error got, %v", err)
			}

			expect, err := NewDefaultResolver().ResolveEndpoint(c.Service, c.RealRegion)
			if err != nil {
				t.Fatalf("failed to get endpoint for default resolver")
			}
			if e, a := expect.URL, actual.URL; e != a {
				t.Errorf("expect %v URL, got %v", e, a)
			}
			if e, a := expect.SigningRegion, actual.SigningRegion; e != a {
				t.Errorf("expect %v signing region, got %v", e, a)
			}

		})
	}
}
