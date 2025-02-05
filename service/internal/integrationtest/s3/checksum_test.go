//go:build integration
// +build integration

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/logging"
)

type retryClient struct {
	isRetriedCall bool
	baseClient    aws.HTTPClient
}

type mockConnectionError struct{ err error }

func (m mockConnectionError) ConnectionError() bool {
	return true
}
func (m mockConnectionError) Error() string {
	return fmt.Sprintf("request error: %v", m.err)
}
func (m mockConnectionError) Unwrap() error {
	return m.err
}

func (c *retryClient) Do(req *http.Request) (*http.Response, error) {
	if !c.isRetriedCall {
		c.isRetriedCall = true
		return nil, mockConnectionError{}
	}
	return c.baseClient.Do(req)
}

func TestInteg_ObjectChecksums(t *testing.T) {
	cases := map[string]map[string]struct {
		disableHTTPS               bool
		retry                      bool
		requestChecksumCalculation aws.RequestChecksumCalculation
		params                     *s3.PutObjectInput

		expectErr                string
		expectReadErr            string
		expectLogged             string
		expectChecksumAlgorithms s3types.ChecksumAlgorithm
		expectPayload            []byte
		expectComputedChecksums  *s3.ComputedInputChecksumsMetadata
		expectAlgorithmsUsed     *s3.ChecksumValidationMetadata
	}{
		"seekable": {
			"no checksum algorithm passed": {
				params: &s3.PutObjectInput{
					Body: strings.NewReader("abc123"),
				},
				expectPayload: []byte("abc123"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32": "zwK7XA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32"},
				},
			},
			"calculate crc64nvme": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("abc123"),
					ChecksumAlgorithm: types.ChecksumAlgorithmCrc64nvme,
				},
				expectPayload: []byte("abc123"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC64NVME": "gwCmMgdcSIQ=",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"no checksum calculation": {
				params: &s3.PutObjectInput{
					Body: strings.NewReader("abc123"),
				},
				requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
				expectPayload:              []byte("abc123"),
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("yZRlqg=="),
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"preset crc64nvme checksum": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("Hello, precomputed checksum!"),
					ChecksumCRC64NVME: aws.String("uxBNEklueLQ="),
				},
				expectPayload: []byte("Hello, precomputed checksum!"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC64NVME": "uxBNEklueLQ=",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"wrong preset checksum": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("RZRlqg=="),
				},
				expectErr: "BadDigest",
			},
			"without TLS autofill header checksum": {
				disableHTTPS: true,
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"autofill trailing checksum": {
				retry: true,
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"content length preset": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ContentLength:     aws.Int64(11),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"with content encoding set": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ContentEncoding:   aws.String("gzip"),
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
		},
		"unseekable": {
			"no checksum algorithm passed": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte("abc123")),
				},
				expectPayload: []byte("abc123"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32": "zwK7XA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32"},
				},
			},
			"calculate crc64nvme": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("abc123")),
					ChecksumAlgorithm: types.ChecksumAlgorithmCrc64nvme,
				},
				expectPayload: []byte("abc123"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC64NVME": "gwCmMgdcSIQ=",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"no checksum calculation": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte("abc123")),
				},
				requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
				expectPayload:              []byte("abc123"),
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("yZRlqg=="),
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"preset crc64nvme checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("Hello, precomputed checksum!")),
					ChecksumCRC64NVME: aws.String("uxBNEklueLQ="),
				},
				expectPayload: []byte("Hello, precomputed checksum!"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC64NVME": "uxBNEklueLQ=",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"wrong preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("RZRlqg=="),
				},
				expectErr: "BadDigest",
			},
			"autofill trailing checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"without TLS": {
				disableHTTPS: true,
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectErr: "unseekable stream is not supported without TLS",
			},
			"content length preset": {
				params: &s3.PutObjectInput{
					Body:              ioutil.NopCloser(strings.NewReader("hello world")),
					ContentLength:     aws.Int64(11),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectPayload: []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"unknown content length": {
				params: &s3.PutObjectInput{
					Body:              ioutil.NopCloser(strings.NewReader("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectErr: "MissingContentLength",
			},
		},
		"nil body": {
			"no checksum algorithm passed": {
				params: &s3.PutObjectInput{},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32"},
				},
			},
			"no checksum calculation": {
				params:                     &s3.PutObjectInput{},
				requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("AAAAAA=="),
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"autofill checksum": {
				params: &s3.PutObjectInput{
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"without TLS autofill checksum": {
				disableHTTPS: true,
				params: &s3.PutObjectInput{
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
		},
		"empty body": {
			"no checksum algorithm passed": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte{}),
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32"},
				},
			},
			"no checksum calculation": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte{}),
				},
				requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC64NVME"},
				},
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte{}),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("AAAAAA=="),
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"autofill checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte{}),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"without TLS autofill checksum": {
				disableHTTPS: true,
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte{}),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "AAAAAA==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
		},
	}

	for groupName, cs := range cases {
		t.Run(groupName, func(t *testing.T) {
			for caseName, c := range cs {
				t.Run(caseName, func(t *testing.T) {
					c.params.Bucket = &setupMetadata.Buckets.Source.Name
					c.params.Key = aws.String(t.Name())

					ctx := context.Background()
					logger, logged := bufferLogger(t)
					s3Options := func(o *s3.Options) {
						o.Logger = logger
						o.EndpointOptions.DisableHTTPS = c.disableHTTPS
						if c.requestChecksumCalculation != 0 {
							o.RequestChecksumCalculation = c.requestChecksumCalculation
						}
					}

					if c.retry {
						opts := s3client.Options()
						opts.HTTPClient = &retryClient{
							baseClient: opts.HTTPClient,
						}
						s3client = s3.New(opts)
					}

					t.Logf("putting bucket: %q, object: %q", *c.params.Bucket, *c.params.Key)
					putResult, err := s3client.PutObject(ctx, c.params, s3Options)
					if err == nil && len(c.expectErr) != 0 {
						t.Fatalf("expect error %v, got none", c.expectErr)
					}
					if err != nil && len(c.expectErr) == 0 {
						t.Fatalf("expect no error, got %v", err)
					}
					if err != nil && !strings.Contains(err.Error(), c.expectErr) {
						t.Fatalf("expect error to contain %v, got %v", c.expectErr, err)
					}
					if c.expectErr != "" {
						return
					}
					// assert computed input checksums metadata
					computedChecksums, ok := s3.GetComputedInputChecksumsMetadata(putResult.ResultMetadata)
					if e, a := (c.expectComputedChecksums != nil), ok; e != a {
						t.Fatalf("expect computed checksum metadata %t, got %t, %v", e, a, computedChecksums)
					}
					if c.expectComputedChecksums != nil {
						if diff := cmpDiff(*c.expectComputedChecksums, computedChecksums); diff != "" {
							t.Errorf("expect computed checksum metadata match\n%s", diff)
						}
					}

					getResult, err := s3client.GetObject(ctx, &s3.GetObjectInput{
						Bucket: c.params.Bucket,
						Key:    c.params.Key,
					}, s3Options)
					if err != nil {
						t.Fatalf("expect no error, got %v", err)
					}

					actualPayload, err := ioutil.ReadAll(getResult.Body)
					if err == nil && len(c.expectReadErr) != 0 {
						t.Fatalf("expected read error: %v, got none", c.expectReadErr)
					}
					if err != nil && len(c.expectReadErr) == 0 {
						t.Fatalf("expect no read error, got %v", err)
					}
					if err != nil && !strings.Contains(err.Error(), c.expectReadErr) {
						t.Fatalf("expected read error %v to contain %v", err, c.expectReadErr)
					}
					if c.expectReadErr != "" {
						return
					}

					if diff := cmpDiff(string(c.expectPayload), string(actualPayload)); diff != "" {
						t.Errorf("expect payload match:\n%s", diff)
					}

					if err = getResult.Body.Close(); err != nil {
						t.Errorf("expect no close error, got %v", err)
					}

					// Only compare string values, since S3 can respond with
					// empty value Content-Encoding header.
					if e, a := aws.ToString(c.params.ContentEncoding), aws.ToString(getResult.ContentEncoding); e != a {
						t.Errorf("expect %v content encoding, got %v", e, a)
					}

					// assert checksum validation metadata
					algorithmsUsed, ok := s3.GetChecksumValidationMetadata(getResult.ResultMetadata)
					if e, a := (c.expectAlgorithmsUsed != nil), ok; e != a {
						t.Fatalf("expect algorithms used metadata %t, got %t, %v", e, a, algorithmsUsed)
					}
					if c.expectAlgorithmsUsed != nil {
						if diff := cmpDiff(*c.expectAlgorithmsUsed, algorithmsUsed); diff != "" {
							t.Errorf("expect algorithms used to match\n%s", diff)
						}
					}

					if c.expectLogged != "" {
						if e, a := c.expectLogged, logged.String(); !strings.Contains(a, e) {
							t.Errorf("expected %q logged in:\n%s", e, a)
						}
					}
				})
			}
		})
	}
}

