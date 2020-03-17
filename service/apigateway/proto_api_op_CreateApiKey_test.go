package apigateway_test

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
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

	reqBody, err := request.HTTPRequest.GetBody()
	if err != nil {
		t.Fatal("failed to read body from request", err)
	}

	protoBody, err := prototypeRequest.HTTPRequest.GetBody()
	if err != nil {
		t.Fatal("failed to read body from prototyped request", err)
	}
	assertJSON(t, reqBody, protoBody)
}

func assertJSON(t testing.TB, a, b io.ReadCloser) {
	t.Helper()
	buf, err := ioutil.ReadAll(a)
	if err != nil {
		t.Fatal("Error reading a", err)
	}
	protoBuff, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal("Error reading b", err)
	}
	var av, bv interface{}

	if err := json.Unmarshal(buf, &av); err != nil {
		t.Fatalf("assertJSON: unable to unmarshal a, %v", err)
	}

	if err := json.Unmarshal(protoBuff, &bv); err != nil {
		t.Fatalf("assertJSON: unable to unmarshal b, %v", err)
	}

	if !reflect.DeepEqual(av, bv) {
		t.Fatalf("JSON are not equal\nexpect:\n%v\nactual:\n%v\n", a, b)
	}
}
