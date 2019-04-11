// +build go1.8,codegen

package api

import (
	"testing"
)

func TestAPI_ServiceID(t *testing.T) {
	cases := map[string]struct {
		Metadata Metadata
		Expect   string
	}{
		"FullName": {
			Metadata: Metadata{
				ServiceFullName: "Amazon Service Name-100",
			},
			Expect: "ServiceName100",
		},
		"Abbreviation": {
			Metadata: Metadata{
				ServiceFullName:     "Amazon Service Name-100",
				ServiceAbbreviation: "AWS SN100",
			},
			Expect: "SN100",
		},
		"Lowercase Name": {
			Metadata: Metadata{
				EndpointPrefix:      "other",
				ServiceFullName:     "AWS Lowercase service",
				ServiceAbbreviation: "lowercase",
			},
			Expect: "Lowercase",
		},
		"Lowercase Name Mixed": {
			Metadata: Metadata{
				EndpointPrefix:      "other",
				ServiceFullName:     "AWS Lowercase service",
				ServiceAbbreviation: "lowercase name Goes heRe",
			},
			Expect: "LowercaseNameGoesHeRe",
		},
	}

	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			a := API{
				Metadata: c.Metadata,
			}

			if e, o := c.Expect, a.ServiceID(); e != o {
				t.Errorf("expect %v structName, got %v", e, o)
			}
		})
	}
}
