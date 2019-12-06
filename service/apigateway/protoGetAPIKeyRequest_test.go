package apigateway_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/google/go-cmp/cmp"
)

func TestProtoGetApiKeyRequest_Diff(t *testing.T) {
	svc := apigateway.New(mock.Config())

	input := types.GetApiKeyInput{
		ApiKey:       aws.String("mock key"),
		IncludeValue: aws.Bool(true),
	}

	request := svc.GetApiKeyRequest(&input)
	_, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	prototypeRequest := svc.ProtoGetAPIKeyRequest(&input)
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
}
