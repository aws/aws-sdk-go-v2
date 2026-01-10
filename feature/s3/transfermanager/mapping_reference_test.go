package transfermanager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
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
		"SSEKMSKeyID",
		"ServerSideEncryption",
		"StorageClass",
		"Tagging",
		"WebsiteRedirectLocation"
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
        "SSEKMSKeyID",
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
	legacyMappings := map[string]string{
		"SSEKMSKeyId": "SSEKMSKeyID",
	}

	rtype := reflect.TypeFor[UploadObjectInput]()

	for _, field := range mappingReference.Definition.UploadRequest.PutObjectRequest {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("UploadInput: missing field %q", field)
		}
	}
}

func TestDefinition_UploadResponse(t *testing.T) {
	legacyMappings := map[string]string{
		"VersionId":   "VersionID",
		"SSEKMSKeyId": "SSEKMSKeyID",
	}

	rtype := reflect.TypeFor[UploadObjectOutput]()

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
	legacyMappings := map[string]string{
		"VersionId": "VersionID",
	}

	rtype := reflect.TypeFor[DownloadObjectInput]()

	for _, field := range mappingReference.Definition.DownloadRequest.GetObjectRequest {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("DownloadInput: missing field %q", field)
		}
	}
}

func TestDefinition_DownloadResponse(t *testing.T) {
	legacyMappings := map[string]string{
		"VersionId":   "VersionID",
		"SSEKMSKeyId": "SSEKMSKeyID",
	}

	rtype := reflect.TypeFor[DownloadObjectOutput]()

	for _, field := range mappingReference.Definition.DownloadResponse.GetObjectResponse {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("DownloadOutput: missing field %q", field)
		}
	}
}

func TestDefinition_GetRequest(t *testing.T) {
	legacyMappings := map[string]string{
		"VersionId": "VersionID",
	}

	rtype := reflect.TypeFor[GetObjectInput]()

	for _, field := range mappingReference.Definition.DownloadRequest.GetObjectRequest {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("GetInput: missing field %q", field)
		}
	}
}

func TestDefinition_GetResponse(t *testing.T) {
	legacyMappings := map[string]string{
		"VersionId":   "VersionID",
		"SSEKMSKeyId": "SSEKMSKeyID",
	}

	rtype := reflect.TypeFor[GetObjectOutput]()

	for _, field := range mappingReference.Definition.DownloadResponse.GetObjectResponse {
		if mapped, ok := legacyMappings[field]; ok {
			field = mapped
		}

		_, ok := rtype.FieldByName(field)
		if !ok {
			t.Errorf("DownloadOutput: missing field %q", field)
		}
	}
}

func TestConversion_UploadRequest_PutObjectRequest(t *testing.T) {
	src := UploadObjectInput{}

	now := time.Now()
	mockFields(&src, mappingReference.Conversion.UploadRequest.PutObjectRequest, map[string]interface{}{
		"Time": &now,
	})

	dst := src.mapSingleUploadInput(strings.NewReader(""), types.ChecksumAlgorithmCrc32)

	expectConvertedFields(t, dst, &src,
		mappingReference.Conversion.UploadRequest.PutObjectRequest,
		map[string]string{
			"SSEKMSKeyID": "SSEKMSKeyId",
		})
}

func TestConversion_UploadRequest_CreateMultipartUploadRequest(t *testing.T) {
	src := UploadObjectInput{}

	now := time.Now()
	mockFields(&src, mappingReference.Conversion.UploadRequest.CreateMultipartRequest, map[string]interface{}{
		"Time": &now,
	})

	dst := src.mapCreateMultipartUploadInput(types.ChecksumAlgorithmCrc32)

	expectConvertedFields(t, dst, &src,
		mappingReference.Conversion.UploadRequest.CreateMultipartRequest,
		map[string]string{
			"SSEKMSKeyID": "SSEKMSKeyId",
		})
}

func TestConversion_UploadRequest_UploadPartRequest(t *testing.T) {
	src := UploadObjectInput{}

	mockFields(&src, mappingReference.Conversion.UploadRequest.UploadPartRequest, map[string]interface{}{})

	dst := src.mapUploadPartInput(strings.NewReader(""), aws.Int32(1), aws.String(""), types.ChecksumAlgorithmCrc32)

	expectConvertedFields(t, dst, &src,
		mappingReference.Conversion.UploadRequest.UploadPartRequest,
		map[string]string{})
}

func TestConversion_UploadRequest_CompleteMultipartUploadRequest(t *testing.T) {
	src := UploadObjectInput{}

	mockFields(&src, mappingReference.Conversion.UploadRequest.CompleteMultipartRequest, map[string]interface{}{})

	dst := src.mapCompleteMultipartUploadInput(aws.String(""), []types.CompletedPart{})

	expectConvertedFields(t, dst, &src,
		mappingReference.Conversion.UploadRequest.CompleteMultipartRequest,
		map[string]string{})
}

