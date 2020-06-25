package ec2_test

import (
	"io/ioutil"
	"net/url"
	"regexp"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/aws/endpoints"
	"github.com/jviney/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/jviney/aws-sdk-go-v2/service/ec2"
)

func TestCopySnapshotPresignedURL(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc := ec2.New(cfg)

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("expect CopySnapshotRequest with nil")
			}
		}()

		// Doesn't panic on nil input
		req := svc.CopySnapshotRequest(nil)
		req.Sign()
	}()

	req := svc.CopySnapshotRequest(&ec2.CopySnapshotInput{
		SourceRegion:     aws.String("us-west-1"),
		SourceSnapshotId: aws.String("snap-id"),
	})
	if err := req.Sign(); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, _ := ioutil.ReadAll(req.HTTPRequest.Body)
	q, _ := url.ParseQuery(string(b))
	u, _ := url.QueryUnescape(q.Get("PresignedUrl"))

	if e, a := "us-west-2", q.Get("DestinationRegion"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if e, a := "us-west-1", q.Get("SourceRegion"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	r := regexp.MustCompile(`^https://ec2\.us-west-1\.amazonaws\.com/.+&DestinationRegion=us-west-2`)
	if !r.MatchString(u) {
		t.Errorf("expect %v to match, got %v", r.String(), u)
	}
}
