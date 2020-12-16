package imds

import (
	"bytes"
	"context"
	"encoding/hex"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetUserData(t *testing.T) {
	cases := map[string]struct {
		RespStatusCode int
		ExpectContent  []byte
		ExpectTrace    []string
		ExpectErr      string
	}{
		"get data": {
			ExpectContent: []byte("success"),
			ExpectTrace: []string{
				getTokenPath,
				getUserDataPath,
			},
		},
		"get data error": {
			RespStatusCode: 400,
			ExpectTrace: []string{
				getTokenPath,
				getUserDataPath,
			},
			ExpectErr: "EC2 IMDS failed",
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
							path:       getUserDataPath,
							method:     "GET",
							statusCode: c.RespStatusCode,
							body:       append([]byte{}, c.ExpectContent...),
						},
					))))
			defer server.Close()

			// Asserts
			client := New(Options{
				Endpoint: server.URL,
			})

			resp, err := client.GetUserData(ctx, nil)
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

			actualContent, err := ioutil.ReadAll(resp.Content)
			if err != nil {
				t.Fatalf("expect to read content, got %v", err)
			}

			if e, a := c.ExpectContent, actualContent; !bytes.Equal(e, a) {
				t.Errorf("expect content to be equal\nexpect:\n%s\nactual:\n%s",
					hex.Dump(e), hex.Dump(a))
			}

			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}
		})
	}
}