func TestInteg_RequireChecksum(t *testing.T) {
	cases := map[string]struct {
		checksumAlgorithm       types.ChecksumAlgorithm
		expectComputedChecksums []string
	}{
		"no algorithm": {
			expectComputedChecksums: []string{"CRC32"},
		},
		"with algorithm": {
			checksumAlgorithm:       types.ChecksumAlgorithmCrc32c,
			expectComputedChecksums: []string{"CRC32C"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := s3client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: &setupMetadata.Buckets.Source.Name,
				Delete: &s3types.Delete{
					Objects: []s3types.ObjectIdentifier{
						{Key: aws.String(t.Name())},
					},
					Quiet: aws.Bool(true),
				},
				ChecksumAlgorithm: c.checksumAlgorithm,
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			// assert computed input checksums metadata
			computedChecksums, ok := s3.GetComputedInputChecksumsMetadata(result.ResultMetadata)
			if e, a := ok, (c.expectComputedChecksums != nil); e != a {
				t.Fatalf("expect computed checksum metadata %t, got %t, %v", e, a, computedChecksums)
			}
			if e, a := len(c.expectComputedChecksums), len(computedChecksums.ComputedChecksums); e != a {
				t.Errorf("expect %v computed checksums, got %v, %v", e, a, computedChecksums)
			}
			for _, e := range c.expectComputedChecksums {
				v, ok := computedChecksums.ComputedChecksums[e]
				if !ok {
					t.Errorf("expect %v algorithm to be computed", e)
				}
				if v == "" {
					t.Errorf("expect %v algorithm to have non-empty computed checksum", e)
				}
			}
		})
	}
}

