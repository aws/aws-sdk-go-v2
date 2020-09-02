package benchmark

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	awsOld "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	ddbOld "github.com/aws/aws-sdk-go/service/dynamodb"
)

func BenchmarkCustomizations_Old(b *testing.B) {
	testdataFilename := filepath.Join("testdata", "scan_response.short.json")
	_, body, err := loadBenchTestData(testdataFilename)
	if err != nil {
		b.Fatalf("failed to load test data, %s, %v", testdataFilename, err)
	}

	sess, err := session.NewSession(&awsOld.Config{
		Region: awsOld.String("us-west-2"),
	})
	sess.Handlers.Send.SwapNamed(request.NamedHandler{
		Name: corehandlers.SendHandler.Name,
		Fn: func(r *request.Request) {
			r.HTTPResponse = &http.Response{
				StatusCode: 200,
				Header: http.Header{
					"X-Amz-Crc32": []string{"512691431"},
				},
				ContentLength: int64(len(body)),
				Body:          ioutil.NopCloser(bytes.NewReader(body)),
			}
		}})
	if err != nil {
		b.Fatalf("failed to load session, %v", err)
	}

	b.Run("default", func(b *testing.B) {
		client := ddbOld.New(sess)

		doBenchScanOld(b, client)
	})

	b.Run("none enabled", func(b *testing.B) {
		client := ddbOld.New(sess, &awsOld.Config{
			DisableComputeChecksums: awsOld.Bool(true),
		})

		doBenchScanOld(b, client)
	})
}

func doBenchScanOld(b *testing.B, client *ddbOld.DynamoDB) {
	b.Helper()

	tableName := "mockTable"
	params := ddbOld.ScanInput{
		TableName: &tableName,
	}
	ctx := context.Background()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.ScanWithContext(ctx, &params)
			if err != nil {
				b.Fatalf("expect no error, %v", err)
			}
		}
	})
}

func BenchmarkCustomizations_Smithy(b *testing.B) {
	testdataFilename := filepath.Join("testdata", "scan_response.short.json")
	gzipBody, body, err := loadBenchTestData(testdataFilename)
	if err != nil {
		b.Fatalf("failed to load test data, %s, %v", testdataFilename, err)
	}

	b.Run("defaults", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ChecksumHeaderValue: []string{"512691431"},
				ScanRespBody:        body,
			},
		})

		doBenchScan(b, client)
	})

	b.Run("all enabled", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ChecksumHeaderValue: []string{"891511383"},
				ScanRespGzipBody:    gzipBody,
			},
			EnableAcceptEncodingGzip: true,
		})

		doBenchScan(b, client)
	})

	b.Run("none enabled", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ScanRespBody: body,
			},
			DisableValidateResponseChecksum: true,
		})

		doBenchScan(b, client)
	})

	b.Run("checksum only", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ChecksumHeaderValue: []string{"512691431"},
				ScanRespBody:        body,
			},
		})

		doBenchScan(b, client)
	})

	b.Run("gzip only", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ScanRespGzipBody: gzipBody,
			},
			DisableValidateResponseChecksum: true,
			EnableAcceptEncodingGzip:        true,
		})

		doBenchScan(b, client)
	})
}

type mockClient struct {
	ChecksumHeaderValue []string
	ScanRespGzipBody    []byte
	ScanRespBody        []byte
}

func (m *mockClient) Do(r *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/x-amz-json-1.0"},
		},
	}

	if m.ChecksumHeaderValue != nil {
		resp.Header["X-Amz-Crc32"] = m.ChecksumHeaderValue
	}

	if m.ScanRespGzipBody != nil {
		resp.Header["Content-Encoding"] = []string{"gzip"}
		resp.ContentLength = int64(len(m.ScanRespGzipBody))
		resp.Body = ioutil.NopCloser(bytes.NewReader(m.ScanRespGzipBody))
	} else if m.ScanRespBody != nil {
		resp.ContentLength = int64(len(m.ScanRespBody))
		resp.Body = ioutil.NopCloser(bytes.NewReader(m.ScanRespBody))
	} else {
		return nil, fmt.Errorf("no client mock response body set")
	}

	return resp, nil
}

func doBenchScan(b *testing.B, client *dynamodb.Client) {
	b.Helper()

	tableName := "mockTable"
	params := dynamodb.ScanInput{
		TableName: &tableName,
	}
	ctx := context.Background()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.Scan(ctx, &params)
			if err != nil {
				b.Fatalf("expect no error, %v", err)
			}
		}
	})
}

func loadBenchTestData(filename string) ([]byte, []byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open test data %v", err)
	}
	defer f.Close()

	var body bytes.Buffer
	if _, err := io.Copy(&body, f); err != nil {
		return nil, nil, fmt.Errorf("failed to read test data %v", err)
	}

	var gzipBody bytes.Buffer
	w := gzip.NewWriter(&gzipBody)
	w.Write(body.Bytes())
	w.Close()

	return gzipBody.Bytes(), body.Bytes(), nil
}
