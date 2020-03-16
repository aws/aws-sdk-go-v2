package apigateway_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/google/go-cmp/cmp"
)

func TestProtoCreateApiKeyRequestMarshaler_Diff(t *testing.T) {
	svc := apigateway.New(mock.Config())
	input := apigateway.CreateApiKeyInput{
		CustomerId:         aws.String("mock id"),
		Description:        aws.String("mock operation description"),
		Enabled:            aws.Bool(true),
		GenerateDistinctId: aws.Bool(true),
		Name:               aws.String("mock name"),
		StageKeys: []apigateway.StageKey{apigateway.StageKey{
			RestApiId: aws.String("mock rest api id"),
			StageName: aws.String("mock stage name"),
		}},
		Tags:  map[string]string{"a": "1", "b": "2"},
		Value: aws.String("mock value"),
	}

	request := svc.CreateApiKeyRequest(&input)
	_, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	prototypeRequest := svc.ProtoCreateAPIKeyRequest(&input)
	_, err = prototypeRequest.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(request.HTTPRequest.Header, prototypeRequest.HTTPRequest.Header); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
	if diff := cmp.Diff(request.HTTPRequest.URL, prototypeRequest.HTTPRequest.URL); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
	if diff := cmp.Diff(request.Body, prototypeRequest.Body, cmp.AllowUnexported(bytes.Reader{})); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
}
