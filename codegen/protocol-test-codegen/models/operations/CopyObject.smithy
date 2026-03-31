$version: "2"

namespace com.amazonaws.sdk.benchmark

use aws.protocols#restXml
use aws.protocols#restJson1
use smithy.test#httpRequestTests
use smithy.test#httpResponseTests

@documentation("""
    From Amazon S3.
    CopyObject serializes a large set of headers.
""")
@http(
    method: "PUT",
    uri: "/{Bucket}/{Key+}?x-id=CopyObject",
    code: 200
)
@httpRequestTests([
    {
        id: "restXml_CopyObjectRequest_Baseline"
        documentation: """
        """
        protocol: restXml
        method: "PUT"
        uri: "/test-bucket/test-key?x-id=CopyObject"
        params: {
            Bucket: "test-bucket"
            Key: "test-key"
            CopySource: "/source-bucket/source-key"
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "restXml_CopyObjectRequest_M"
        documentation: """
        Serialization of a large set of headers.
        """
        protocol: restXml
        method: "PUT"
        uri: "/dest-bucket/dest-key?x-id=CopyObject"
        params: {
            Bucket: "dest-bucket"
            Key: "dest-key"
            ACL: "private"
            CacheControl: "max-age=3600"
            ChecksumAlgorithm: "SHA256"
            ContentDisposition: "attachment; filename=\"example.txt\""
            ContentEncoding: "gzip"
            ContentLanguage: "en-US"
            ContentType: "text/plain"
            CopySource: "/source-bucket/source-key"
            CopySourceIfMatch: "\"9bb58f26192e4ba00f01e2e7b136bbd8\""
            CopySourceIfModifiedSince: 1609459200
            CopySourceIfNoneMatch: "\"different-etag\""
            CopySourceIfUnmodifiedSince: 1640995199
            Expires: 1641024000
            GrantFullControl: "id=canonical-user-id"
            GrantRead: "id=read-user-id"
            GrantReadACP: "id=read-acp-user-id"
            GrantWriteACP: "id=write-acp-user-id"
            IfMatch: "\"target-etag\""
            IfNoneMatch: "\"different-target-etag\""
            MetadataDirective: "REPLACE"
            TaggingDirective: "REPLACE"
            ServerSideEncryption: "AES256"
            StorageClass: "STANDARD_IA"
            WebsiteRedirectLocation: "https://example.com/redirect"
            SSECustomerAlgorithm: "AES256"
            SSECustomerKey: "customer-key-base64"
            SSECustomerKeyMD5: "customer-key-md5-hash"
            SSEKMSKeyId: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
            SSEKMSEncryptionContext: "encryption-context"
            BucketKeyEnabled: true
            CopySourceSSECustomerAlgorithm: "AES256"
            CopySourceSSECustomerKey: "source-customer-key-base64"
            CopySourceSSECustomerKeyMD5: "source-customer-key-md5-hash"
            RequestPayer: "BucketOwner"
            Tagging: "key1=value1&key2=value2"
            ObjectLockMode: "GOVERNANCE"
            ObjectLockRetainUntilDate: 1641024000
            ObjectLockLegalHoldStatus: "ON"
            ExpectedBucketOwner: "123456789012"
            ExpectedSourceBucketOwner: "123456789012"
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "restJson1_CopyObjectRequest_Baseline"
        documentation: """
        """
        protocol: restJson1
        method: "PUT"
        uri: "/test-bucket/test-key?x-id=CopyObject"
        params: {
            Bucket: "test-bucket"
            Key: "test-key"
            CopySource: "/source-bucket/source-key"
        }
        tags: ["serde-benchmark"]
    }
    {
        id: "restJson1_CopyObjectRequest_M"
        documentation: """
        Serialization of a large set of headers.
        """
        protocol: restJson1
        method: "PUT"
        uri: "/dest-bucket/dest-key?x-id=CopyObject"
        params: {
            Bucket: "dest-bucket"
            Key: "dest-key"
            ACL: "private"
            CacheControl: "max-age=3600"
            ChecksumAlgorithm: "SHA256"
            ContentDisposition: "attachment; filename=\"example.txt\""
            ContentEncoding: "gzip"
            ContentLanguage: "en-US"
            ContentType: "text/plain"
            CopySource: "/source-bucket/source-key"
            CopySourceIfMatch: "\"9bb58f26192e4ba00f01e2e7b136bbd8\""
            CopySourceIfModifiedSince: 1609459200
            CopySourceIfNoneMatch: "\"different-etag\""
            CopySourceIfUnmodifiedSince: 1640995199
            Expires: 1641024000
            GrantFullControl: "id=canonical-user-id"
            GrantRead: "id=read-user-id"
            GrantReadACP: "id=read-acp-user-id"
            GrantWriteACP: "id=write-acp-user-id"
            IfMatch: "\"target-etag\""
            IfNoneMatch: "\"different-target-etag\""
            MetadataDirective: "REPLACE"
            TaggingDirective: "REPLACE"
            ServerSideEncryption: "AES256"
            StorageClass: "STANDARD_IA"
            WebsiteRedirectLocation: "https://example.com/redirect"
            SSECustomerAlgorithm: "AES256"
            SSECustomerKey: "customer-key-base64"
            SSECustomerKeyMD5: "customer-key-md5-hash"
            SSEKMSKeyId: "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
            SSEKMSEncryptionContext: "encryption-context"
            BucketKeyEnabled: true
            CopySourceSSECustomerAlgorithm: "AES256"
            CopySourceSSECustomerKey: "source-customer-key-base64"
            CopySourceSSECustomerKeyMD5: "source-customer-key-md5-hash"
            RequestPayer: "BucketOwner"
            Tagging: "key1=value1&key2=value2"
            ObjectLockMode: "GOVERNANCE"
            ObjectLockRetainUntilDate: 1641024000
            ObjectLockLegalHoldStatus: "ON"
            ExpectedBucketOwner: "123456789012"
            ExpectedSourceBucketOwner: "123456789012"
        }
        tags: ["serde-benchmark"]
    }
])
@httpResponseTests([
    {
        id: "restXml_CopyObjectOutput_Baseline"
        documentation: """
        """
        protocol: restXml
        code: 200
        tags: ["serde-benchmark"]
    }
    {
        id: "restXml_CopyObjectOutput_M"
        documentation: """
        Deserialization of headers and XML body.
        """
        protocol: restXml
        code: 200
        headers: {
            "x-amz-expiration": "expiry-date=\"Fri, 01 Jan 2022 00:00:00 GMT\", rule-id=\"rule1\""
            "x-amz-copy-source-version-id": "source-version-id-12345"
            "x-amz-version-id": "dest-version-id-67890"
            "x-amz-server-side-encryption": "AES256"
            "x-amz-server-side-encryption-customer-algorithm": "AES256"
            "x-amz-server-side-encryption-customer-key-MD5": "customer-key-md5-hash"
            "x-amz-server-side-encryption-aws-kms-key-id": "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
            "x-amz-server-side-encryption-context": "encryption-context"
            "x-amz-server-side-encryption-bucket-key-enabled": "true"
            "x-amz-request-charged": "requester"
        }
        body: """
<CopyObjectResult>
    <ETag>"9bb58f26192e4ba00f01e2e7b136bbd8"</ETag>
    <LastModified>2021-01-01T00:00:00.000Z</LastModified>
    <ChecksumType>SHA256</ChecksumType>
    <ChecksumCRC32>checksum-crc32</ChecksumCRC32>
    <ChecksumCRC32C>checksum-crc32c</ChecksumCRC32C>
    <ChecksumCRC64NVME>checksum-crc64nvme</ChecksumCRC64NVME>    
    <ChecksumSHA1>checksum-sha1</ChecksumSHA1>
    <ChecksumSHA256>checksum-sha256</ChecksumSHA256>
</CopyObjectResult>
        """
        tags: ["serde-benchmark"]
    }
    {
        id: "restJson1_CopyObjectOutput_Baseline"
        documentation: """
        """
        protocol: restJson1
        code: 200
        tags: ["serde-benchmark"]
    }
    {
        id: "restJson1_CopyObjectOutput_M"
        documentation: """
        Deserialization of headers and JSON body.
        """
        protocol: restJson1
        code: 200
        headers: {
            "x-amz-expiration": "expiry-date=\"Fri, 01 Jan 2022 00:00:00 GMT\", rule-id=\"rule1\""
            "x-amz-copy-source-version-id": "source-version-id-12345"
            "x-amz-version-id": "dest-version-id-67890"
            "x-amz-server-side-encryption": "AES256"
            "x-amz-server-side-encryption-customer-algorithm": "AES256"
            "x-amz-server-side-encryption-customer-key-MD5": "customer-key-md5-hash"
            "x-amz-server-side-encryption-aws-kms-key-id": "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
            "x-amz-server-side-encryption-context": "encryption-context"
            "x-amz-server-side-encryption-bucket-key-enabled": "true"
            "x-amz-request-charged": "requester"
        }
        body: """
{
    "ETag": "\\"9bb58f26192e4ba00f01e2e7b136bbd8\\"",
    "LastModified": 1609459200,
    "ChecksumType": "SHA256",
    "ChecksumCRC32": "checksum-crc32",
    "ChecksumCRC32C": "checksum-crc32c",
    "ChecksumCRC64NVME": "checksum-crc64nvme",
    "ChecksumSHA1": "checksum-sha1",
    "ChecksumSHA256": "checksum-sha256"
}
        """
        tags: ["serde-benchmark"]
    }
])
operation CopyObject {
    input: CopyObjectRequest
    output: CopyObjectOutput
}

structure CopyObjectRequest {
    @httpHeader("x-amz-acl")
    ACL: String
    @httpLabel
    @required
    Bucket: String
    @httpHeader("Cache-Control")
    CacheControl: String
    @httpHeader("x-amz-checksum-algorithm")
    ChecksumAlgorithm: String
    @httpHeader("Content-Disposition")
    ContentDisposition: String
    @httpHeader("Content-Encoding")
    ContentEncoding: String
    @httpHeader("Content-Language")
    ContentLanguage: String
    @httpHeader("Content-Type")
    ContentType: String
    @httpHeader("x-amz-copy-source")
    @required
    CopySource: String
    @httpHeader("x-amz-copy-source-if-match")
    CopySourceIfMatch: String
    @httpHeader("x-amz-copy-source-if-modified-since")
    CopySourceIfModifiedSince: Timestamp
    @httpHeader("x-amz-copy-source-if-none-match")
    CopySourceIfNoneMatch: String
    @httpHeader("x-amz-copy-source-if-unmodified-since")
    CopySourceIfUnmodifiedSince: Timestamp
    @httpHeader("Expires")
    Expires: Timestamp
    @httpHeader("x-amz-grant-full-control")
    GrantFullControl: String
    @httpHeader("x-amz-grant-read")
    GrantRead: String
    @httpHeader("x-amz-grant-read-acp")
    GrantReadACP: String
    @httpHeader("x-amz-grant-write-acp")
    GrantWriteACP: String
    @httpHeader("If-Match")
    IfMatch: String
    @httpHeader("If-None-Match")
    IfNoneMatch: String
    @httpLabel
    @required
    Key: String
    @httpHeader("x-amz-metadata-directive")
    MetadataDirective: String
    @httpHeader("x-amz-tagging-directive")
    TaggingDirective: String
    @httpHeader("x-amz-server-side-encryption")
    ServerSideEncryption: String
    @httpHeader("x-amz-storage-class")
    StorageClass: String
    @httpHeader("x-amz-website-redirect-location")
    WebsiteRedirectLocation: String
    @httpHeader("x-amz-server-side-encryption-customer-algorithm")
    SSECustomerAlgorithm: String
    @httpHeader("x-amz-server-side-encryption-customer-key")
    SSECustomerKey: String
    @httpHeader("x-amz-server-side-encryption-customer-key-MD5")
    SSECustomerKeyMD5: String
    @httpHeader("x-amz-server-side-encryption-aws-kms-key-id")
    SSEKMSKeyId: String
    @httpHeader("x-amz-server-side-encryption-context")
    SSEKMSEncryptionContext: String
    @httpHeader("x-amz-server-side-encryption-bucket-key-enabled")
    BucketKeyEnabled: Boolean
    @httpHeader("x-amz-copy-source-server-side-encryption-customer-algorithm")
    CopySourceSSECustomerAlgorithm: String
    @httpHeader("x-amz-copy-source-server-side-encryption-customer-key")
    CopySourceSSECustomerKey: String
    @httpHeader("x-amz-copy-source-server-side-encryption-customer-key-MD5")
    CopySourceSSECustomerKeyMD5: String
    @httpHeader("x-amz-request-payer")
    RequestPayer: String
    @httpHeader("x-amz-tagging")
    Tagging: String
    @httpHeader("x-amz-object-lock-mode")
    ObjectLockMode: String
    @httpHeader("x-amz-object-lock-retain-until-date")
    ObjectLockRetainUntilDate: Timestamp
    @httpHeader("x-amz-object-lock-legal-hold")
    ObjectLockLegalHoldStatus: String
    @httpHeader("x-amz-expected-bucket-owner")
    ExpectedBucketOwner: String
    @httpHeader("x-amz-source-expected-bucket-owner")
    ExpectedSourceBucketOwner: String
}

structure CopyObjectOutput {
    @httpHeader("x-amz-expiration")
    Expiration: String
    @httpHeader("x-amz-copy-source-version-id")
    CopySourceVersionId: String
    @httpHeader("x-amz-version-id")
    VersionId: String
    @httpHeader("x-amz-server-side-encryption")
    ServerSideEncryption: String
    @httpHeader("x-amz-server-side-encryption-customer-algorithm")
    SSECustomerAlgorithm: String
    @httpHeader("x-amz-server-side-encryption-customer-key-MD5")
    SSECustomerKeyMD5: String
    @httpHeader("x-amz-server-side-encryption-aws-kms-key-id")
    SSEKMSKeyId: String
    @httpHeader("x-amz-server-side-encryption-context")
    SSEKMSEncryptionContext: String
    @httpHeader("x-amz-server-side-encryption-bucket-key-enabled")
    BucketKeyEnabled: Boolean
    @httpHeader("x-amz-request-charged")
    RequestCharged: String
    ETag: String
    LastModified: Timestamp
    ChecksumType: String
    ChecksumCRC32: String
    ChecksumCRC32C: String
    ChecksumCRC64NVME: String
    ChecksumSHA1: String
    ChecksumSHA256: String
}
