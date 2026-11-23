package v4

import (
	"net/http"
	"testing"
)

func TestStripPort(t *testing.T) {
	cases := map[string]struct {
		hostport string
		expect   string
	}{
		"hostname with port": {
			hostport: "example.com:443",
			expect:   "example.com",
		},
		"ipv6 with default https port": {
			hostport: "[2001:db8::1]:443",
			expect:   "[2001:db8::1]",
		},
		"ipv6 loopback with default http port": {
			hostport: "[::1]:80",
			expect:   "[::1]",
		},
		"ipv6 without port": {
			hostport: "[2001:db8::1]",
			expect:   "[2001:db8::1]",
		},
		"ipv6 with zone identifier and port": {
			hostport: "[fe80::1%25eth0]:443",
			expect:   "[fe80::1%25eth0]",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if e, a := c.expect, stripPort(c.hostport); e != a {
				t.Errorf("expect %q, got %q", e, a)
			}
		})
	}
}

func TestSanitizeHostForHeader(t *testing.T) {
	cases := map[string]struct {
		url        string
		expectHost string
	}{
		"https ipv6 default port": {
			url:        "https://[2001:db8::1]:443/path",
			expectHost: "[2001:db8::1]",
		},
		"http ipv6 loopback default port": {
			url:        "http://[::1]:80/path",
			expectHost: "[::1]",
		},
		"https ipv6 non-default port": {
			url:        "https://[2001:db8::1]:8443/path",
			expectHost: "[2001:db8::1]:8443",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", c.url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			SanitizeHostForHeader(req)

			if e, a := c.expectHost, req.Host; e != a {
				t.Errorf("expect req.Host %q, got %q", e, a)
			}
		})
	}
}