func TestConversion_UploadRequest_AbortMultipartUploadRequest(t *testing.T) {
	src := UploadObjectInput{}

	mockFields(&src, mappingReference.Conversion.UploadRequest.AbortMultipartRequest, map[string]interface{}{})

	dst := src.mapAbortMultipartUploadInput(aws.String(""))

	expectConvertedFields(t, dst, &src,
		mappingReference.Conversion.UploadRequest.AbortMultipartRequest,
		map[string]string{})
}

func TestConversion_CompleteMultipartUploadResponse_UploadResponse(t *testing.T) {
	dst := UploadObjectOutput{}
	src := s3.CompleteMultipartUploadOutput{}

	mockFields(&src, mappingReference.Conversion.CompleteMultipartResponse.UploadResponse, map[string]interface{}{})

	dst.mapFromCompleteMultipartUploadOutput(&src, aws.String(""), aws.String(""), 0, []types.CompletedPart{})

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.CompleteMultipartResponse.UploadResponse,
		map[string]string{
			"VersionId":   "VersionID",
			"SSEKMSKeyId": "SSEKMSKeyID",
		})
}

func TestConversion_PutObjectResponse_UploadResponse(t *testing.T) {
	dst := UploadObjectOutput{}
	src := s3.PutObjectOutput{}

	mockFields(&src, mappingReference.Conversion.PutObjectResponse.UploadResponse, map[string]interface{}{})

	dst.mapFromPutObjectOutput(&src, aws.String(""), aws.String(""), 0)

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.PutObjectResponse.UploadResponse,
		map[string]string{
			"VersionId":   "VersionID",
			"SSEKMSKeyId": "SSEKMSKeyID",
		})
}

func TestConversion_GetObjectResponse_DownloadResponse(t *testing.T) {
	dst := DownloadObjectOutput{}
	src := s3.GetObjectOutput{}

	now := time.Now()
	mockFields(&src, mappingReference.Conversion.GetObjectResponse.DownloadResponse, map[string]interface{}{
		"Time": &now,
	})

	dst.mapFromGetObjectOutput(&src, "")

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.GetObjectResponse.DownloadResponse,
		map[string]string{
			"VersionId":   "VersionID",
			"SSEKMSKeyId": "SSEKMSKeyID",
		})
}

func TestConversion_GetObjectResponse_GetResponse(t *testing.T) {
	dst := GetObjectOutput{}
	src := s3.GetObjectOutput{}

	now := time.Now()
	mockFields(&src, mappingReference.Conversion.GetObjectResponse.DownloadResponse, map[string]interface{}{
		"Time": &now,
	})

	dst.mapFromGetObjectOutput(&src, "")

	expectConvertedFields(t, &dst, &src,
		mappingReference.Conversion.GetObjectResponse.DownloadResponse,
		map[string]string{
			"VersionId":   "VersionID",
			"SSEKMSKeyId": "SSEKMSKeyID",
		})
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
		if dstf.Type() != srcf.Type() {
			// s3 and tm use their own string type, which is always deeply unequal
			// so their underlying string value is compared directly
			expect = srcf.String()
			actual = dstf.String()
		}
		if !reflect.DeepEqual(expect, actual) {
			t.Errorf("field %q: %v != %v", srcField, expect, actual)
		}
	}
}

// v must be a pointer to a struct
func mockFields(v any, fields []string, legacyTypes map[string]interface{}) {
	rv := reflect.ValueOf(v).Elem()
	for _, field := range fields {
		mockValue(rv.FieldByName(field), field, legacyTypes)
	}
}

func mockValue(v reflect.Value, field string, legacyTypes map[string]interface{}) {
	switch v.Kind() {
	case reflect.Pointer:
		switch v.Type().Elem().Kind() {
		case reflect.Bool:
			vv := true
			v.Set(reflect.ValueOf(&vv))
		case reflect.String:
			v.Set(reflect.ValueOf(&field))
		case reflect.Int64:
			vv := int64(1)
			v.Set(reflect.ValueOf(&vv))
		case reflect.Int32:
			vv := int32(2)
			v.Set(reflect.ValueOf(&vv))
		case reflect.Struct:
			vv, ok := legacyTypes[v.Type().Elem().Name()]
			if !ok {
				panic(fmt.Sprintf("need to handle %v for field %s", v.Type().Elem().Kind(), field))
			}
			v.Set(reflect.ValueOf(vv))
		default:
			panic(fmt.Sprintf("need to handle %v", v.Type().Elem().Kind()))
		}
	case reflect.String:
		// Convert() because it's probably a string enum
		v.Set(reflect.ValueOf(field).Convert(v.Type()))
	case reflect.Map:
		v.Set(reflect.ValueOf(map[string]string{"a": "b"}))
	default:
		panic(fmt.Sprintf("need to handle %v", v.Type().Elem().Kind()))
	}
}
