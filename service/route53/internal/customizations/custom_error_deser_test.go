package customizations_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func TestCustomErrorDeserialization(t *testing.T) {
	cases := map[string]struct {
		responseStatus     int
		responseBody       []byte
		expectedError      string
		expectedRequestID  string
		expectedResponseID string
	}{
		"invalidChangeBatchError": {
			responseStatus: 500,
			responseBody: []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InvalidChangeBatch xmlns="https://route53.amazonaws.com/doc/2013-04-01/">
		  <Messages>
		    <Message>Tried to create resource record set duplicate.example.com. type A, but it already exists</Message>
		  </Messages>
		  <RequestId>b25f48e8-84fd-11e6-80d9-574e0c4664cb</RequestId>
		</InvalidChangeBatch>`),
			expectedError:     "InvalidChangeBatch: ChangeBatch errors occurred",
			expectedRequestID: "b25f48e8-84fd-11e6-80d9-574e0c4664cb",
		},

		"standardRestXMLError": {
			responseStatus: 500,
			responseBody: []byte(`<?xml version="1.0"?>
		<ErrorResponse xmlns="http://route53.amazonaws.com/doc/2016-09-07/">
		  <Error>
		    <Type>Sender</Type>
		    <Code>MalformedXML</Code>
		    <Message>1 validation error detected: Value null at 'route53#ChangeSet' failed to satisfy constraint: Member must not be null</Message>
		  </Error>
		  <RequestId>b25f48e8-84fd-11e6-80d9-574e0c4664cb</RequestId>
		</ErrorResponse>
		`),
			expectedError:     "1 validation error detected:",
			expectedRequestID: "b25f48e8-84fd-11e6-80d9-574e0c4664cb",
		},

		"Success response": {
			responseStatus: 200,
			responseBody: []byte(`<?xml version="1.0" encoding="UTF-8"?>
		<ChangeResourceRecordSetsResponse>
   			<ChangeInfo>
      		<Comment>mockComment</Comment>
      		<Id>mockID</Id>
   		</ChangeInfo>
		</ChangeResourceRecordSetsResponse>`),
			expectedResponseID: "mockID",
		},
	}

	for name, c := range cases {
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.responseStatus)
				w.Write(c.responseBody)
			}))
		defer server.Close()

		t.Run(name, func(t *testing.T) {
			svc := route53.NewFromConfig(aws.Config{
				Region: "us-east-1",
				EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:         server.URL,
						SigningName: "route53",
					}, nil
				}),
				Retryer: func() aws.Retryer {
					return aws.NopRetryer{}
				},
				Credentials: &fakeCredentials{},
			})
			resp, err := svc.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
				ChangeBatch: &types.ChangeBatch{
					Changes: []types.Change{},
					Comment: aws.String("mock"),
				},
				HostedZoneId: aws.String("zone"),
			})

			if err == nil && len(c.expectedError) != 0 {
				t.Fatalf("expected err, got none")
			}

			if len(c.expectedError) != 0 {
				if e, a := c.expectedError, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expected error to be %s, got %s", e, a)
				}

				var responseError interface {
					ServiceRequestID() string
				}

				if !errors.As(err, &responseError) {
					t.Fatalf("expected error to be of type %T, was not", responseError)
				}

				if e, a := c.expectedRequestID, responseError.ServiceRequestID(); !strings.EqualFold(e, a) {
					t.Fatalf("expected request id to be %s, got %s", e, a)
				}
			}

			if len(c.expectedResponseID) != 0 {
				if e, a := c.expectedResponseID, *resp.ChangeInfo.Id; !strings.EqualFold(e, a) {
					t.Fatalf("expected response to have id %v, got %v", e, a)
				}
			}

		})
	}
}

type fakeCredentials struct{}

func (*fakeCredentials) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{}, nil
}
