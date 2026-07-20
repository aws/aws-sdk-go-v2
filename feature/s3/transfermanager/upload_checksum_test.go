package transfermanager

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// checksumResolutionCases covers how the transfer manager resolves the upload
// checksum algorithm from its RequestChecksumCalculation and ChecksumAlgorithm
// options. Single-part and multipart uploads resolve identically.
var checksumResolutionCases = map[string]struct {
	requestChecksumCalculation aws.RequestChecksumCalculation
	optionsChecksumAlgorithm   types.ChecksumAlgorithm
	inputChecksumAlgorithm     types.ChecksumAlgorithm
	expect                     s3types.ChecksumAlgorithm
}{
	"default (unset) applies CRC32": {
		expect: s3types.ChecksumAlgorithmCrc32,
	},
	"WhenSupported applies CRC32": {
		requestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
		expect:                     s3types.ChecksumAlgorithmCrc32,
	},
	"WhenRequired sets no algorithm": {
		requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		expect:                     "",
	},
	"explicit input algorithm is honored": {
		requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		inputChecksumAlgorithm:     types.ChecksumAlgorithmSha256,
		expect:                     s3types.ChecksumAlgorithmSha256,
	},
	"explicit options algorithm is honored": {
		requestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		optionsChecksumAlgorithm:   types.ChecksumAlgorithmSha1,
		expect:                     s3types.ChecksumAlgorithmSha1,
	},
}

func TestUploadObjectSinglePartChecksum(t *testing.T) {
	for name, c := range checksumResolutionCases {
		t.Run(name, func(t *testing.T) {
			client, invocations, args := s3testing.NewUploadLoggingClient(nil)
			mgr := New(client, func(o *Options) {
				o.RequestChecksumCalculation = c.requestChecksumCalculation
				o.ChecksumAlgorithm = c.optionsChecksumAlgorithm
			})

			_, err := mgr.UploadObject(context.Background(), &UploadObjectInput{
				Bucket:            aws.String("Bucket"),
				Key:               aws.String("Key"),
				Body:              bytes.NewReader(buf2MB),
				ChecksumAlgorithm: c.inputChecksumAlgorithm,
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if diff := cmpDiff([]string{"PutObject"}, *invocations); len(diff) > 0 {
				t.Fatal(diff)
			}

			put := (*args)[0].(*s3.PutObjectInput)
			if e, a := c.expect, put.ChecksumAlgorithm; e != a {
				t.Errorf("expect PutObject checksum algorithm %q, got %q", e, a)
			}
		})
	}
}

func TestUploadObjectMultipartChecksum(t *testing.T) {
	for name, c := range checksumResolutionCases {
		t.Run(name, func(t *testing.T) {
			client, invocations, args := s3testing.NewUploadLoggingClient(nil)
			mgr := New(client, func(o *Options) {
				o.RequestChecksumCalculation = c.requestChecksumCalculation
				o.ChecksumAlgorithm = c.optionsChecksumAlgorithm
			})

			_, err := mgr.UploadObject(context.Background(), &UploadObjectInput{
				Bucket:            aws.String("Bucket"),
				Key:               aws.String("Key"),
				Body:              bytes.NewReader(buf20MB),
				ChecksumAlgorithm: c.inputChecksumAlgorithm,
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			expectOps := []string{"CreateMultipartUpload", "UploadPart", "UploadPart", "UploadPart", "CompleteMultipartUpload"}
			if diff := cmpDiff(expectOps, *invocations); len(diff) > 0 {
				t.Fatal(diff)
			}

			cmu := (*args)[0].(*s3.CreateMultipartUploadInput)
			if e, a := c.expect, cmu.ChecksumAlgorithm; e != a {
				t.Errorf("expect CreateMultipartUpload checksum algorithm %q, got %q", e, a)
			}

			for i := 1; i <= 3; i++ {
				part := (*args)[i].(*s3.UploadPartInput)
				if e, a := c.expect, part.ChecksumAlgorithm; e != a {
					t.Errorf("expect UploadPart %d checksum algorithm %q, got %q", i, e, a)
				}
			}
		})
	}
}

type cannedHTTPClient struct{ puts []http.Header }

func (c *cannedHTTPClient) Do(r *http.Request) (*http.Response, error) {
	if r.Method == http.MethodPut {
		c.puts = append(c.puts, r.Header.Clone())
	}
	if r.Body != nil {
		io.Copy(io.Discard, r.Body)
	}
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Etag": {`"abc"`}},
		Body:       io.NopCloser(bytes.NewReader(nil)),
		Request:    r,
	}, nil
}

// TestUploadObjectChecksumClientPrecedence drives a real s3.Client through a
// canned HTTP transport (no network) to verify that the transfer manager's
// RequestChecksumCalculation takes precedence over the S3 client's setting on
// the wire, as required by the s3-transfer-manager SEP.
func TestUploadObjectChecksumClientPrecedence(t *testing.T) {
	cases := map[string]struct {
		tmSetting     aws.RequestChecksumCalculation
		clientSetting aws.RequestChecksumCalculation
		expectCRC32   bool
	}{
		"tm WhenRequired overrides client WhenSupported": {
			tmSetting:     aws.RequestChecksumCalculationWhenRequired,
			clientSetting: aws.RequestChecksumCalculationWhenSupported,
			expectCRC32:   false,
		},
		"tm WhenSupported overrides client WhenRequired": {
			tmSetting:     aws.RequestChecksumCalculationWhenSupported,
			clientSetting: aws.RequestChecksumCalculationWhenRequired,
			expectCRC32:   true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cap := &cannedHTTPClient{}
			cfg, err := config.LoadDefaultConfig(context.Background(),
				config.WithRegion("us-west-2"),
				config.WithHTTPClient(cap),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET", "")),
			)
			if err != nil {
				t.Fatalf("load config: %v", err)
			}

			s3c := s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.RequestChecksumCalculation = c.clientSetting
			})
			mgr := New(s3c, func(o *Options) {
				o.RequestChecksumCalculation = c.tmSetting
			})

			if _, err := mgr.UploadObject(context.Background(), &UploadObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
				Body:   bytes.NewReader([]byte("hello")),
			}); err != nil {
				t.Fatalf("upload: %v", err)
			}

			if len(cap.puts) == 0 {
				t.Fatal("no PUT request captured")
			}
			h := cap.puts[len(cap.puts)-1]
			gotCRC32 := h.Get("X-Amz-Sdk-Checksum-Algorithm") == "CRC32"
			if gotCRC32 != c.expectCRC32 {
				t.Errorf("expectCRC32=%v, got x-amz-sdk-checksum-algorithm=%q trailer=%q",
					c.expectCRC32, h.Get("X-Amz-Sdk-Checksum-Algorithm"), h.Get("X-Amz-Trailer"))
			}
		})
	}
}

