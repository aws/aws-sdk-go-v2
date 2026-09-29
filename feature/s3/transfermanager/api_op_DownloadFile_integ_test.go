//go:build integration

package transfermanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestInteg_DownloadFile(t *testing.T) {
	const mib = 1024 * 1024
	bucket := setupMetadata.Buckets.Source.Name

	ranges := []func(*Options){func(o *Options) { o.GetObjectType = types.GetObjectRanges }}
	parts := []func(*Options){func(o *Options) { o.GetObjectType = types.GetObjectParts }}

	// Above the direct I/O threshold, with an unaligned tail so the final write
	// is padded. Only ranges downloads on Linux use direct I/O.
	over := make([]byte, 96*mib+13)
	if _, err := rand.Read(over); err != nil {
		t.Fatal(err)
	}
	under := largeObjectBuf

	small := []byte("hello world")
	empty := []byte{}
	keys := map[string]string{}
	for name, body := range map[string][]byte{"small": small, "empty": empty, "under": under, "over": over} {
		key := UniqueID()
		if _, err := s3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader(body),
		}); err != nil {
			t.Fatalf("upload %s: %v", name, err)
		}
		keys[name] = key
	}

	unequal := bytes.Join([][]byte{
		bytes.Repeat([]byte{'A'}, 6*mib),
		bytes.Repeat([]byte{'B'}, 5*mib),
		bytes.Repeat([]byte{'C'}, 1*mib),
	}, nil)
	keys["unequal"] = UniqueID()
	uploadWithPartSizes(t, bucket, keys["unequal"], unequal, []int64{6 * mib, 5 * mib, 1 * mib})

	cases := map[string]downloadFileTestData{
		"small":                       {Key: keys["small"], ExpectBody: small},
		"empty ranges":                {Key: keys["empty"], OptFns: ranges, ExpectBody: empty},
		"empty parts":                 {Key: keys["empty"], OptFns: parts, ExpectBody: empty},
		"under threshold ranges":      {Key: keys["under"], OptFns: ranges, ExpectBody: under},
		"under threshold parts":       {Key: keys["under"], OptFns: parts, ExpectBody: under},
		"over threshold ranges":       {Key: keys["over"], OptFns: ranges, ExpectBody: over},
		"over threshold parts":        {Key: keys["over"], OptFns: parts, ExpectBody: over},
		"range under threshold":       {Key: keys["under"], OptFns: ranges, Range: "bytes=1-10485760", ExpectBody: under[1:10485761]},
		"range over threshold":        {Key: keys["over"], OptFns: ranges, Range: "bytes=10485760-94371852", ExpectBody: over[10485760:94371853]},
		"unequal part sizes ranges":   {Key: keys["unequal"], OptFns: ranges, ExpectBody: unequal},
		"unequal part sizes parts":    {Key: keys["unequal"], OptFns: parts, ExpectBody: unequal},
		"replaces larger file":        {Key: keys["under"], OptFns: ranges, Existing: bytes.Repeat([]byte{'x'}, 2*len(under)), ExpectBody: under},
		"missing key preserves file":  {Key: UniqueID(), OptFns: ranges, Existing: []byte("previous contents"), ExpectError: "NoSuchKey"},
		"missing key creates nothing": {Key: UniqueID(), OptFns: ranges, ExpectError: "NoSuchKey"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadFile(t, bucket, c)
		})
	}
}
