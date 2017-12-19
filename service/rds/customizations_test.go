package rds

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestPresignWithPresignNotSet(t *testing.T) {
	reqs := map[string]*request.Request{}

	cfg := unit.Config()
	cfg.Region = "us-west-2"
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc := New(cfg)

	f := func() {
		// Doesn't panic on nil input
		req := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	}
	if paniced, p := awstesting.DidPanic(f); paniced {
		t.Errorf("expect no panic, got %v", p)
	}

	reqs[opCopyDBSnapshot] = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
	}).Request

	reqs[opCreateDBInstanceReadReplica] = svc.CreateDBInstanceReadReplicaRequest(&CreateDBInstanceReadReplicaInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBInstanceIdentifier: aws.String("foo"),
		DBInstanceIdentifier:       aws.String("bar"),
	}).Request

	for op, req := range reqs {
		req.Sign()
		b, _ := ioutil.ReadAll(req.HTTPRequest.Body)
		q, _ := url.ParseQuery(string(b))

		u, _ := url.QueryUnescape(q.Get("PreSignedUrl"))

		exp := fmt.Sprintf(`^https://rds.us-west-1\.amazonaws\.com/\?Action=%s.+?DestinationRegion=us-west-2.+`, op)
		if re, a := regexp.MustCompile(exp), u; !re.MatchString(a) {
			t.Errorf("expect %s to match %s", re, a)
		}
	}
}

func TestPresignWithPresignSet(t *testing.T) {
	reqs := map[string]*request.Request{}

	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := New(cfg)

	f := func() {
		// Doesn't panic on nil input
		req := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	}
	if paniced, p := awstesting.DidPanic(f); paniced {
		t.Errorf("expect no panic, got %v", p)
	}

	reqs[opCopyDBSnapshot] = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
		PreSignedUrl:               aws.String("presignedURL"),
	}).Request

	reqs[opCreateDBInstanceReadReplica] = svc.CreateDBInstanceReadReplicaRequest(&CreateDBInstanceReadReplicaInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBInstanceIdentifier: aws.String("foo"),
		DBInstanceIdentifier:       aws.String("bar"),
		PreSignedUrl:               aws.String("presignedURL"),
	}).Request

	for _, req := range reqs {
		req.Sign()

		b, _ := ioutil.ReadAll(req.HTTPRequest.Body)
		q, _ := url.ParseQuery(string(b))

		u, _ := url.QueryUnescape(q.Get("PreSignedUrl"))
		if e, a := "presignedURL", u; !strings.Contains(a, e) {
			t.Errorf("expect %s to be in %s", e, a)
		}
	}
}

func TestPresignWithSourceNotSet(t *testing.T) {
	reqs := map[string]*request.Request{}

	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := New(cfg)

	f := func() {
		// Doesn't panic on nil input
		req := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	}
	if paniced, p := awstesting.DidPanic(f); paniced {
		t.Errorf("expect no panic, got %v", p)
	}

	reqs[opCopyDBSnapshot] = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
	}).Request

	for _, req := range reqs {
		_, err := req.Presign(5 * time.Minute)
		if err != nil {
			t.Fatal(err)
		}
	}
}
