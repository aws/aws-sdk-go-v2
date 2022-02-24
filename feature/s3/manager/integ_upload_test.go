//go:build integration
// +build integration

package manager_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var integBuf12MB = make([]byte, 1024*1024*12)
var integMD512MB = fmt.Sprintf("%x", md5.Sum(integBuf12MB))

func hexEncodeSum(sum []byte) string {
	sumHex := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(sumHex, sum)
	return string(sumHex)
}

func base64EncodeSum(sum []byte) string {
	sum64 := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
	base64.StdEncoding.Encode(sum64, sum)
	return string(sum64)
}

func base64Sum(h hash.Hash, b []byte) string {
	h.Write(b)
	return base64EncodeSum(h.Sum(nil))
}
func hexSum(h hash.Hash, b []byte) string {
	h.Write(b)
	return hexEncodeSum(h.Sum(nil))
}

func base64StringDecode(v string) []byte {
	vv, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		panic(err.Error())
	}
	return vv
}
func hexStringDecode(v string) []byte {
	vv, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		panic(err.Error())
	}
	return vv
}

func base64SumOfSums(h hash.Hash, sums []string) string {
	for _, v := range sums {
		h.Write(base64StringDecode(v))
	}
	return base64EncodeSum(h.Sum(nil)) + "-" + strconv.Itoa(len(sums))
}

func hexSumOfSums(h hash.Hash, sums []string) string {
	for _, v := range sums {
		h.Write(hexStringDecode(unquote(v)))
	}
	return hexEncodeSum(h.Sum(nil)) + "-" + strconv.Itoa(len(sums))
}

func TestInteg_UploadConcurrently(t *testing.T) {
	key := "12mb-1"
	mgr := manager.NewUploader(client)
	out, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:            bucketName,
		Key:               &key,
		Body:              bytes.NewReader(integBuf12MB),
		ChecksumAlgorithm: s3types.ChecksumAlgorithmCrc32c,
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if len(out.UploadID) == 0 {
		t.Errorf("expect upload ID but was empty")
	}

	re := regexp.MustCompile(`^https?://.+/` + key + `$`)
	if e, a := re.String(), out.Location; !re.MatchString(a) {
		t.Errorf("expect %s to match URL regexp %q, did not", e, a)
	}

	validate(t, key, integMD512MB)
}

func unquote(v string) string {
	log.Printf("unquote: %v", v)
	vv, err := strconv.Unquote(v)
	if err != nil {
		// Unquote returns error if string doesn't contain quotes
		if err == strconv.ErrSyntax {
			return v
		}
		panic(err.Error())
	}
	return vv
}

