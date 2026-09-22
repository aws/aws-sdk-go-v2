//go:build integration

package transfermanager

import (
	"bytes"
	"strings"
	"testing"

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
