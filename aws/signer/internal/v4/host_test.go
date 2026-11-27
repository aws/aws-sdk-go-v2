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
		"empty": {
			hostport: "",
			expect:   "",
		},
		"hostname without port": {
			hostport: "example.com",
			expect:   "example.com",
		},
		"hostname with port": {
			hostport: "example.com:443",
			expect:   "example.com",
		},
		"ipv4 without port": {
			hostport: "127.0.0.1",
			expect:   "127.0.0.1",
		},
		"ipv4 with port": {
			hostport: "127.0.0.1:80",
			expect:   "127.0.0.1",
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
		hostHeader string
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
		"http ipv6 non-default port": {
			url:        "http://[::1]:8080/path",
			expectHost: "[::1]:8080",
		},
		"https ipv6 without port": {
			url:        "https://[2001:db8::1]/path",
			expectHost: "[2001:db8::1]",
		},
		"https ipv6 explicit Host header with default port": {
			url:        "https://example.com/path",
			hostHeader: "[2001:db8::1]:443",
			expectHost: "[2001:db8::1]",
		},
		"https ipv6 with zone identifier and default port": {
			url:        "https://[fe80::1%25eth0]:443/path",
			expectHost: "[fe80::1%eth0]",
		},
		"https hostname default port": {
			url:        "https://example.com:443/path",
			expectHost: "example.com",
		},
		"http hostname default port": {
			url:        "http://example.com:80/path",
			expectHost: "example.com",
		},
		"https hostname non-default port": {
			url:        "https://example.com:8443/path",
			expectHost: "example.com:8443",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", c.url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			if c.hostHeader != "" {
				req.Host = c.hostHeader
			}

			SanitizeHostForHeader(req)

			if e, a := c.expectHost, req.Host; e != a {
				t.Errorf("expect req.Host %q, got %q", e, a)
			}
		})
	}
}
