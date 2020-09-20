package ec2imds

import (
	"bytes"
	"context"
	"encoding/hex"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetDynamicData(t *testing.T) {
	cases := map[string]struct {
		Path          string
		ExpectContent []byte
		ExpectTrace   []string
	}{
		"empty path": {
			ExpectContent: []byte("success"),
			ExpectTrace: []string{
				getTokenPath,
				getDynamicDataPath,
			},
		},
		"with path": {
			Path:          "/abc",
			ExpectContent: []byte("success"),
			ExpectTrace: []string{
				getTokenPath,
				getDynamicDataPath + "/abc",
			},
		},
		"with path trailing slash": {
			Path:          "/abc/",
			ExpectContent: []byte("success"),
			ExpectTrace: []string{
				getTokenPath,
				getDynamicDataPath + "/abc/",
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
							path:   getDynamicDataPath + c.Path,
							method: "GET",
							body:   append([]byte{}, c.ExpectContent...),
						},
					))))
			defer server.Close()

			// Asserts
			client := New(Options{
				Endpoint: server.URL,
			})

			resp, err := client.GetDynamicData(ctx, &GetDynamicDataInput{
				Path: c.Path,
			})
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
