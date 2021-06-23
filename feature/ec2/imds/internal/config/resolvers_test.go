package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEndpointMode_SetFromString(t *testing.T) {
	cases := map[string]struct {
		Value   string
		Expect  EndpointMode
		WantErr bool
	}{
		"empty value": {
			Expect: EndpointModeUnset,
		},
		"unknown value": {
			Value:   "foobar",
			WantErr: true,
		},
		"IPv4": {
			Value:  "IPv4",
			Expect: EndpointModeIPv4,
		},
		"IPv6": {
			Value:  "IPv6",
			Expect: EndpointModeIPv6,
		},
		"IPv4 case-insensitive": {
			Value:  "iPv4",
			Expect: EndpointModeIPv4,
		},
		"IPv6 case-insensitive": {
			Value:  "iPv6",
			Expect: EndpointModeIPv6,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var em EndpointMode
			if err := em.SetFromString(tt.Value); (err != nil) != tt.WantErr {
				t.Fatalf("WantErr=%v, got err=%v", tt.WantErr, err)
			}
			if diff := cmp.Diff(em, tt.Expect); len(diff) > 0 {
				t.Errorf(diff)
			}
		})
	}
}
