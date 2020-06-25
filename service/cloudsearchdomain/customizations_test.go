package cloudsearchdomain_test

import (
	"errors"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/defaults"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/service/cloudsearchdomain"
)

func TestRequireEndpointIfRegionProvided(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "mock-region"
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("")

	svc := cloudsearchdomain.New(cfg)
	req := svc.SearchRequest(nil)
	err := req.Build()

	if e, a := "", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err == nil {
		t.Errorf("expect error, got none")
	}
	var expected *aws.MissingEndpointError
	if !errors.As(err, &expected) {
		t.Fatalf("expected %T, got %T", expected, err)
	}
}

func TestRequireEndpointIfNoRegionProvided(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("")

	svc := cloudsearchdomain.New(cfg)
	req := svc.SearchRequest(nil)
	err := req.Build()

	if e, a := "", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err == nil {
		t.Errorf("expect error, got none")
	}
	var expected *aws.MissingEndpointError
	if !errors.As(err, &expected) {
		t.Fatalf("expected %T, got %T", expected, err)
	}
}

func TestRequireEndpointUsed(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "mock-region"
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("https://endpoint")

	svc := cloudsearchdomain.New(cfg)
	req := svc.SearchRequest(nil)
	err := req.Build()

	if e, a := "https://endpoint", req.Endpoint.URL; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
