//go:build integration

package transfermanager

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
)

func TestInteg_GetObject(t *testing.T) {
	cases := map[string]getObjectTestData{
		"part get seekable body":     {Body: strings.NewReader("hello world"), ExpectBody: []byte("hello world")},
		"part get empty string body": {Body: strings.NewReader(""), ExpectBody: []byte("")},
		"part get multipart body":    {Body: bytes.NewReader(largeObjectBuf), ExpectBody: largeObjectBuf},
		"range get seekable body": {
			Body:       strings.NewReader("hello world"),
			ExpectBody: []byte("hello world"),
			OptFns: []func(*Options){
				func(opt *Options) {
					opt.GetObjectType = types.GetObjectRanges
				},
			},
		},
		"range get empty string body": {
			Body:       strings.NewReader(""),
			ExpectBody: []byte(""),
			OptFns: []func(*Options){
				func(opt *Options) {
					opt.GetObjectType = types.GetObjectRanges
				},
			},
		},
		"range get multipart body": {
			Body:       bytes.NewReader(largeObjectBuf),
			ExpectBody: largeObjectBuf,
			OptFns: []func(*Options){
				func(opt *Options) {
					opt.GetObjectType = types.GetObjectRanges
				},
			},
		},
		"range get multipart body with range input": {
			Body:       bytes.NewReader(largeObjectBuf),
			Range:      "bytes=8388608-16777215",
			ExpectBody: largeObjectBuf[8388608:16777216],
			OptFns: []func(*Options){
				func(opt *Options) {
					opt.GetObjectType = types.GetObjectRanges
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testGetObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}

// TestInteg_GetObject_UnequalSize verifies that streaming a multipart object
// whose parts have unequal sizes returns exactly the uploaded content, in both
// parts and ranges modes (#3526).
func TestInteg_GetObject_UnequalSize(t *testing.T) {
	const mib = 1024 * 1024
	// Distinct byte patterns per part so any misplacement corrupts the result.
	// All parts but the last are >= 5MB to satisfy the S3 minimum part size.
	partA := bytes.Repeat([]byte{'A'}, 6*mib)
	partB := bytes.Repeat([]byte{'B'}, 5*mib)
	partC := bytes.Repeat([]byte{'C'}, 1*mib)
	body := bytes.Join([][]byte{partA, partB, partC}, nil)
	partSizes := []int64{6 * mib, 5 * mib, 1 * mib}

	cases := map[string]getObjectTestData{
		"parts get unequal part sizes": {
			Body:       bytes.NewReader(body),
			ExpectBody: body,
			PartSizes:  partSizes,
		},
		"ranges get unequal part sizes": {
			Body:       bytes.NewReader(body),
			ExpectBody: body,
			PartSizes:  partSizes,
			OptFns: []func(*Options){
				func(opt *Options) {
					opt.GetObjectType = types.GetObjectRanges
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testGetObjectWithChangingPartSize(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}

// TestInteg_GetObject_RandomizedReassembly is the property-based durability
// test for https://github.com/aws/aws-sdk-go-v2/issues/3526, driving the GetObject (streaming reader) path.
// It uploads randomized multipart layouts and asserts the streamed bytes equal
// the uploaded bytes byte for byte. See runRandomizedReassembly in
// setup_integ_test.go. Reproduce a failure with TM_REASSEMBLY_SEED.
func TestInteg_GetObject_RandomizedReassembly(t *testing.T) {
	runRandomizedReassembly(t, "GetObject", func(t *testing.T, ctx context.Context, bucket, key string, c reassemblyCase) []byte {
		out, err := s3TransferManagerClient.GetObject(ctx, &GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  c.rangePtr(),
		}, c.optFns()...)
		if err != nil {
			t.Fatalf("GetObject: %v", err)
		}
		got, err := io.ReadAll(out.Body)
		if err != nil {
			t.Fatalf("read GetObject body: %v", err)
		}
		return got
	})
}

// TestInteg_GetObject_ZeroLengthRead uploads a 40 MiB (5-part) multipart object
// and verifies that a zero-length Read on the GetObject body returns (0, nil)
// without deadlocking, and that a subsequent io.ReadAll still streams the whole
// object back byte for byte (GH #3641). Both GetObject branches are covered
// since each has its own dispatch path.
func TestInteg_GetObject_ZeroLengthRead(t *testing.T) {
	// largeObjectBuf is 20 MiB; doubled it is 40 MiB, which is 5 parts at the
	// 8 MiB default part size and matches the object shape from the issue.
	body := append(append([]byte(nil), largeObjectBuf...), largeObjectBuf...)

	cases := map[string]getObjectTestData{
		"parts": {Body: bytes.NewReader(body)},
		"ranges": {
			Body: bytes.NewReader(body),
			OptFns: []func(*Options){
				func(o *Options) { o.GetObjectType = types.GetObjectRanges },
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testGetObjectZeroLengthRead(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
