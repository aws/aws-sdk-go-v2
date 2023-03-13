package imds

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const instanceIdentityDocument = `{
	"devpayProductCodes" : ["abc","123"],
	"marketplaceProductCodes" : [ "1a2bc3" ],
	"availabilityZone" : "us-east-1d",
	"privateIp" : "10.158.112.84",
	"version" : "2010-08-31",
	"region" : "us-east-1",
	"instanceId" : "i-1234567890abcdef0",
	"billingProducts" : ["123"],
	"instanceType" : "t1.micro",
	"accountId" : "123456789012",
	"pendingTime" : "2015-11-19T16:32:11Z",
	"imageId" : "ami-5fb8c835",
	"kernelId" : "aki-919dcaf8",
	"ramdiskId" : "abc123",
	"architecture" : "x86_64"
}`

func TestGetInstanceIdentityDocument(t *testing.T) {

	cases := map[string]struct {
		Body         []byte
		ExpectResult InstanceIdentityDocument
		ExpectTrace  []string
		ExpectErr    string
	}{
		"success": {
			Body: []byte(instanceIdentityDocument),
			ExpectResult: InstanceIdentityDocument{
				DevpayProductCodes:      []string{"abc", "123"},
				MarketplaceProductCodes: []string{"1a2bc3"},
				AvailabilityZone:        "us-east-1d",
				PrivateIP:               "10.158.112.84",
				Version:                 "2010-08-31",
				Region:                  "us-east-1",
				InstanceID:              "i-1234567890abcdef0",
				BillingProducts:         []string{"123"},
				InstanceType:            "t1.micro",
				AccountID:               "123456789012",
				PendingTime:             time.Date(2015, 11, 19, 16, 32, 11, 0, time.UTC),
				ImageID:                 "ami-5fb8c835",
				KernelID:                "aki-919dcaf8",
				RamdiskID:               "abc123",
				Architecture:            "x86_64",
			},
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

			resp, err := client.GetInstanceIdentityDocument(ctx, nil)
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

			if diff := cmp.Diff(c.ExpectResult, resp.InstanceIdentityDocument); len(diff) != 0 {
				t.Errorf("expect result to match\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}
		})
	}
}
