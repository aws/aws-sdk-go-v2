package s3_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/enums"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var s3LocationTests = []struct {
	body string
	loc  string
}{
	{`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/"/>`, ``},
	{`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">EU</LocationConstraint>`, `EU`},
}

func TestGetBucketLocation(t *testing.T) {
	for _, test := range s3LocationTests {
		s := s3.New(unit.Config())
		s.Handlers.Send.Clear()
		s.Handlers.Send.PushBack(func(r *request.Request) {
			reader := ioutil.NopCloser(bytes.NewReader([]byte(test.body)))
			r.HTTPResponse = &http.Response{StatusCode: 200, Body: reader}
		})

		req := s.GetBucketLocationRequest(&types.GetBucketLocationInput{Bucket: aws.String("bucket")})
		resp, err := req.Send(context.Background())
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		if test.loc == "" {
			if v := resp.LocationConstraint; len(v) > 0 {
				t.Errorf("expect location constraint to be empty, got %v", v)
			}
		} else {
			if e, a := test.loc, string(resp.LocationConstraint); e != a {
				t.Errorf("expect %s location constraint, got %v", e, a)
			}
		}
	}
}

func TestNormalizeBucketLocation(t *testing.T) {
	cases := []struct {
		In, Out string
	}{
		{"", "us-east-1"},
		{"EU", "eu-west-1"},
		{"us-east-1", "us-east-1"},
		{"something", "something"},
	}

	for i, c := range cases {
		actual := s3.NormalizeBucketLocation(enums.BucketLocationConstraint(c.In))
		if e, a := c.Out, string(actual); e != a {
			t.Errorf("%d, expect %s bucket location, got %s", i, e, a)
		}
	}
}

func TestWithNormalizeBucketLocation(t *testing.T) {
	req := &request.Request{}
	req.ApplyOptions(s3.WithNormalizeBucketLocation)

	cases := []struct {
		In, Out string
	}{
		{"", "us-east-1"},
		{"EU", "eu-west-1"},
		{"us-east-1", "us-east-1"},
		{"something", "something"},
	}

	for i, c := range cases {
		req.Data = &types.GetBucketLocationOutput{
			LocationConstraint: enums.BucketLocationConstraint(c.In),
		}
		req.Handlers.Unmarshal.Run(req)

		v := req.Data.(*types.GetBucketLocationOutput).LocationConstraint
		if e, a := c.Out, string(v); e != a {
			t.Errorf("%d, expect %s bucket location, got %s", i, e, a)
		}
	}
}

func TestPopulateLocationConstraint(t *testing.T) {
	s := s3.New(unit.Config())
	in := &types.CreateBucketInput{
		Bucket: aws.String("bucket"),
	}
	req := s.CreateBucketRequest(in)
	if err := req.Build(); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	v, _ := awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")
	if e, a := "mock-region", string(v[0].(enums.BucketLocationConstraint)); e != a {
		t.Errorf("expect %s location constraint, got %s", e, a)
	}
	if v := in.CreateBucketConfiguration; v != nil {
		// don't modify original params
		t.Errorf("expect create bucket Configuration to be nil, got %s", *v)
	}
}

func TestNoPopulateLocationConstraintIfProvided(t *testing.T) {
	s := s3.New(unit.Config())
	req := s.CreateBucketRequest(&types.CreateBucketInput{
		Bucket:                    aws.String("bucket"),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{},
	})
	if err := req.Build(); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	v, _ := awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")
	if l := len(v); l != 1 {
		t.Errorf("expect empty string only, got %d elements", l)
	}
	if v[0].(enums.BucketLocationConstraint) != enums.BucketLocationConstraint("") {
		t.Errorf("expected empty string, but received %v", v[0])
	}
}

func TestNoPopulateLocationConstraintIfClassic(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-east-1"

	s := s3.New(cfg)
	req := s.CreateBucketRequest(&types.CreateBucketInput{
		Bucket: aws.String("bucket"),
	})
	if err := req.Build(); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	v, _ := awsutil.ValuesAtPath(req.Params, "CreateBucketConfiguration.LocationConstraint")
	if l := len(v); l != 0 {
		t.Errorf("expect no values, got %d", l)
	}
}
