//go:build integration
// +build integration

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestExpressRoundTripObject(t *testing.T) {
	const key = "TestExpressRoundTripObject"
	const value = "TestExpressRoundTripObjectValue"

	_, err := s3client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
		Body:   strings.NewReader(value),
	}, withAssertExpress)
	if err != nil {
		t.Fatalf("put object: %v", err)
	}

	resp, err := s3client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
	}, withAssertExpress)
	if err != nil {
		t.Fatalf("get object: %v", err)
	}

	obj, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read object response body: %v", err)
	}

	if string(obj) != value {
		t.Fatalf("round-trip object didn't match: %q", obj)
	}
}

func TestExpressPresignGetObject(t *testing.T) {
	const key = "TestExpressPresignGetObject"
	const value = "TestExpressPresignGetObjectValue"

	_, err := s3client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
		Body:   strings.NewReader(value),
	}, withAssertExpress)
	if err != nil {
		t.Fatalf("put object: %v", err)
	}

	presigner := s3.NewPresignClient(s3client)
	req, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatalf("presign get object: %v", err)
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		log.Fatalf("parse url: %v", err)
	}

	resp, err := http.DefaultClient.Do(&http.Request{
		Method: req.Method,
		URL:    u,
	})
	if err != nil {
		log.Fatalf("call presigned get object: %v", err)
	}

	obj, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response obj: %v", err)
	}

	if string(obj) != value { // ignore the status code, response body wouldn't match anyway
		t.Fatalf("presigned get didn't match: %q", obj)
	}
}

func TestExpressPresignPutObject(t *testing.T) {
	const key = "TestExpressPresignPutObject"
	const value = "TestExpressPresignPutObjectValue"

	presigner := s3.NewPresignClient(s3client)
	req, err := presigner.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatalf("presign put object: %v", err)
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		log.Fatal(err)
	}

	presp, err := http.DefaultClient.Do(&http.Request{
		Method:        req.Method,
		URL:           u,
		Body:          io.NopCloser(strings.NewReader(value)),
		ContentLength: int64(len(value)),
	})
	if err != nil {
		log.Fatal(err)
	}
	if presp.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(presp.Body)
		if err != nil {
			log.Fatalf("read presigned put object response body: %s", msg)
		}

		log.Fatalf("call presigned put object: %s", msg)
	}

	gresp, err := s3client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
	}, withAssertExpress)
	if err != nil {
		log.Fatalf("get object: %v", err)
	}

	obj, err := io.ReadAll(gresp.Body)
	if err != nil {
		log.Fatalf("read response body: %v", err)
	}

	if string(obj) != value {
		t.Fatalf("presigned put didn't match: %q", obj)
	}
}

func TestExpressUploaderDefaultChecksum(t *testing.T) {
	const key = "TestExpressUploaderDefaultChecksum"
	const valueLen = 12 * 1024 * 1024 // default/min part size is 5MiB, guarantee 2 full + 1 partial

	value := make([]byte, valueLen)

	uploader := manager.NewUploader(s3client, func(u *manager.Uploader) {
		u.ClientOptions = append(u.ClientOptions, withAssertExpress)
	})
	out, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &setupMetadata.ExpressBucket,
		Key:    aws.String(key),
		Body:   bytes.NewBuffer(value),
	})
	if err != nil {
		log.Fatal(err)
	}

	if out.ChecksumCRC32 == nil {
		log.Fatal("upload didn't default to crc32")
	}
}

func TestExpressUploaderManualChecksum(t *testing.T) {
	const key = "TestExpressUploaderManualChecksum"
	const valueLen = 12 * 1024 * 1024

	value := make([]byte, valueLen)

	uploader := manager.NewUploader(s3client, func(u *manager.Uploader) {
		u.ClientOptions = append(u.ClientOptions, withAssertExpress)
	})
	out, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:            &setupMetadata.ExpressBucket,
		Key:               aws.String(key),
		Body:              bytes.NewBuffer(value),
		ChecksumAlgorithm: types.ChecksumAlgorithmCrc32c,
	})
	if err != nil {
		log.Fatal(err)
	}

	if out.ChecksumCRC32C == nil {
		log.Fatal("upload didn't use explicit crc32c")
	}
}

var withAssertExpress = s3.WithAPIOptions(func(s *middleware.Stack) error {
	return s.Finalize.Add(&assertExpress{}, middleware.After)
})

type assertExpress struct{}

func (*assertExpress) ID() string {
	return "assertExpress"
}

func (m *assertExpress) HandleFinalize(ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unexpected transport type %T", in.Request)
	}

	sig := awstesting.ParseSigV4Signature(req.Header)
	if sig.SigningName != "s3express" {
		return out, metadata, fmt.Errorf("signing name is not s3express: %q", sig.SigningName)
	}

	return next.HandleFinalize(ctx, in)
}
