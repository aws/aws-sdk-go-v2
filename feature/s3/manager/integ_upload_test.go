// +build integration

package manager_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"regexp"
	"testing"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/smithy-go/middleware"
)

var integBuf12MB = make([]byte, 1024*1024*12)
var integMD512MB = fmt.Sprintf("%x", md5.Sum(integBuf12MB))

func TestInteg_UploadConcurrently(t *testing.T) {
	key := "12mb-1"
	mgr := manager.NewUploader(client)
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
	if input, ok := in.Parameters.(*s3.UploadPartInput); ok && input.PartNumber == 2 {
		ctx = v4.SetPayloadHash(ctx, "000")
	}

	return next.HandleSerialize(ctx, in)
}

func TestInteg_UploadFailCleanup(t *testing.T) {
	key := "12mb-leave"
	mgr := manager.NewUploader(client, func(u *manager.Uploader) {
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
	var uf manager.MultiUploadFailure
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
