package config

import (
	"reflect"
	"testing"
)

func TestEndpointMode_SetFromString(t *testing.T) {
	cases := map[string]struct {
		Value   string
		Expect  EndpointModeState
		WantErr bool
	}{
		"empty value": {
			Expect: EndpointModeStateUnset,
		},
		"unknown value": {
			Value:   "foobar",
			WantErr: true,
		},
		"IPv4": {
			Value:  "IPv4",
			Expect: EndpointModeStateIPv4,
		},
		"IPv6": {
			Value:  "IPv6",
			Expect: EndpointModeStateIPv6,
		},
		"IPv4 case-insensitive": {
			Value:  "iPv4",
			Expect: EndpointModeStateIPv4,
		},
		"IPv6 case-insensitive": {
			Value:  "iPv6",
			Expect: EndpointModeStateIPv6,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			var em EndpointModeState
			if err := em.SetFromString(tt.Value); (err != nil) != tt.WantErr {
				t.Fatalf("WantErr=%v, got err=%v", tt.WantErr, err)
			}
			if !reflect.DeepEqual(tt.Expect, em) {
				t.Errorf("%v != %v", tt.Expect, em)
			}
		})
	}
}
