package imds

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetToken(t *testing.T) {
	cases := map[string]struct {
		TokenTTL       time.Duration
		Header         http.Header
		Body           []byte
		ExpectToken    string
		ExpectTokenTTL time.Duration
		ExpectTrace    []string
		ExpectErr      string
	}{
		"success": {
			TokenTTL: 10 * time.Second,
			Header: http.Header{
				tokenTTLHeader: []string{"10"},
			},
			Body:           []byte("tokenABC"),
			ExpectToken:    "tokenABC",
			ExpectTokenTTL: 10 * time.Second,
			ExpectTrace: []string{
				getTokenPath,
			},
		},
	}

	ctx := context.Background()

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			trace := newRequestTrace()
			server := httptest.NewServer(trace.WrapHandler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					actualTTL := r.Header.Get(tokenTTLHeader)
					expectTTL := strconv.Itoa(int(c.TokenTTL / time.Second))
					if expectTTL != actualTTL {
						t.Errorf("expect %v token TTL request header, got %v",
							expectTTL, actualTTL)
						http.Error(w, http.StatusText(400), 400)
						return
					}

					(&successAPIResponseHandler{t: t,
						path:   getTokenPath,
						method: "PUT",
						header: c.Header,
						body:   append([]byte{}, c.Body...),
					}).ServeHTTP(w, r)
				})))
			defer server.Close()

			// Asserts
			client := New(Options{
				Endpoint: server.URL,
			})

			resp, err := client.getToken(ctx, &getTokenInput{
				TokenTTL: c.TokenTTL,
			})
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

			if e, a := c.ExpectToken, resp.Token; e != a {
				t.Errorf("expect %v token, got %v", e, a)
			}
			if e, a := c.ExpectTokenTTL, resp.TokenTTL; e != a {
				t.Errorf("expect %v token TTL, got %v", e, a)
			}

			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}
		})
	}
}
