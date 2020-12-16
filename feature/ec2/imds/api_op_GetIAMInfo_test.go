package imds

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetIAMInfo(t *testing.T) {
	const validIamInfo = `{
		"Code" : "Success",
		"LastUpdated" : "2016-03-17T12:27:32Z",
		"InstanceProfileArn" : "arn:aws:iam::123456789012:instance-profile/my-instance-profile",
		"InstanceProfileId" : "AIPAABCDEFGHIJKLMN123"
	}`

	const unsuccessfulIamInfo = `{
		"Code" : "Failed"
	}`

	cases := map[string]struct {
		Body         []byte
		ExpectResult IAMInfo
		ExpectTrace  []string
		ExpectErr    string
	}{
		"success": {
			Body: []byte(validIamInfo),
			ExpectResult: IAMInfo{
				Code:               "Success",
				LastUpdated:        time.Date(2016, 3, 17, 12, 27, 32, 0, time.UTC),
				InstanceProfileArn: "arn:aws:iam::123456789012:instance-profile/my-instance-profile",
				InstanceProfileID:  "AIPAABCDEFGHIJKLMN123",
			},
			ExpectTrace: []string{
				getTokenPath,
				getIAMInfoPath,
			},
		},
		"not success code": {
			Body:      []byte(unsuccessfulIamInfo),
			ExpectErr: "Failed",
			ExpectTrace: []string{
				getTokenPath,
				getIAMInfoPath,
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
							path:   getIAMInfoPath,
							method: "GET",
							body:   append([]byte{}, c.Body...),
						},
					))))
			defer server.Close()

			// Asserts
			client := New(Options{
				Endpoint: server.URL,
			})

			resp, err := client.GetIAMInfo(ctx, nil)
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

			if diff := cmp.Diff(c.ExpectResult, resp.IAMInfo); len(diff) != 0 {
				t.Errorf("expect result to match\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}
		})
	}
}
