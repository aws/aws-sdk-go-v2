//go:build integration
// +build integration

package manager_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_GetBucketRegion(t *testing.T) {
	expectRegion := integConfig.Region
	cases := map[string]struct {
		optFns []func(*s3.Options)
	}{
		"normal credentials": {[]func(o *s3.Options){}},
		"invalid credentials": {[]func(o *s3.Options){
			func(o *s3.Options) {
				o.Credentials = credentials.NewStaticCredentialsProvider("fakeDummy", "fakeDummy", "")
			},
		}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			region, err := manager.GetBucketRegion(context.Background(), s3.NewFromConfig(integConfig, c.optFns...), aws.ToString(bucketName))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if e, a := expectRegion, region; e != a {
				t.Errorf("expect %s bucket region, got %s", e, a)
			}
		})
	}
}
