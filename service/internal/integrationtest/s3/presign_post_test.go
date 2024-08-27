//go:build integration
// +build integration

package s3

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_PresigPost(t *testing.T) {

	const filePath = "sample.txt"

	cases := map[string]struct {
		params             s3.PutObjectInput
		conditions         []interface{}
		expectedStatusCode int
	}{
		"standard": {
			params: s3.PutObjectInput{},
		},
		"extra conditions, fail upload": {
			params: s3.PutObjectInput{},
			conditions: []interface{}{
				[]interface{}{
					// any number larger than the small sample
					"content-length-range",
					100000,
					200000,
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
			if err != nil {
				t.Fatalf("failed to load config, %v", err)
			}

			client := s3.NewFromConfig(cfg)

			// construct a put object
			presignerClient := s3.NewPresignClient(client)

			params := c.params
			if params.Key == nil {
				params.Key = aws.String(integrationtest.UniqueID())
			}
			params.Bucket = &setupMetadata.Buckets.Source.Name
			var presignRequest *s3.PresignedPostRequest
			if c.conditions != nil {
				presignRequest, err = presignerClient.PresignPostObject(ctx, &params, func(opts *s3.PresignPostOptions) {
					opts.Conditions = c.conditions
				})

			} else {
				presignRequest, err = presignerClient.PresignPostObject(ctx, &params)
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			resp, err := sendMultipartRequest(presignRequest.URL, presignRequest.Values, filePath)
			if err != nil {
				t.Fatalf("expect no error while sending HTTP request using presigned url, got %v", err)
			}

			defer resp.Body.Close()
			if c.expectedStatusCode != 0 {
				if resp.StatusCode != c.expectedStatusCode {
					t.Fatalf("expect status code %v, got %v", c.expectedStatusCode, resp.StatusCode)
				}
				// don't check the rest of the tests if there's a custom status code
				return
			} else {
				// expected result is 204 on POST requests
				if resp.StatusCode != http.StatusNoContent {
					t.Fatalf("failed to put S3 object, %d:%s", resp.StatusCode, resp.Status)
				}
			}

			// construct a get object
			getObjectInput := &s3.GetObjectInput{
				Bucket: params.Bucket,
				Key:    params.Key,
			}

			// This could be a regular GetObject call, but since we already have a presigner client available
			getRequest, err := presignerClient.PresignGetObject(ctx, getObjectInput)
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}

			resp, err = sendHTTPRequest(getRequest, nil)
			if err != nil {
				t.Errorf("expect no error while sending HTTP request using presigned url, got %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("failed to get S3 object, %d:%s", resp.StatusCode, resp.Status)
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("expect no error reading local file %v, got %v", filePath, err)
			}
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("expect no error reading response %v, got %v", resp.Body, err)
			}
			if !bytes.Equal(content, respBytes) {
				t.Fatalf("expect response body %v, got %v", content, resp.Body)
			}
		})
	}
}

func sendMultipartRequest(url string, fields map[string]string, filePath string) (*http.Response, error) {
	// Create a buffer to hold the multipart data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields
	for key, val := range fields {
		err := writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}

	// Add the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Always has to be named like this
	fileField := "file"
	part, err := writer.CreateFormFile(fileField, filePath)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Close the writer to finalize the multipart message
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	return client.Do(req)
}