func TestInteg_RequireChecksumWithoutRequestAlgorithmMember(t *testing.T) {
	params := &s3.PutBucketOwnershipControlsInput{
		Bucket: &setupMetadata.Buckets.Source.Name,
		OwnershipControls: &types.OwnershipControls{
			Rules: []types.OwnershipControlsRule{
				{
					ObjectOwnership: types.ObjectOwnershipBucketOwnerPreferred,
				},
			},
		},
	}

	t.Logf("putting bucket ownership control: %q", *params.Bucket)
	result, err := s3client.PutBucketOwnershipControls(context.Background(), params)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	computedChecksums, ok := s3.GetComputedInputChecksumsMetadata(result.ResultMetadata)
	if !ok {
		t.Fatalf("expect computed checksums metadata present, got %q", result)
	}

	expectComputedChecksums := s3.ComputedInputChecksumsMetadata{
		ComputedChecksums: map[string]string{
			"CRC32": "cK9COg==",
		},
	}
	if diff := cmpDiff(expectComputedChecksums, computedChecksums); diff != "" {
		t.Errorf("expect computed checksum metadata match: %s\n", diff)
	}
}

func bufferLogger(t *testing.T) (logging.Logger, *bytes.Buffer) {
	var logged bytes.Buffer

	logger := logging.LoggerFunc(
		func(classification logging.Classification, format string, v ...interface{}) {
			fmt.Fprintf(&logged, format, v...)
			t.Logf(format, v...)
		})

	return logger, &logged
}
