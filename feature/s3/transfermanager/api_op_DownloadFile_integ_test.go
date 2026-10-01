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

	// Large enough for sustained direct I/O, with an unaligned tail so the
	// final write is padded. Direct I/O is only used on Linux, when opted in,
	// for ranges downloads.
	large := make([]byte, 96*mib+13)
	if _, err := rand.Read(large); err != nil {
		t.Fatal(err)
	}
	medium := largeObjectBuf

	small := []byte("hello world")
	empty := []byte{}
	keys := map[string]string{}
	for name, body := range map[string][]byte{"small": small, "empty": empty, "medium": medium, "large": large} {
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
		"small ranges direct":         {Key: keys["small"], OptFns: ranges, DirectIO: true, ExpectBody: small},
		"empty ranges":                {Key: keys["empty"], OptFns: ranges, ExpectBody: empty},
		"empty ranges direct":         {Key: keys["empty"], OptFns: ranges, DirectIO: true, ExpectBody: empty},
		"empty parts":                 {Key: keys["empty"], OptFns: parts, ExpectBody: empty},
		"medium ranges":               {Key: keys["medium"], OptFns: ranges, ExpectBody: medium},
		"medium ranges direct":        {Key: keys["medium"], OptFns: ranges, DirectIO: true, ExpectBody: medium},
		"medium parts":                {Key: keys["medium"], OptFns: parts, ExpectBody: medium},
		"large ranges":                {Key: keys["large"], OptFns: ranges, ExpectBody: large},
		"large ranges direct":         {Key: keys["large"], OptFns: ranges, DirectIO: true, ExpectBody: large},
		"large parts":                 {Key: keys["large"], OptFns: parts, ExpectBody: large},
		"large parts direct":          {Key: keys["large"], OptFns: parts, DirectIO: true, ExpectBody: large},
		"large default direct":        {Key: keys["large"], DirectIO: true, ExpectBody: large},
		"range medium":                {Key: keys["medium"], OptFns: ranges, Range: "bytes=1-10485760", ExpectBody: medium[1:10485761]},
		"range medium direct":         {Key: keys["medium"], OptFns: ranges, DirectIO: true, Range: "bytes=1-10485760", ExpectBody: medium[1:10485761]},
		"range large direct":          {Key: keys["large"], OptFns: ranges, DirectIO: true, Range: "bytes=10485760-94371852", ExpectBody: large[10485760:94371853]},
		"unequal part sizes ranges":   {Key: keys["unequal"], OptFns: ranges, ExpectBody: unequal},
		"unequal part sizes direct":   {Key: keys["unequal"], OptFns: ranges, DirectIO: true, ExpectBody: unequal},
		"unequal part sizes parts":    {Key: keys["unequal"], OptFns: parts, ExpectBody: unequal},
		"replaces larger file":        {Key: keys["medium"], OptFns: ranges, Existing: bytes.Repeat([]byte{'x'}, 2*len(medium)), ExpectBody: medium},
		"replaces larger file direct": {Key: keys["medium"], OptFns: ranges, DirectIO: true, Existing: bytes.Repeat([]byte{'x'}, 2*len(medium)), ExpectBody: medium},
		"missing key preserves file":  {Key: UniqueID(), OptFns: ranges, Existing: []byte("previous contents"), ExpectError: "NoSuchKey"},
		"missing key creates nothing": {Key: UniqueID(), OptFns: ranges, ExpectError: "NoSuchKey"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			testDownloadFile(t, bucket, c)
		})
	}
}
