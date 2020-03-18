package iotdataplane_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

func TestRequireEndpointIfRegionProvided(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("")

	svc := iotdataplane.New(cfg)
	req := svc.GetThingShadowRequest(nil)
	err := req.Build()

	if e, a := "", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err == nil {
		t.Errorf("expect error, got none")
	}
	var awsErr awserr.Error
	if !errors.As(err, &awsErr) {
		t.Fatalf("expected awserr.Error, got %T", err)
	}
	if e, a := aws.ErrCodeMissingEndpoint, awsErr.Code(); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
}

func TestRequireEndpointIfNoRegionProvided(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("")

	svc := iotdataplane.New(cfg)

	req := svc.GetThingShadowRequest(nil)
	err := req.Build()

	if e, a := "", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err == nil {
		t.Errorf("expect error, got none")
	}
	var awsErr awserr.Error
	if !errors.As(err, &awsErr) {
		t.Fatalf("expected awserr.Error, got %T", err)
	}
	if e, a := aws.ErrCodeMissingEndpoint, awsErr.Code(); e != a {
		t.Errorf("expect %v to be %v", e, a)
	}
}

func TestRequireEndpointUsed(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://endpoint")

	svc := iotdataplane.New(cfg)
	req := svc.GetThingShadowRequest(nil)
	err := req.Build()

	if e, a := "https://endpoint", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
