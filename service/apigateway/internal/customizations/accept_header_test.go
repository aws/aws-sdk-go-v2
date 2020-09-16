package customizations_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"

	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

type mockClient struct {
}

func (m *mockClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
	"_links": {
		"curies": {
			"href": "https://docs.aws.amazon.com/apigateway/latest/developerguide/account-apigateway-{rel}.html",
			"name": "account",
			"templated": true
		},
		"self": {
			"href": "/account"
		},
		"account:update": {
			"href": "/account"
		}
	},
	"cloudwatchRoleArn": "arn:aws:iam::123456789012:role/apigAwsProxyRole"
}`))),
	}, nil
}

func TestAddAcceptHeader(t *testing.T) {
	options := apigateway.Options{
		Credentials: aws.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
				Source: "unit test credentials",
			}},
		Retryer:    aws.NoOpRetryer{},
		HTTPClient: &mockClient{},
		Region:     "mock-region",
	}
	svc := apigateway.New(options)
	fm := requestRetrieverMiddleware{}

	_, err := svc.GetAccount(context.Background(), &apigateway.GetAccountInput{
		Name:  ptr.String("mock_name"),
		Title: ptr.String("mock_title"),
	}, func(options *apigateway.Options) {
		options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
			stack.Build.Add(&fm, middleware.After)
			return nil
		})
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	req := fm.request

	if e, a := "application/json", req.Header.Get("Accept"); e != a {
		t.Fatalf("Expected Accept header to be set to %v, got %v", e, a)
	}
}

type requestRetrieverMiddleware struct {
	request *smithyhttp.Request
}

func (*requestRetrieverMiddleware) ID() string { return "S3:requestRetrieverMiddleware" }

func (rm *requestRetrieverMiddleware) HandleBuild(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}
	rm.request = req
	return next.HandleBuild(ctx, in)
}
