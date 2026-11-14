//go:build integration

package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

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

// TestInteg_GetObjectMultipart downloads a multipart object with equal part
// sizes in both modes.
func TestInteg_GetObjectMultipart(t *testing.T) {
	const mib = 1024 * 1024
	partSizes := []int64{8 * mib, 8 * mib, 4 * mib}

	cases := map[string]getObjectTestData{
		"parts get": {
			Body:       bytes.NewReader(largeObjectBuf),
			ExpectBody: largeObjectBuf,
			PartSizes:  partSizes,
		},
		"ranges get": {
			Body:       bytes.NewReader(largeObjectBuf),
			ExpectBody: largeObjectBuf,
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

// TestInteg_GetObjectEarlyClose closes the body partway through a download.
// Close must return promptly, and Read must then fail without blocking.
func TestInteg_GetObjectEarlyClose(t *testing.T) {
	const mib = 1024 * 1024
	ctx := context.Background()
	bucket := setupMetadata.Buckets.Source.Name
	key := UniqueID()
	uploadMultipartExact(t, ctx, bucket, key, largeObjectBuf, []int64{8 * mib, 8 * mib, 4 * mib})

	out, err := s3TransferManagerClient.GetObject(ctx, &GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("expect no error getting object, got %v", err)
	}

	head := make([]byte, mib)
	if _, err := io.ReadFull(out.Body, head); err != nil {
		t.Fatalf("expect no error reading first MiB, got %v", err)
	}
	if !bytes.Equal(head, largeObjectBuf[:mib]) {
		t.Fatalf("first MiB does not match the uploaded bytes")
	}

	closed := make(chan error, 1)
	go func() { closed <- out.Body.(io.Closer).Close() }()
	select {
	case err := <-closed:
		if err != nil {
			t.Fatalf("expect no error closing body, got %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatalf("Close did not return within 30s")
	}

	// Read may still return parts that were already downloaded, but it must
	// fail or return the rest of the object.
	readErr := make(chan error, 1)
	go func() {
		rest, err := io.ReadAll(out.Body)
		if err == nil && !bytes.Equal(append(head, rest...), largeObjectBuf) {
			err = fmt.Errorf("read %d bytes after Close without an error", len(rest))
		}
		readErr <- err
	}()
	select {
	case err := <-readErr:
		if err != nil && !errors.Is(err, errReadClosed) {
			t.Fatalf("expect read error %q after Close, got %v", errReadClosed, err)
		}
	case <-time.After(30 * time.Second):
		t.Fatalf("Read after Close did not return within 30s")
	}
}
