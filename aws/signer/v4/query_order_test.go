package v4

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	v4Internal "github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4"
)

// TestSign_canonicalQueryOrder asserts that the canonical query string orders
// parameters by their URI-encoded form. The SigV4 specification requires each
// name and value to be URI-encoded individually and states that "the sorting
// occurs after encoding":
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
		"repeated key, values diverging at '/' and '-', reversed input": {
			query:  url.Values{"k": {"a-b", "a/b"}},
			expect: "k=a%2Fb&k=a-b",
		},
		"repeated key, shared prefix": {
			query: url.Values{"filter[]": {
				"team/group/1",
				"team-scout/group/1",
			}},
			expect: "filter%5B%5D=team%2Fgroup%2F1&filter%5B%5D=team-scout%2Fgroup%2F1",
		},
		// Key axis: 'é' encodes to "%C3%A9", and '%' sorts before 'a'.
		"keys ordered on encoded form": {
			query:  url.Values{"a": {"1"}, "é": {"2"}},
			expect: "%C3%A9=2&a=1",
		},
		// Controls: values that order identically either way must not move.
		"values already in order": {
			query:  url.Values{"k": {"alpha", "zeta"}},
			expect: "k=alpha&k=zeta",
		},
		"repeated values still sorted within a key": {
			query:  url.Values{"b": {"1"}, "a": {"2", "1"}},
			expect: "a=1&a=2&b=1",
		},
		"space encodes to %20": {
			query:  url.Values{"k": {"a b"}},
			expect: "k=a%20b",
		},
		"literal plus is escaped": {
			query:  url.Values{"k": {"a+b"}},
			expect: "k=a%2Bb",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "https://service.us-east-1.amazonaws.com/", nil)
			req.URL.RawQuery = c.query.Encode()

			signer := &httpSigner{
				Request:     req,
				ServiceName: "service",
				Region:      "us-east-1",
				Credentials: testCredentials,
				Time: v4Internal.NewSigningTime(
					time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				KeyDerivator: v4Internal.NewSigningKeyDeriver(),
			}

			build, err := signer.Build()
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			// The canonical query string is the third line of the canonical request.
			lines := strings.Split(build.CanonicalString, "\n")
			if len(lines) < 3 {
				t.Fatalf("malformed canonical string:\n%s", build.CanonicalString)
			}
			if e, a := c.expect, lines[2]; e != a {
				t.Errorf("canonical query:\n  expect %v\n  actual %v", e, a)
			}

			// The query actually sent must match what was signed, otherwise the
			// service recomputes a different canonical request.
			if e, a := c.expect, req.URL.RawQuery; e != a {
				t.Errorf("wire query:\n  expect %v\n  actual %v", e, a)
			}
		})
	}
}
