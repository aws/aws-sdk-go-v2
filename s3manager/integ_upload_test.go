// +build integration

package s3manager_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

var integBuf12MB = make([]byte, 1024*1024*12)
var integMD512MB = fmt.Sprintf("%x", md5.Sum(integBuf12MB))

func TestUploadConcurrently(t *testing.T) {
	key := "12mb-1"
	mgr := s3manager.NewUploader(client)
	out, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
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

type invalidateHash struct{}

func (b *invalidateHash) ID() string {
	return "s3manager:InvalidateHash"
}

func (b *invalidateHash) RegisterMiddleware(stack *middleware.Stack) error {
	return stack.Serialize.Add(b, middleware.After)
}

func (b *invalidateHash) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if input, ok := in.Parameters.(*s3.UploadPartInput); ok && aws.ToInt32(input.PartNumber) == 1 {
		in.Request.(*smithyhttp.Request).Header.Set("X-Amz-Content-Sha256", "000")
	}

	return next.HandleSerialize(ctx, in)
}

func TestUploadFailCleanup(t *testing.T) {
	key := "12mb-leave"
	mgr := s3manager.NewUploader(client, func(u *s3manager.Uploader) {
		u.LeavePartsOnError = false
		u.ClientOptions = append(u.ClientOptions, func(options *s3.Options) {
			options.APIOptions = append(options.APIOptions, (&invalidateHash{}).RegisterMiddleware)
		})
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	})
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}

	uploadID := ""
	var uf s3manager.MultiUploadFailure
	if !errors.As(err, &uf) {
		t.Errorf("")
	} else if uploadID = uf.UploadID(); len(uploadID) == 0 {
		t.Errorf("expect upload ID to not be empty, but was")
	}

	_, err = client.ListParts(context.Background(), &s3.ListPartsInput{
		Bucket: bucketName, Key: &key, UploadId: &uploadID,
	})
	if err == nil {
		t.Errorf("expect error for list parts, but got none")
	}
}