type captureFinalizeHandler struct {
	ctx context.Context
}

func (h *captureFinalizeHandler) HandleFinalize(
	ctx context.Context, _ smithymiddleware.FinalizeInput,
) (smithymiddleware.FinalizeOutput, smithymiddleware.Metadata, error) {
	h.ctx = ctx
	return smithymiddleware.FinalizeOutput{}, smithymiddleware.Metadata{}, nil
}

// TestSetS3ExpressDefaultChecksum verifies that the finalize middleware still
// forces CRC32 for S3 Express uploads even when the general checksum default is
// suppressed, and that it is a no-op for non-Express backends.
func TestSetS3ExpressDefaultChecksum(t *testing.T) {
	newInput := func() smithymiddleware.FinalizeInput {
		return smithymiddleware.FinalizeInput{Request: smithyhttp.NewStackRequest()}
	}

	t.Run("non-Express backend is a no-op", func(t *testing.T) {
		next := &captureFinalizeHandler{}
		_, _, err := (&setS3ExpressDefaultChecksum{}).HandleFinalize(context.Background(), newInput(), next)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if a := internalcontext.GetChecksumInputAlgorithm(next.ctx); a != "" {
			t.Errorf("expect no checksum algorithm to be set, got %q", a)
		}
	})

	t.Run("Express backend forces CRC32 when none set", func(t *testing.T) {
		ctx := internalcontext.SetS3Backend(context.Background(), internalcontext.S3BackendS3Express)
		next := &captureFinalizeHandler{}
		_, _, err := (&setS3ExpressDefaultChecksum{}).HandleFinalize(ctx, newInput(), next)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if e, a := string(s3types.ChecksumAlgorithmCrc32), internalcontext.GetChecksumInputAlgorithm(next.ctx); e != a {
			t.Errorf("expect checksum algorithm %q, got %q", e, a)
		}
	})

	t.Run("Express backend preserves an explicit algorithm", func(t *testing.T) {
		ctx := internalcontext.SetS3Backend(context.Background(), internalcontext.S3BackendS3Express)
		ctx = internalcontext.SetChecksumInputAlgorithm(ctx, string(s3types.ChecksumAlgorithmSha256))
		next := &captureFinalizeHandler{}
		_, _, err := (&setS3ExpressDefaultChecksum{}).HandleFinalize(ctx, newInput(), next)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

		if e, a := string(s3types.ChecksumAlgorithmSha256), internalcontext.GetChecksumInputAlgorithm(next.ctx); e != a {
			t.Errorf("expect checksum algorithm %q preserved, got %q", e, a)
		}
	})
}
