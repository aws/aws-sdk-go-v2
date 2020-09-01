package customizations_test

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
)

func doBenchScan(b *testing.B, client *dynamodb.Client) {
	b.Helper()

	tableName := "mockTable"
	params := dynamodb.ScanInput{
		TableName: &tableName,
	}
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.Scan(ctx, &params)
		if err != nil {
			b.Fatalf("expect no error, %v", err)
		}
	}
}

func BenchmarkCustomizations(b *testing.B) {
	testdataFilename := filepath.Join("testdata", "scan_response.short.json")
	gzipBody, body, err := loadBenchTestData(testdataFilename)
	if err != nil {
		b.Fatalf("failed to load test data, %s, %v", testdataFilename, err)
	}

	b.Run("all", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				//ChecksumHeaderValue: []string{"2981943876"},
				ChecksumHeaderValue: []string{"891511383"},
				ScanRespGzipBody:    gzipBody,
			},
		})

		doBenchScan(b, client)
	})

	b.Run("none", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ScanRespBody: body,
			},
			DisableAcceptEncodingGzip:       true,
			DisableValidateResponseChecksum: true,
		})

		doBenchScan(b, client)
	})

	b.Run("validate checksum only", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				//ChecksumHeaderValue: []string{"3977635235"},
				ChecksumHeaderValue: []string{"512691431"},
				ScanRespBody:        body,
			},
			DisableAcceptEncodingGzip: true,
		})

		doBenchScan(b, client)
	})

	b.Run("accept encoding gzip only", func(b *testing.B) {
		client := dynamodb.New(dynamodb.Options{
			Credentials: aws.NewStaticCredentialsProvider("AKID", "SECRET", ""),
			HTTPClient: &mockClient{
				ScanRespGzipBody: gzipBody,
			},
			DisableValidateResponseChecksum: true,
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
