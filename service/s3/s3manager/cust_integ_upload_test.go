// +build integration

package s3manager_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var integBuf12MB = make([]byte, 1024*1024*12)
var integMD512MB = fmt.Sprintf("%x", md5.Sum(integBuf12MB))

func TestInteg_UploadConcurrently(t *testing.T) {
	key := "12mb-1"
	mgr := s3manager.NewUploader(integCfg)
	out, err := mgr.Upload(&s3manager.UploadInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if len(out.UploadID) == 0 {
		t.Errorf("expect upload ID but was empty")
	}

	re := regexp.MustCompile(`^https?://.+/` + key + `$`)
	if e, a := re.String(), out.Location; !re.MatchString(a) {
		t.Errorf("expect %s to match URL regexp %q, did not", e, a)
	}

	validate(t, key, integMD512MB)
}

func TestInteg_UploadFailCleanup(t *testing.T) {
	// Break checksum on 2nd part so it fails
	part := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		if r.Operation.Name == "UploadPart" {
			if part == 1 {
				r.HTTPRequest.Header.Set("X-Amz-Content-Sha256", "000")
			}
			part++
		}
	})

	key := "12mb-leave"
	mgr := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		u.LeavePartsOnError = false
	})
	_, err := mgr.Upload(&s3manager.UploadInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	})
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}

	aerr := err.(awserr.Error)
	if e, a := "MissingRegion", aerr.Code(); strings.Contains(a, e) {
		t.Errorf("expect %q to not be in error code %q", e, a)
	}

	uploadID := ""
	merr := err.(s3manager.MultiUploadFailure)
	if uploadID = merr.UploadID(); len(uploadID) == 0 {
		t.Errorf("expect upload ID to not be empty, but was")
	}

	_, err = svc.ListPartsRequest(&types.ListPartsInput{
		Bucket: bucketName, Key: &key, UploadId: &uploadID,
	}).Send(context.Background())
	if err == nil {
		t.Errorf("expect error for list parts, but got none")
	}
}
