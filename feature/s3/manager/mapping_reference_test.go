package manager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// all tests use a single package-level instance of this spec
//
// legacy name mappings may mutate some of the field names but it doesn't break
// anything from test to test right now
var mappingReference struct {
	Definition struct {
		UploadRequest struct {
			PutObjectRequest []string
		}
		UploadResponse struct {
			PutObjectResponse []string
		}
		DownloadRequest struct {
			GetObjectRequest []string
		}
		DownloadResponse struct {
			GetObjectResponse []string
		}
	}
	Conversion struct {
		UploadRequest struct {
			PutObjectRequest         []string
			CreateMultipartRequest   []string
			UploadPartRequest        []string
			CompleteMultipartRequest []string
			AbortMultipartRequest    []string
		}
		CompleteMultipartResponse struct {
			UploadResponse []string
		}
		PutObjectResponse struct {
			UploadResponse []string
		}
		GetObjectResponse struct {
			DownloadResponse []string
		}
	}
}

const mappingReferenceJSON = `
{
  "Definition": {
    "UploadRequest": {
      "PutObjectRequest": [
        "ACL",
        "Bucket",
        "BucketKeyEnabled",
        "CacheControl",
        "ChecksumAlgorithm",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ContentDisposition",
        "ContentEncoding",
        "ContentLanguage",
        "ContentType",
        "ExpectedBucketOwner",
        "Expires",
        "GrantFullControl",
        "GrantRead",
        "GrantReadACP",
        "GrantWriteACP",
        "IfMatch",
        "IfNoneMatch",
        "Key",
        "Metadata",
        "ObjectLockLegalHoldStatus",
        "ObjectLockMode",
        "ObjectLockRetainUntilDate",
        "RequestPayer",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5",
        "SSEKMSEncryptionContext",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "StorageClass",
        "Tagging",
        "WebsiteRedirectLocation"
      ]
    },
    "UploadResponse": {
      "PutObjectResponse": [
        "BucketKeyEnabled",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ChecksumType",
        "ETag",
        "Expiration",
        "RequestCharged",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "VersionId"
      ]
    },
    "DownloadRequest": {
      "GetObjectRequest": [
        "Bucket",
        "ChecksumMode",
        "ExpectedBucketOwner",
        "IfMatch",
        "IfModifiedSince",
        "IfNoneMatch",
        "IfUnmodifiedSince",
        "Key",
        "RequestPayer",
        "ResponseCacheControl",
        "ResponseContentDisposition",
        "ResponseContentEncoding",
        "ResponseContentLanguage",
        "ResponseContentType",
        "ResponseExpires",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5",
        "VersionId"
      ]
    },
    "DownloadResponse": {
      "GetObjectResponse": [
        "AcceptRanges",
        "BucketKeyEnabled",
        "CacheControl",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ChecksumType",
        "ContentDisposition",
        "ContentEncoding",
        "ContentLanguage",
        "ContentLength",
        "ContentRange",
        "ContentType",
        "DeleteMarker",
        "ETag",
        "Expiration",
        "Expires",
        "LastModified",
        "Metadata",
        "MissingMeta",
        "ObjectLockLegalHoldStatus",
        "ObjectLockMode",
        "ObjectLockRetainUntilDate",
        "PartsCount",
        "ReplicationStatus",
        "RequestCharged",
        "Restore",
        "SSECustomerAlgorithm",
        "SSECustomerKeyMD5",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "StorageClass",
        "TagCount",
        "VersionId",
        "WebsiteRedirectLocation"
      ]
    }
  },
  "Conversion": {
    "UploadRequest": {
      "PutObjectRequest": [
        "Bucket",
        "ChecksumAlgorithm",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ExpectedBucketOwner",
        "Key",
        "RequestPayer",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5"
      ],
      "CreateMultipartRequest": [
        "ACL",
        "Bucket",
        "BucketKeyEnabled",
        "CacheControl",
        "ChecksumAlgorithm",
        "ContentDisposition",
        "ContentEncoding",
        "ContentLanguage",
        "ContentType",
        "ExpectedBucketOwner",
        "Expires",
        "GrantFullControl",
        "GrantRead",
        "GrantReadACP",
        "GrantWriteACP",
        "Key",
        "Metadata",
        "ObjectLockLegalHoldStatus",
        "ObjectLockMode",
        "ObjectLockRetainUntilDate",
        "RequestPayer",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5",
        "SSEKMSEncryptionContext",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "StorageClass",
        "Tagging",
        "WebsiteRedirectLocation"
      ],
      "UploadPartRequest": [
        "Bucket",
        "ChecksumAlgorithm",
        "ExpectedBucketOwner",
        "Key",
        "RequestPayer",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5"
      ],
      "CompleteMultipartRequest": [
        "Bucket",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ExpectedBucketOwner",
        "IfMatch",
        "IfNoneMatch",
        "Key",
        "RequestPayer",
        "SSECustomerAlgorithm",
        "SSECustomerKey",
        "SSECustomerKeyMD5"
      ],
      "AbortMultipartRequest": [
        "Bucket",
        "ExpectedBucketOwner",
        "Key",
        "RequestPayer"
      ]
    },
    "CompleteMultipartResponse": {
      "UploadResponse": [
        "BucketKeyEnabled",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ChecksumType",
        "ETag",
        "Expiration",
        "RequestCharged",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "VersionId"
      ]
    },
    "PutObjectResponse": {
      "UploadResponse": [
        "BucketKeyEnabled",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ChecksumType",
        "ETag",
        "Expiration",
        "RequestCharged",
        "SSECustomerAlgorithm",
        "SSECustomerKeyMD5",
        "SSEKMSEncryptionContext",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "VersionId"
      ]
    },
    "GetObjectResponse": {
      "DownloadResponse": [
        "AcceptRanges",
        "BucketKeyEnabled",
        "CacheControl",
        "ChecksumCRC32",
        "ChecksumCRC32C",
        "ChecksumCRC64NVME",
        "ChecksumSHA1",
        "ChecksumSHA256",
        "ChecksumType",
        "ContentDisposition",
        "ContentEncoding",
        "ContentLanguage",
        "ContentLength",
        "ContentRange",
        "ContentType",
        "DeleteMarker",
        "ETag",
        "Expiration",
        "Expires",
        "ExpiresString",
        "LastModified",
        "Metadata",
        "MissingMeta",
        "ObjectLockLegalHoldStatus",
        "ObjectLockMode",
        "ObjectLockRetainUntilDate",
        "PartsCount",
        "ReplicationStatus",
        "RequestCharged",
        "Restore",
        "SSECustomerAlgorithm",
        "SSECustomerKeyMD5",
        "SSEKMSKeyId",
        "ServerSideEncryption",
        "StorageClass",
        "TagCount",
        "VersionId",
        "WebsiteRedirectLocation"
      ]
    }
  }
}
`

