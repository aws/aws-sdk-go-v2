package imds

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetRegion(t *testing.T) {
	cases := map[string]struct {
		Body         []byte
		ExpectRegion string
		ExpectTrace  []string
		ExpectErr    string
	}{
		"success": {
			Body:         []byte(instanceIdentityDocument),
			ExpectRegion: "us-east-1",
			ExpectTrace: []string{
				getTokenPath,
				getInstanceIdentityDocumentPath,
			},
		},
	}

	ctx := context.Background()

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			trace := newRequestTrace()
			server := httptest.NewServer(trace.WrapHandler(
				newTestServeMux(t,
					newSecureAPIHandler(t,
						[]string{"tokenA"},
						5*time.Minute,
						&successAPIResponseHandler{t: t,
							path:   getInstanceIdentityDocumentPath,
							method: "GET",
							body:   append([]byte{}, c.Body...),
						},
					))))
			defer server.Close()

			// Asserts
			client := New(Options{
				Endpoint: server.URL,
			})

			resp, err := client.GetRegion(ctx, nil)
			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if resp == nil {
				t.Fatalf("expect resp, got none")
			}

			if diff := cmp.Diff(c.ExpectRegion, resp.Region); len(diff) != 0 {
				t.Errorf("expect region to match\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}
		})
	}
}
