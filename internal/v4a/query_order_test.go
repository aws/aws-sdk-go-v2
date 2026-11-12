package v4a

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestSign_canonicalQueryOrder asserts that the SigV4a canonical query string
// orders parameters by their URI-encoded form. The specification requires each
// name and value to be URI-encoded individually, with the sorting occurring
// after encoding:
// https://docs.aws.amazon.com/IAM/latest/UserGuide/create-signed-request.html
func TestSign_canonicalQueryOrder(t *testing.T) {
	cases := map[string]struct {
		query  url.Values
		expect string
	}{
		// '/' (0x2F) sorts after '-' (0x2D), but once encoded '%' (0x25)
		// sorts before it, so the two orders disagree.
		"repeated key, values diverging at '/' and '-'": {
			query:  url.Values{"k": {"a/b", "a-b"}},
			expect: "k=a%2Fb&k=a-b",
		},
		"keys ordered on encoded form": {
			query:  url.Values{"a": {"1"}, "é": {"2"}},
			expect: "%C3%A9=2&a=1",
		},
		"values already in order": {
			query:  url.Values{"k": {"alpha", "zeta"}},
			expect: "k=alpha&k=zeta",
		},
		"repeated values still sorted within a key": {
			query:  url.Values{"b": {"1"}, "a": {"2", "1"}},
			expect: "a=1&a=2&b=1",
		},
	}

	creds, err := stubCredentials.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "https://service.us-east-1.amazonaws.com/", nil)
			req.URL.RawQuery = c.query.Encode()

			h := sha256.Sum256(nil)
			signer := NewSigner()
			err := signer.SignHTTP(context.Background(), creds, req,
				hex.EncodeToString(h[:]), "service", []string{"us-east-1"}, time.Now())
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			// The query sent on the wire must match what was signed, otherwise
			// the service recomputes a different canonical request.
			if e, a := c.expect, req.URL.RawQuery; e != a {
				t.Errorf("wire query:\n  expect %v\n  actual %v", e, a)
			}
			if strings.Contains(req.URL.RawQuery, "+") {
				t.Errorf("space must encode as %%20, got %v", req.URL.RawQuery)
			}
		})
	}
}