func init() {
	if err := json.Unmarshal([]byte(mappingReferenceJSON), &mappingReference); err != nil {
		panic(err)
	}
}

func TestDefinition_UploadRequest(t *testing.T) {
	t.Skip("[non-SEP-compliant] Upload uses s3.PutObjectInput directly")
}

func TestDefinition_UploadResponse(t *testing.T) {
	legacyMappings := map[string]string{
		"VersionId": "VersionID",
	}

	rtype := reflect.TypeOf(UploadOutput{})

	for _, field := range mappingReference.Definition.UploadResponse.PutObjectResponse {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("UploadOutput: missing field %q", field)
		}
	}
}

func TestDefinition_DownloadRequest(t *testing.T) {
	t.Skip("[non-SEP-compliant] Download uses s3.GetObjectInput directly")
}

func TestDefinition_DownloadResponse(t *testing.T) {
	t.Skip("[non-SEP-compliant] Download does not have structured output")
}

func TestConversion_UploadRequest_All(t *testing.T) {
	t.Skip("[non-SEP-compliant] Upload uses s3.PutObjectInput directly")
}

func TestConversion_CompleteMultipartUploadResponse_UploadResponse(t *testing.T) {
	dst := UploadOutput{}
	src := s3.CompleteMultipartUploadOutput{}

	mockFields(&src, mappingReference.Conversion.CompleteMultipartResponse.UploadResponse)

	convertCompleteMultipartUploadResponse(&dst, &src)

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.CompleteMultipartResponse.UploadResponse,
		map[string]string{
			"VersionId": "VersionID",
		})
}

func TestConversion_PutObjectResponse_UploadResponse(t *testing.T) {
	t.Skip("mapping reference appears to be invalid - it specifies 3 SSE-related fields to copy that it does not define on UploadResponse. need clarification from SEP authors")

	dst := UploadOutput{}
	src := s3.PutObjectOutput{}

	mockFields(&src, mappingReference.Conversion.PutObjectResponse.UploadResponse)

	convertPutObjectResponse(&dst, &src)

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.PutObjectResponse.UploadResponse,
		map[string]string{
			"VersionId": "VersionID",
		})
}

func TestConversion_GetObjectResponse_DownloadResponse(t *testing.T) {
	t.Skip("[non-SEP-compliant] Download does not have structured output")
}

func expectConvertedFields(t *testing.T, dst, src any, fields []string, legacyNames map[string]string) {
	t.Helper()

	dstv := reflect.ValueOf(dst).Elem()
	srcv := reflect.ValueOf(src).Elem()
	for _, srcField := range fields {
		dstField := srcField
		if legacy, ok := legacyNames[srcField]; ok {
			dstField = legacy
		}

		// indirect for any fields that happen to different in pointerness
		dstf := reflect.Indirect(dstv.FieldByName(dstField))
		srcf := reflect.Indirect(srcv.FieldByName(srcField))
		if !dstf.IsValid() {
			t.Fatalf("dst is missing field %q, do you need a legacy name mapping?", srcField)
		}

		expect := srcf.Interface()
		actual := dstf.Interface()
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("field %q: %v != %v", srcField, expect, actual)
		}
	}
}

// v must be a pointer to a struct
func mockFields(v any, fields []string) {
	rv := reflect.ValueOf(v).Elem()
	for _, field := range fields {
		mockValue(rv.FieldByName(field), field)
	}
}

func mockValue(v reflect.Value, field string) {
	switch v.Kind() {
	case reflect.Pointer:
		switch v.Type().Elem().Kind() {
		case reflect.Bool:
			vv := true
			v.Set(reflect.ValueOf(&vv))
		case reflect.String:
			v.Set(reflect.ValueOf(&field))
		default:
			panic(fmt.Sprintf("need to handle %v", v.Type().Elem().Kind()))
		}
	case reflect.String:
		// Convert() because it's probably a string enum
		v.Set(reflect.ValueOf(field).Convert(v.Type()))
	default:
		panic(fmt.Sprintf("need to handle %v", v.Type().Elem().Kind()))
	}
}
