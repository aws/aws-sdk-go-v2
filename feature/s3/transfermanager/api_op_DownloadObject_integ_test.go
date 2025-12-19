//go:build integration
// +build integration

package transfermanager

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"strings"
	"testing"
)

func TestInteg_DownloadObject(t *testing.T) {
	cases := map[string]downloadObjectTestData{
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
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadObject(t, setupMetadata.Buckets.Source.Name, c)
		})
	}
}