func TestInteg_UploadPresetChecksum(t *testing.T) {
	singlePartBytes := integBuf12MB[0:manager.DefaultUploadPartSize]
	singlePartCRC32 := base64Sum(crc32.NewIEEE(), singlePartBytes)
	singlePartCRC32C := base64Sum(crc32.New(crc32.MakeTable(crc32.Castagnoli)), singlePartBytes)
	singlePartSHA1 := base64Sum(sha1.New(), singlePartBytes)
	singlePartSHA256 := base64Sum(sha256.New(), singlePartBytes)
	singlePartMD5 := base64Sum(md5.New(), singlePartBytes)
	singlePartETag := fmt.Sprintf("%q", hexSum(md5.New(), singlePartBytes))

	multiPartTailBytes := integBuf12MB[manager.DefaultUploadPartSize*2:]
	multiPartTailCRC32 := base64Sum(crc32.NewIEEE(), multiPartTailBytes)
	multiPartTailCRC32C := base64Sum(crc32.New(crc32.MakeTable(crc32.Castagnoli)), multiPartTailBytes)
	multiPartTailSHA1 := base64Sum(sha1.New(), multiPartTailBytes)
	multiPartTailSHA256 := base64Sum(sha256.New(), multiPartTailBytes)
	multiPartTailETag := fmt.Sprintf("%q", hexSum(md5.New(), multiPartTailBytes))

	multiPartBytes := integBuf12MB
	multiPartCRC32 := base64SumOfSums(crc32.NewIEEE(), []string{singlePartCRC32, singlePartCRC32, multiPartTailCRC32})
	multiPartCRC32C := base64SumOfSums(crc32.New(crc32.MakeTable(crc32.Castagnoli)), []string{singlePartCRC32C, singlePartCRC32C, multiPartTailCRC32C})
	multiPartSHA1 := base64SumOfSums(sha1.New(), []string{singlePartSHA1, singlePartSHA1, multiPartTailSHA1})
	multiPartSHA256 := base64SumOfSums(sha256.New(), []string{singlePartSHA256, singlePartSHA256, multiPartTailSHA256})
	multiPartETag := `"4e982d58b6c2ce178ae042c23f9bca6e-3"` // Not obvious how this is computed

	cases := map[string]map[string]struct {
		algorithm            s3types.ChecksumAlgorithm
		payload              io.Reader
		checksumCRC32        string
		checksumCRC32C       string
		checksumSHA1         string
		checksumSHA256       string
		contentMD5           string
		expectParts          []s3types.CompletedPart
		expectChecksumCRC32  string
		expectChecksumCRC32C string
		expectChecksumSHA1   string
		expectChecksumSHA256 string
		expectETag           string
	}{
		"auto single part": {
			"no checksum": {
				payload:    bytes.NewReader(singlePartBytes),
				expectETag: singlePartETag,
			},
			"CRC32": {
				algorithm:           s3types.ChecksumAlgorithmCrc32,
				payload:             bytes.NewReader(singlePartBytes),
				expectChecksumCRC32: singlePartCRC32,
				expectETag:          singlePartETag,
			},
			"CRC32C": {
				algorithm:            s3types.ChecksumAlgorithmCrc32c,
				payload:              bytes.NewReader(singlePartBytes),
				expectChecksumCRC32C: singlePartCRC32C,
				expectETag:           singlePartETag,
			},
			"SHA1": {
				algorithm:          s3types.ChecksumAlgorithmSha1,
				payload:            bytes.NewReader(singlePartBytes),
				expectChecksumSHA1: singlePartSHA1,
				expectETag:         singlePartETag,
			},
			"SHA256": {
				algorithm:            s3types.ChecksumAlgorithmSha256,
				payload:              bytes.NewReader(singlePartBytes),
				expectChecksumSHA256: singlePartSHA256,
				expectETag:           singlePartETag,
			},
		},
		"preset single part": {
			"CRC32": {
				payload:             bytes.NewReader(singlePartBytes),
				checksumCRC32:       singlePartCRC32,
				expectChecksumCRC32: singlePartCRC32,
				expectETag:          singlePartETag,
			},
			"CRC32C": {
				payload:              bytes.NewReader(singlePartBytes),
				checksumCRC32C:       singlePartCRC32C,
				expectChecksumCRC32C: singlePartCRC32C,
				expectETag:           singlePartETag,
			},
			"SHA1": {
				payload:            bytes.NewReader(singlePartBytes),
				checksumSHA1:       singlePartSHA1,
				expectChecksumSHA1: singlePartSHA1,
				expectETag:         singlePartETag,
			},
			"SHA256": {
				payload:              bytes.NewReader(singlePartBytes),
				checksumSHA256:       singlePartSHA256,
				expectChecksumSHA256: singlePartSHA256,
				expectETag:           singlePartETag,
			},
			"MD5": {
				payload:    bytes.NewReader(singlePartBytes),
				contentMD5: singlePartMD5,
				expectETag: singlePartETag,
			},
		},
		"auto multipart part": {
			"no checksum": {
				payload: bytes.NewReader(multiPartBytes),
				expectParts: []s3types.CompletedPart{
					{
						ETag:       aws.String(singlePartETag),
						PartNumber: 1,
					},
					{
						ETag:       aws.String(singlePartETag),
						PartNumber: 2,
					},
					{
						ETag:       aws.String(multiPartTailETag),
						PartNumber: 3,
					},
				},
				expectETag: multiPartETag,
			},
			"CRC32": {
				algorithm: s3types.ChecksumAlgorithmCrc32,
				payload:   bytes.NewReader(multiPartBytes),
				expectParts: []s3types.CompletedPart{
					{
						ChecksumCRC32: aws.String(singlePartCRC32),
						ETag:          aws.String(singlePartETag),
						PartNumber:    1,
					},
					{
						ChecksumCRC32: aws.String(singlePartCRC32),
						ETag:          aws.String(singlePartETag),
						PartNumber:    2,
					},
					{
						ChecksumCRC32: aws.String(multiPartTailCRC32),
						ETag:          aws.String(multiPartTailETag),
						PartNumber:    3,
					},
				},
				expectChecksumCRC32: multiPartCRC32,
				expectETag:          multiPartETag,
			},
			"CRC32C": {
				algorithm: s3types.ChecksumAlgorithmCrc32c,
				payload:   bytes.NewReader(multiPartBytes),
				expectParts: []s3types.CompletedPart{
					{
						ChecksumCRC32C: aws.String(singlePartCRC32C),
						ETag:           aws.String(singlePartETag),
						PartNumber:     1,
					},
					{
						ChecksumCRC32C: aws.String(singlePartCRC32C),
						ETag:           aws.String(singlePartETag),
						PartNumber:     2,
					},
					{
						ChecksumCRC32C: aws.String(multiPartTailCRC32C),
						ETag:           aws.String(multiPartTailETag),
						PartNumber:     3,
					},
				},
				expectChecksumCRC32C: multiPartCRC32C,
				expectETag:           multiPartETag,
			},
			"SHA1": {
				algorithm: s3types.ChecksumAlgorithmSha1,
				payload:   bytes.NewReader(multiPartBytes),
				expectParts: []s3types.CompletedPart{
					{
						ChecksumSHA1: aws.String(singlePartSHA1),
						ETag:         aws.String(singlePartETag),
						PartNumber:   1,
					},
					{
						ChecksumSHA1: aws.String(singlePartSHA1),
						ETag:         aws.String(singlePartETag),
						PartNumber:   2,
					},
					{
						ChecksumSHA1: aws.String(multiPartTailSHA1),
						ETag:         aws.String(multiPartTailETag),
						PartNumber:   3,
					},
				},
				expectChecksumSHA1: multiPartSHA1,
				expectETag:         multiPartETag,
			},
			"SHA256": {
				algorithm: s3types.ChecksumAlgorithmSha256,
				payload:   bytes.NewReader(multiPartBytes),
				expectParts: []s3types.CompletedPart{
					{
						ChecksumSHA256: aws.String(singlePartSHA256),
						ETag:           aws.String(singlePartETag),
						PartNumber:     1,
					},
					{
						ChecksumSHA256: aws.String(singlePartSHA256),
						ETag:           aws.String(singlePartETag),
						PartNumber:     2,
					},
					{
						ChecksumSHA256: aws.String(multiPartTailSHA256),
						ETag:           aws.String(multiPartTailETag),
						PartNumber:     3,
					},
				},
				expectChecksumSHA256: multiPartSHA256,
				expectETag:           multiPartETag,
			},
		},
		"preset multipart part": {
			"CRC32": {
				algorithm:     s3types.ChecksumAlgorithmCrc32,
				payload:       bytes.NewReader(multiPartBytes),
				checksumCRC32: multiPartCRC32,
				expectParts: []s3types.CompletedPart{
					{
						ChecksumCRC32: aws.String(singlePartCRC32),
						ETag:          aws.String(singlePartETag),
						PartNumber:    1,
					},
					{
						ChecksumCRC32: aws.String(singlePartCRC32),
						ETag:          aws.String(singlePartETag),
						PartNumber:    2,
					},
					{
						ChecksumCRC32: aws.String(multiPartTailCRC32),
						ETag:          aws.String(multiPartTailETag),
						PartNumber:    3,
					},
				},
				expectChecksumCRC32: multiPartCRC32,
				expectETag:          multiPartETag,
			},
			"CRC32C": {
				algorithm:      s3types.ChecksumAlgorithmCrc32c,
				payload:        bytes.NewReader(multiPartBytes),
				checksumCRC32C: multiPartCRC32C,
				expectParts: []s3types.CompletedPart{
					{
						ChecksumCRC32C: aws.String(singlePartCRC32C),
						ETag:           aws.String(singlePartETag),
						PartNumber:     1,
					},
					{
						ChecksumCRC32C: aws.String(singlePartCRC32C),
						ETag:           aws.String(singlePartETag),
						PartNumber:     2,
					},
					{
						ChecksumCRC32C: aws.String(multiPartTailCRC32C),
						ETag:           aws.String(multiPartTailETag),
						PartNumber:     3,
					},
				},
				expectChecksumCRC32C: multiPartCRC32C,
				expectETag:           multiPartETag,
			},
			"SHA1": {
				algorithm:    s3types.ChecksumAlgorithmSha1,
				payload:      bytes.NewReader(multiPartBytes),
				checksumSHA1: multiPartSHA1,
				expectParts: []s3types.CompletedPart{
					{
						ChecksumSHA1: aws.String(singlePartSHA1),
						ETag:         aws.String(singlePartETag),
						PartNumber:   1,
					},
					{
						ChecksumSHA1: aws.String(singlePartSHA1),
						ETag:         aws.String(singlePartETag),
						PartNumber:   2,
					},
					{
						ChecksumSHA1: aws.String(multiPartTailSHA1),
						ETag:         aws.String(multiPartTailETag),
						PartNumber:   3,
					},
				},
				expectChecksumSHA1: multiPartSHA1,
				expectETag:         multiPartETag,
			},
			"SHA256": {
				algorithm:      s3types.ChecksumAlgorithmSha256,
				payload:        bytes.NewReader(multiPartBytes),
				checksumSHA256: multiPartSHA256,
				expectParts: []s3types.CompletedPart{
					{
						ChecksumSHA256: aws.String(singlePartSHA256),
						ETag:           aws.String(singlePartETag),
						PartNumber:     1,
					},
					{
						ChecksumSHA256: aws.String(singlePartSHA256),
						ETag:           aws.String(singlePartETag),
						PartNumber:     2,
					},
					{
						ChecksumSHA256: aws.String(multiPartTailSHA256),
						ETag:           aws.String(multiPartTailETag),
						PartNumber:     3,
					},
				},
				expectChecksumSHA256: multiPartSHA256,
				expectETag:           multiPartETag,
			},
		},
	}

	for group, cs := range cases {
		t.Run(group, func(t *testing.T) {
			for name, c := range cs {
				t.Run(name, func(t *testing.T) {
					mgr := manager.NewUploader(client)
					out, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
						Bucket:            bucketName,
						Key:               aws.String(t.Name()),
						Body:              c.payload,
						ChecksumAlgorithm: c.algorithm,
						ChecksumCRC32:     toStringPtr(c.checksumCRC32),
						ChecksumCRC32C:    toStringPtr(c.checksumCRC32C),
						ChecksumSHA1:      toStringPtr(c.checksumSHA1),
						ChecksumSHA256:    toStringPtr(c.checksumSHA256),
						ContentMD5:        toStringPtr(c.contentMD5),
					})
					if err != nil {
						t.Fatalf("expect no error, got %v", err)
					}

					if diff := cmp.Diff(c.expectParts, out.CompletedParts, cmpopts.IgnoreUnexported(types.CompletedPart{})); diff != "" {
						t.Errorf("expect parts match\n%s", diff)
					}

					if e, a := c.expectChecksumCRC32, aws.ToString(out.ChecksumCRC32); e != a {
						t.Errorf("expect %v CRC32 checksum, got %v", e, a)
					}
					if e, a := c.expectChecksumCRC32C, aws.ToString(out.ChecksumCRC32C); e != a {
						t.Errorf("expect %v CRC32C checksum, got %v", e, a)
					}
					if e, a := c.expectChecksumSHA1, aws.ToString(out.ChecksumSHA1); e != a {
						t.Errorf("expect %v SHA1 checksum, got %v", e, a)
					}
					if e, a := c.expectChecksumSHA256, aws.ToString(out.ChecksumSHA256); e != a {
						t.Errorf("expect %v SHA256 checksum, got %v", e, a)
					}
					if e, a := c.expectETag, aws.ToString(out.ETag); e != a {
						t.Errorf("expect %v ETag, got %v", e, a)
					}
				})
			}
		})
	}
}

func toStringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

type invalidateHash struct{}

func (b *invalidateHash) ID() string {
	return "s3manager:InvalidateHash"
}

func (b *invalidateHash) RegisterMiddleware(stack *middleware.Stack) error {
	return stack.Serialize.Add(b, middleware.After)
}

func (b *invalidateHash) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if input, ok := in.Parameters.(*s3.UploadPartInput); ok && input.PartNumber == 2 {
		ctx = v4.SetPayloadHash(ctx, "000")
	}

	return next.HandleSerialize(ctx, in)
}

func TestInteg_UploadFailCleanup(t *testing.T) {
	key := "12mb-leave"
	mgr := manager.NewUploader(client, func(u *manager.Uploader) {
		u.LeavePartsOnError = false
		u.ClientOptions = append(u.ClientOptions, func(options *s3.Options) {
			options.APIOptions = append(options.APIOptions, (&invalidateHash{}).RegisterMiddleware)
		})
	})
	_, err := mgr.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: bucketName,
		Key:    &key,
		Body:   bytes.NewReader(integBuf12MB),
	})
	if err == nil {
		t.Fatalf("expect error, but did not get one")
	}

	uploadID := ""
	var uf manager.MultiUploadFailure
	if !errors.As(err, &uf) {
		t.Errorf("")
	} else if uploadID = uf.UploadID(); len(uploadID) == 0 {
		t.Errorf("expect upload ID to not be empty, but was")
	}

	_, err = client.ListParts(context.Background(), &s3.ListPartsInput{
		Bucket: bucketName, Key: &key, UploadId: &uploadID,
	})
	if err == nil {
		t.Errorf("expect error for list parts, but got none")
	}
}
