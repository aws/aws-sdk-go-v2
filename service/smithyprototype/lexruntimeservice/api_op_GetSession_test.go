package lexruntimeservice_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	lexruntime "github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

/*
/// Mock smithy tests
apply GetSession @httpRequestTests([
	{
		id: "RestJsonUriEncode",
		documentation: "serialize URL",
		protocol: "aws.rest-json-1.1",
		params: {
			"botName": "ABotName",
			"botAlias": "ABotAlias",
			"userId": "AUserId"
		},
		method: "GET",
		uri: "/bot/ABotName/alias/ABotAlias/user/AUserID/session",
		forbidHeaders: [
			"Content-Type"
		]
	}
	... Other tests
])
*/

func TestClient_GetSession_serialize_RestJsonUriEncode(t *testing.T) {
	t.Skip("middleware not implemented")

	cfg := unit.Config()

	var actualReq *http.Request
	client := lexruntime.NewFromConfig(cfg, func(o *lexruntime.Options) {
		o.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
			actualReq = r
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{},
			}, nil
		})
	})

	params := &lexruntime.GetSessionInput{
		BotName:  aws.String("ABotName"),
		BotAlias: aws.String("ABotAlias"),
		UserId:   aws.String("AUserId"),
	}

	_, err := client.GetSession(context.Background(), params)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if actualReq == nil {
		t.Fatalf("expect built request, got none")
	}

	// assert method
	if e, a := "GET", actualReq.Method; e != a {
		t.Errorf("expect %v method, got %v", e, a)
	}

	// assert uri
	if e, a := "/bot/ABotName/alias/ABotAlias/user/AUserID/session", actualReq.URL.Path; e != a {
		t.Errorf("expect %v url, got %v", e, a)
	}

	// assert forbibHeaders
	if v := actualReq.Header.Get("Content-Type"); len(v) != 0 {
		t.Errorf("expect no header value for Content-Type, got %v", v)
	}
}

/*
/// Mock smithy tests
apply GetSession @httpResponseTests([
	{
		id: "RestJsonError",
		documentation: "deserialize error response body",
        protocol: "aws.rest-json-1.1",
        code: 200,
        body: """
			{"SessionId": "some sessionID"}
             """,
        bodyMediaType: "application/json",
        headers: {
            "Content-Type": "application/json",
			"X-Amz-RequestId": "abc123"
        },
		vendorParams: {
			"RequestID": "abc123"
		},
        params: {
			SessionId: "some sessionID",
        }
	}
	... Other tests
])
*/

func TestClient_GetSession_deserialize_RestJsonDecode(t *testing.T) {
	t.Skip("middleware not implemented")

	cfg := unit.Config()

	client := lexruntime.NewFromConfig(cfg, func(o *lexruntime.Options) {
		o.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"Content-Type":    []string{"application/json"},
					"X-Amz-RequestId": []string{"abc123"},
				},
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"SessionId": "some sessionId"}`))),
			}, nil
		})
	})

	// TODO disable input parameter validation.
	params := &lexruntime.GetSessionInput{}

	result, err := client.GetSession(context.Background(), params)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// Assert params
	if result.SessionId == nil {
		t.Fatalf("expect value for SessionId, got none")
	}
	if e, a := `some sessionId`, *result.SessionId; e != a {
		t.Errorf("expect %v, sessionId, got %v", e, a)
	}

	// Assert Vendor Metadata
	if e, a := "abc123", getRequestID(result.ResultMetadata); e != a {
		t.Errorf("expect %v RequestID, got %v", e, a)
	}
}

func getRequestID(middleware.Metadata) string {
	// TODO replace with actual metadata helper
	return "abc123"
}
