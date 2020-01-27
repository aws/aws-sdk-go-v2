package s3_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
)

func TestProtoGetObjectRequest_Send(t *testing.T) {
	svc := s3.New(mock.Config())
	input := s3.GetObjectInput{
		Bucket:                     aws.String("mock bucket"),
		IfMatch:                    aws.String("mock value"),
		IfModifiedSince:            aws.Time(time.Now()),
		IfNoneMatch:                aws.String("mock value for no match"),
		IfUnmodifiedSince:          aws.Time(time.Now().Add(-10 * time.Minute)),
		Key:                        aws.String("mock key"),
		PartNumber:                 aws.Int64(10),
		Range:                      aws.String("mock range"),
		RequestPayer:               s3.RequestPayerRequester,
		ResponseCacheControl:       aws.String("mock value"),
		ResponseContentDisposition: aws.String("mock value"),
		ResponseContentEncoding:    aws.String("mock value"),
		ResponseContentLanguage:    aws.String("mock value"),
		ResponseContentType:        aws.String("mock value"),
		ResponseExpires:            aws.Time(time.Now().Add(10 * time.Minute)),
		SSECustomerAlgorithm:       aws.String("AES256"),
		VersionId:                  aws.String("mock version value"),
	}

	request := svc.GetObjectRequest(&input)
	expectedResponse, err := request.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	prototypeRequest := svc.ProtoGetObjectRequest(&input)
	prototypeResponse, err := prototypeRequest.Send(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(expectedResponse.GetObjectOutput, prototypeResponse.GetObjectOutput); diff != "" {
		t.Errorf("Found diff: %v", diff)
	}
}
