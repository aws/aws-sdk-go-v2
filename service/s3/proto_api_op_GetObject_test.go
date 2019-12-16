package s3_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting/mock"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestProtoGetObjectRequest_Send(t *testing.T) {
	svc := s3.New(mock.Config())
	input := types.GetObjectInput{
		Bucket:                     aws.String("mock bucket"),
		IfMatch:                    nil,
		IfModifiedSince:            nil,
		IfNoneMatch:                nil,
		IfUnmodifiedSince:          nil,
		Key:                        aws.String("mock key"),
		PartNumber:                 nil,
		Range:                      nil,
		RequestPayer:               "",
		ResponseCacheControl:       nil,
		ResponseContentDisposition: nil,
		ResponseContentEncoding:    nil,
		ResponseContentLanguage:    nil,
		ResponseContentType:        nil,
		ResponseExpires:            nil,
		SSECustomerAlgorithm:       nil,
		SSECustomerKey:             nil,
		SSECustomerKeyMD5:          nil,
		VersionId:                  nil,
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

	fmt.Println(expectedResponse.GetObjectOutput, prototypeResponse.GetObjectOutput)

}
