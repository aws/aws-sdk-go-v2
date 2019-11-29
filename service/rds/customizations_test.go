package rds

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func TestCopyDBSnapshotNoPanic(t *testing.T) {

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

}

func TestPresignCrossRegionRequest(t *testing.T) {

	cfg := unit.Config()
	cfg.Region = "us-west-2"
	cfg.EndpointResolver = endpoints.NewDefaultResolver()

	svc := New(cfg)
	const regexPattern = `^https://rds.us-west-1\.amazonaws\.com/\?Action=%s.+?DestinationRegion=%s.+`

	cases := map[string]struct {
		Req    *request.Request
		Assert func(*testing.T, string)
	}{
		opCopyDBSnapshot: {
			Req: func() *request.Request {
				req := svc.CopyDBSnapshotRequest(&types.CopyDBSnapshotInput{
					SourceRegion:               aws.String("us-west-1"),
					SourceDBSnapshotIdentifier: aws.String("foo"),
					TargetDBSnapshotIdentifier: aws.String("bar"),
				})
				return req.Request
			}(),
			Assert: assertAsRegexMatch(fmt.Sprintf(regexPattern,
				opCopyDBSnapshot, cfg.Region)),
		},

		opCreateDBInstanceReadReplica: {
			Req: func() *request.Request {
				req := svc.CreateDBInstanceReadReplicaRequest(
					&types.CreateDBInstanceReadReplicaInput{
						SourceRegion:               aws.String("us-west-1"),
						SourceDBInstanceIdentifier: aws.String("foo"),
						DBInstanceIdentifier:       aws.String("bar"),
					})
				return req.Request
			}(),
			Assert: assertAsRegexMatch(fmt.Sprintf(regexPattern,
				opCreateDBInstanceReadReplica, cfg.Region)),
		},
		opCopyDBClusterSnapshot: {
			Req: func() *request.Request {
				req := svc.CopyDBClusterSnapshotRequest(
					&types.CopyDBClusterSnapshotInput{
						SourceRegion:                      aws.String("us-west-1"),
						SourceDBClusterSnapshotIdentifier: aws.String("foo"),
						TargetDBClusterSnapshotIdentifier: aws.String("bar"),
					})
				return req.Request
			}(),
			Assert: assertAsRegexMatch(fmt.Sprintf(regexPattern,
				opCopyDBClusterSnapshot, cfg.Region)),
		},
		opCreateDBCluster: {
			Req: func() *request.Request {
				req := svc.CreateDBClusterRequest(
					&types.CreateDBClusterInput{
						SourceRegion:        aws.String("us-west-1"),
						DBClusterIdentifier: aws.String("foo"),
						Engine:              aws.String("bar"),
					})
				return req.Request
			}(),
			Assert: assertAsRegexMatch(fmt.Sprintf(regexPattern,
				opCreateDBCluster, cfg.Region)),
		},
		opCopyDBSnapshot + " same region": {
			Req: func() *request.Request {
				req := svc.CopyDBSnapshotRequest(&types.CopyDBSnapshotInput{
					SourceRegion:               aws.String("us-west-2"),
					SourceDBSnapshotIdentifier: aws.String("foo"),
					TargetDBSnapshotIdentifier: aws.String("bar"),
				})
				return req.Request
			}(),
			Assert: assertAsEmpty(),
		},
		opCreateDBInstanceReadReplica + " same region": {
			Req: func() *request.Request {
				req := svc.CreateDBInstanceReadReplicaRequest(&types.CreateDBInstanceReadReplicaInput{
					SourceRegion:               aws.String("us-west-2"),
					SourceDBInstanceIdentifier: aws.String("foo"),
					DBInstanceIdentifier:       aws.String("bar"),
				})
				return req.Request
			}(),
			Assert: assertAsEmpty(),
		},
		opCopyDBClusterSnapshot + " same region": {
			Req: func() *request.Request {
				req := svc.CopyDBClusterSnapshotRequest(
					&types.CopyDBClusterSnapshotInput{
						SourceRegion:                      aws.String("us-west-2"),
						SourceDBClusterSnapshotIdentifier: aws.String("foo"),
						TargetDBClusterSnapshotIdentifier: aws.String("bar"),
					})
				return req.Request
			}(),
			Assert: assertAsEmpty(),
		},
		opCreateDBCluster + " same region": {
			Req: func() *request.Request {
				req := svc.CreateDBClusterRequest(
					&types.CreateDBClusterInput{
						SourceRegion:        aws.String("us-west-2"),
						DBClusterIdentifier: aws.String("foo"),
						Engine:              aws.String("bar"),
					})
				return req.Request
			}(),
			Assert: assertAsEmpty(),
		},
		opCopyDBSnapshot + " presignURL set": {
			Req: func() *request.Request {
				req := svc.CopyDBSnapshotRequest(&types.CopyDBSnapshotInput{
					SourceRegion:               aws.String("us-west-1"),
					SourceDBSnapshotIdentifier: aws.String("foo"),
					TargetDBSnapshotIdentifier: aws.String("bar"),
					PreSignedUrl:               aws.String("mockPresignedURL"),
				})
				return req.Request
			}(),
			Assert: assertAsEqual("mockPresignedURL"),
		},
		opCreateDBInstanceReadReplica + " presignURL set": {
			Req: func() *request.Request {
				req := svc.CreateDBInstanceReadReplicaRequest(&types.CreateDBInstanceReadReplicaInput{
					SourceRegion:               aws.String("us-west-1"),
					SourceDBInstanceIdentifier: aws.String("foo"),
					DBInstanceIdentifier:       aws.String("bar"),
					PreSignedUrl:               aws.String("mockPresignedURL"),
				})
				return req.Request
			}(),
			Assert: assertAsEqual("mockPresignedURL"),
		},
		opCopyDBClusterSnapshot + " presignURL set": {
			Req: func() *request.Request {
				req := svc.CopyDBClusterSnapshotRequest(
					&types.CopyDBClusterSnapshotInput{
						SourceRegion:                      aws.String("us-west-1"),
						SourceDBClusterSnapshotIdentifier: aws.String("foo"),
						TargetDBClusterSnapshotIdentifier: aws.String("bar"),
						PreSignedUrl:                      aws.String("mockPresignedURL"),
					})
				return req.Request
			}(),
			Assert: assertAsEqual("mockPresignedURL"),
		},
		opCreateDBCluster + " presignURL set": {
			Req: func() *request.Request {
				req := svc.CreateDBClusterRequest(
					&types.CreateDBClusterInput{
						SourceRegion:        aws.String("us-west-1"),
						DBClusterIdentifier: aws.String("foo"),
						Engine:              aws.String("bar"),
						PreSignedUrl:        aws.String("mockPresignedURL"),
					})
				return req.Request
			}(),
			Assert: assertAsEqual("mockPresignedURL"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if err := c.Req.Sign(); err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			b, _ := ioutil.ReadAll(c.Req.HTTPRequest.Body)
			q, _ := url.ParseQuery(string(b))

			u, _ := url.QueryUnescape(q.Get("PreSignedUrl"))

			c.Assert(t, u)

		})
	}
}

func TestPresignWithSourceNotSet(t *testing.T) {
	reqs := map[string]*request.Request{}

	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := New(cfg)

	reqs[opCopyDBSnapshot] = svc.CopyDBSnapshotRequest(&types.CopyDBSnapshotInput{
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

func assertAsRegexMatch(exp string) func(*testing.T, string) {
	return func(t *testing.T, v string) {
		t.Helper()

		if re, a := regexp.MustCompile(exp), v; !re.MatchString(a) {
			t.Errorf("expect %s to match %s", re, a)
		}
	}
}

func assertAsEmpty() func(*testing.T, string) {
	return func(t *testing.T, v string) {
		t.Helper()

		if len(v) != 0 {
			t.Errorf("expect empty, got %v", v)
		}
	}
}

func assertAsEqual(expect string) func(*testing.T, string) {
	return func(t *testing.T, v string) {
		t.Helper()

		if e, a := expect, v; e != a {
			t.Errorf("expect %v, got %v", e, a)
		}
	}
}
