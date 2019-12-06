package s3_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/enums"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/go-cmp/cmp"
)

func TestProtoPutObjectRequest_Diff(t *testing.T) {
	svc := s3.New(mock.Config())

	input := types.PutObjectInput{
		ACL:                enums.ObjectCannedACLAuthenticatedRead,
		Body:               nil,
		ContentLength:      aws.Int64(0),
		Bucket:             aws.String("mock bucket"),
		CacheControl:       aws.String("mock cache control"),
		ContentDisposition: aws.String("mock content disposition"),
		ContentLanguage:    aws.String("english"),
		ContentMD5:         aws.String("mock MD5"),
		ContentType:        aws.String("mock content type"),
		Expires:            aws.Time(time.Now()),
		GrantFullControl:   aws.String("mock full control"),
		GrantRead:          aws.String("mock read grant"),
		GrantReadACP:       aws.String("mock acp read"),
		GrantWriteACP:      aws.String("mock write acp"),
		Key:                aws.String("mock key"),
		Metadata: map[string]string{
			"mockMetaKey01": "mockMetaValue01",
			"mockMetaKey02": "mockMetaValue02",
			"mockMetaKey03": "mockMetaValue03",
		},
		ObjectLockLegalHoldStatus: enums.ObjectLockLegalHoldStatusOn,
		ObjectLockMode:            enums.ObjectLockModeCompliance,
		ObjectLockRetainUntilDate: aws.Time(time.Now()),
		RequestPayer:              enums.RequestPayerRequester,
		SSECustomerAlgorithm:      aws.String("mock sse cust Algo"),
		SSECustomerKey:            nil,
		SSECustomerKeyMD5:         aws.String("mock sse MD5"),
		SSEKMSEncryptionContext:   aws.String("mock encryption content"),
		SSEKMSKeyId:               aws.String("mock ssekmskey id"),
		ServerSideEncryption:      enums.ServerSideEncryptionAes256,
		StorageClass:              enums.StorageClassGlacier,
		Tagging:                   aws.String("mock tagging"),
		WebsiteRedirectLocation:   aws.String("mock redirection"),
	}

	// request created for existing put object request
	request := svc.PutObjectRequest(&input)
	_, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	// request created for prototyped put object request
	prototypeRequest := svc.ProtoPutObjectRequest(&input)
	_, err = prototypeRequest.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(request.HTTPRequest.Header, prototypeRequest.HTTPRequest.Header); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}

	if diff := cmp.Diff(request.HTTPRequest.URL, prototypeRequest.HTTPRequest.URL); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
}
