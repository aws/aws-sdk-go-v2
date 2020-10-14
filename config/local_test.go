package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
)

func TestValidateLocalURL(t *testing.T) {
	origFn := lookupHostFn
	defer func() { lookupHostFn = origFn }()

	lookupHostFn = func(host string) ([]string, error) {
		m := map[string]struct {
			Addrs []string
			Err   error
		}{
			"localhost":       {Addrs: []string{"::1", "127.0.0.1"}},
			"actuallylocal":   {Addrs: []string{"127.0.0.2"}},
			"notlocal":        {Addrs: []string{"::1", "127.0.0.1", "192.168.1.10"}},
			"www.example.com": {Addrs: []string{"10.10.10.10"}},
		}

		h, ok := m[host]
		if !ok {
			return nil, fmt.Errorf("unknown host")
		}

		return h.Addrs, h.Err
	}

	cases := []struct {
		Host string
		Fail bool
	}{
		{"localhost", false},
		{"actuallylocal", false},
		{"127.0.0.1", false},
		{"127.1.1.1", false},
		{"[::1]", false},
		{"www.example.com", true},
		{"169.254.170.2", true},
	}

	restoreEnv := awstesting.StashEnv()
	defer awstesting.PopEnv(restoreEnv)

	for _, c := range cases {
		t.Run(c.Host, func(t *testing.T) {
			u := fmt.Sprintf("http://%s/abc/123", c.Host)

			err := validateLocalURL(u)
			if c.Fail {
				if err == nil {
					t.Fatalf("expect error, got none")
				} else {
					if e, a := "invalid endpoint host", err.Error(); !strings.Contains(a, e) {
						t.Errorf("expect %s to be in %s", e, a)
					}
				}
			} else {
				if err != nil {
					t.Fatalf("expect no error, got %v", err)
				}
			}
		})
	}
}
