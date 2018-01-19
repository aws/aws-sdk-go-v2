package s3crypto_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3crypto"
)

func TestHeaderV2SaveStrategy(t *testing.T) {
	cases := []struct {
		env      s3crypto.Envelope
		expected map[string]string
	}{
		{
			s3crypto.Envelope{
				CipherKey:             "Foo",
				IV:                    "Bar",
				MatDesc:               "{}",
				WrapAlg:               s3crypto.KMSWrap,
				CEKAlg:                s3crypto.AESGCMNoPadding,
				TagLen:                "128",
				UnencryptedMD5:        "hello",
				UnencryptedContentLen: "0",
			},
			map[string]string{
				"X-Amz-Key-V2":                     "Foo",
				"X-Amz-Iv":                         "Bar",
				"X-Amz-Matdesc":                    "{}",
				"X-Amz-Wrap-Alg":                   s3crypto.KMSWrap,
				"X-Amz-Cek-Alg":                    s3crypto.AESGCMNoPadding,
				"X-Amz-Tag-Len":                    "128",
				"X-Amz-Unencrypted-Content-Md5":    "hello",
				"X-Amz-Unencrypted-Content-Length": "0",
			},
		},
		{
			s3crypto.Envelope{
				CipherKey:             "Foo",
				IV:                    "Bar",
				MatDesc:               "{}",
				WrapAlg:               s3crypto.KMSWrap,
				CEKAlg:                s3crypto.AESGCMNoPadding,
				UnencryptedMD5:        "hello",
				UnencryptedContentLen: "0",
			},
			map[string]string{
				"X-Amz-Key-V2":                     "Foo",
				"X-Amz-Iv":                         "Bar",
				"X-Amz-Matdesc":                    "{}",
				"X-Amz-Wrap-Alg":                   s3crypto.KMSWrap,
				"X-Amz-Cek-Alg":                    s3crypto.AESGCMNoPadding,
				"X-Amz-Unencrypted-Content-Md5":    "hello",
				"X-Amz-Unencrypted-Content-Length": "0",
			},
		},
	}

	for _, c := range cases {
		params := &s3.PutObjectInput{}
		req := &aws.Request{
			Params: params,
		}
		strat := s3crypto.HeaderV2SaveStrategy{}
		err := strat.Save(c.env, req)
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		if !reflect.DeepEqual(c.expected, params.Metadata) {
			t.Errorf("expected %v, but received %v", c.expected, params.Metadata)
		}
	}
}
