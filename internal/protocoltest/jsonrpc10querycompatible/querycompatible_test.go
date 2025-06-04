package jsonrpc10querycompatible

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	smithy "github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyprivateprotocol "github.com/aws/smithy-go/private/protocol"
)

// Implements awsQuery-compatible SEP tests
// (query-protocol-migration-compatibility.md)

type mockHTTP struct {
	QueryError            string
	BodyType, BodyMessage string
}

func (m mockHTTP) Do(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 400,
		Header: http.Header{
			"X-Amzn-Query-Error": {m.QueryError},
		},
		Body: io.NopCloser(strings.NewReader(fmt.Sprintf(`{
  "__type": %q,
  "message": %q
}`, m.BodyType, m.BodyMessage))),
	}, nil
}

// TC1: Validate SDK can parse Code field
func TestQueryCompatible_SEP1(t *testing.T) {
	svc := New(Options{
		HTTPClient: mockHTTP{
			QueryError: "AWS.SimpleQueueService.NonExistentQueue;Sender",
			BodyType:   "aws.protocoltests.json#ItemNotFound",
		},
	})
	_, err := svc.GetItem(context.Background(), &GetItemInput{})
	if err == nil {
		t.Fatal("expect err, got none")
	}

	var terr smithy.APIError
	if !errors.As(err, &terr) {
		t.Fatalf("expect smithy.APIError, got %T", err)
	}

	expect := "AWS.SimpleQueueService.NonExistentQueue"
	if actual := terr.ErrorCode(); expect != actual {
		t.Errorf("ErrorCode: %q != %q", expect, actual)
	}
}

// TC2: Validate SDK can handle missing Code field
// (missing X-Amzn-Query-Error)
func TestQueryCompatible_SEP2(t *testing.T) {
	svc := New(Options{
		HTTPClient: mockHTTP{
			BodyType: "aws.protocoltests.json#ItemNotFound",
		},
	})
	_, err := svc.GetItem(context.Background(), &GetItemInput{})
	if err == nil {
		t.Fatal("expect err, got none")
	}

	var terr smithy.APIError
	if !errors.As(err, &terr) {
		t.Fatalf("expect smithy.APIError, got %T", err)
	}

	expect := "ItemNotFound"
	if actual := terr.ErrorCode(); expect != actual {
		t.Errorf("ErrorCode: %q != %q", expect, actual)
	}
}

// TC3: Validate SDK can parse Type field
func TestQueryCompatible_SEP3(t *testing.T) {
	svc := New(Options{
		HTTPClient: mockHTTP{
			QueryError: "AWS.SimpleQueueService.NonExistentQueue;Sender",
			BodyType:   "aws.protocoltests.json#ItemNotFound",
		},
	})
	_, err := svc.GetItem(context.Background(), &GetItemInput{})
	if err == nil {
		t.Fatal("expect err, got none")
	}

	var terr smithy.APIError
	if !errors.As(err, &terr) {
		t.Fatalf("expect smithy.APIError, got %T", err)
	}

	expectCode := "AWS.SimpleQueueService.NonExistentQueue"
	if actual := terr.ErrorCode(); expectCode != actual {
		t.Errorf("ErrorCode: %q != %q", expectCode, actual)
	}

	expectFault := smithy.FaultClient
	if actual := terr.ErrorFault(); expectFault != actual {
		t.Errorf("ErrorFault: %v != %v", expectFault, actual)
	}
}

// TC4: Validate SDK sends x-amzn-query-mode header when service has
// @awsQueryCompatible trait
func TestQueryCompatible_SEP4(t *testing.T) {
	var req http.Request
	svc := New(Options{
		HTTPClient: mockHTTP{},
		APIOptions: []func(*middleware.Stack) error{
			func(s *middleware.Stack) error {
				return smithyprivateprotocol.AddCaptureRequestMiddleware(s, &req)
			},
		},
	})
	svc.GetItem(context.Background(), &GetItemInput{})

	expect := "true"
	if actual := req.Header.Get("X-Amzn-Query-Mode"); expect != actual {
		t.Errorf("X-Amzn-Query-Mode header: %q != %q", expect, actual)
	}
}

// TC5 covered in internal/protocoltest/jsonrpc10/querycompatible_test.go
