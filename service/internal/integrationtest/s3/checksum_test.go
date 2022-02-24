//go:build integration
// +build integration

package s3

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/logging"
	"github.com/google/go-cmp/cmp"
)

func TestInteg_ObjectChecksums(t *testing.T) {
	cases := map[string]map[string]struct {
		disableHTTPS bool
		params       *s3.PutObjectInput
		expectErr    string

		getObjectChecksumMode    s3types.ChecksumMode
		expectReadErr            string
		expectLogged             string
		expectChecksumAlgorithms s3types.ChecksumAlgorithm
		expectPayload            []byte
		expectComputedChecksums  *s3.ComputedInputChecksumsMetadata
		expectAlgorithmsUsed     *s3.ChecksumValidationMetadata
	}{
		"seekable": {
			"no checksum": {
				params: &s3.PutObjectInput{
					Body: strings.NewReader("abc123"),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("abc123"),
				expectLogged:          "Response has no supported checksum.",
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("yZRlqg=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"wrong preset checksum": {
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("RZRlqg=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectErr:             "BadDigest",
			},
			"without TLS autofill header checksum": {
				disableHTTPS: true,
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
				params: &s3.PutObjectInput{
					Body:              strings.NewReader("hello world"),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
					ContentLength:     11,
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
			"no checksum": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte("abc123")),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("abc123"),
				expectLogged:          "Response has no supported checksum.",
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("yZRlqg=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
				expectComputedChecksums: &s3.ComputedInputChecksumsMetadata{
					ComputedChecksums: map[string]string{
						"CRC32C": "yZRlqg==",
					},
				},
				expectAlgorithmsUsed: &s3.ChecksumValidationMetadata{
					AlgorithmsUsed: []string{"CRC32C"},
				},
			},
			"wrong preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("RZRlqg=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectErr:             "BadDigest",
			},
			"autofill trailing checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte("hello world")),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
					ContentLength:     11,
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectPayload:         []byte("hello world"),
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
			"no checksum": {
				params:                &s3.PutObjectInput{},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectLogged:          "Response has no supported checksum.",
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("AAAAAA=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
			"no checksum": {
				params: &s3.PutObjectInput{
					Body: bytes.NewBuffer([]byte{}),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
				expectLogged:          "Response has no supported checksum.",
			},
			"preset checksum": {
				params: &s3.PutObjectInput{
					Body:              bytes.NewBuffer([]byte{}),
					ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
					ChecksumCRC32C:    aws.String("AAAAAA=="),
				},
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
				getObjectChecksumMode: s3types.ChecksumModeEnabled,
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
					if e, a := ok, (c.expectComputedChecksums != nil); e != a {
						t.Fatalf("expect computed checksum metadata %t, got %t, %v", e, a, computedChecksums)
					}
					if c.expectComputedChecksums != nil {
						if diff := cmp.Diff(*c.expectComputedChecksums, computedChecksums); diff != "" {
							t.Errorf("expect computed checksum metadata match\n%s", diff)
						}
					}

					getResult, err := s3client.GetObject(ctx, &s3.GetObjectInput{
						Bucket:       c.params.Bucket,
						Key:          c.params.Key,
						ChecksumMode: c.getObjectChecksumMode,
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

					if diff := cmp.Diff(string(c.expectPayload), string(actualPayload)); diff != "" {
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
					if e, a := ok, (c.expectAlgorithmsUsed != nil); e != a {
						t.Fatalf("expect algorithms used metadata %t, got %t, %v", e, a, algorithmsUsed)
					}
					if c.expectAlgorithmsUsed != nil {
						if diff := cmp.Diff(*c.expectAlgorithmsUsed, algorithmsUsed); diff != "" {
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
			expectComputedChecksums: []string{"MD5"},
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
					Quiet: true,
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

func bufferLogger(t *testing.T) (logging.Logger, *bytes.Buffer) {
	var logged bytes.Buffer

	logger := logging.LoggerFunc(
		func(classification logging.Classification, format string, v ...interface{}) {
			fmt.Fprintf(&logged, format, v...)
			t.Logf(format, v...)
		})

	return logger, &logged
}
