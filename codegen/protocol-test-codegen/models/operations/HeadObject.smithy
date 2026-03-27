$version: "2"

namespace com.amazonaws.sdk.benchmark

@documentation("""
    From Amazon S3.
    HeadObject deserializes a large set of headers.
""")
@http(method: "HEAD", uri: "/{Bucket}/{Key+}", code: 200)
operation HeadObject {
    input: HeadObjectRequest
    output: HeadObjectOutput
}

structure HeadObjectRequest {
    @httpLabel
    @required
    Bucket: String
    @httpHeader("If-Match")
    IfMatch: String
    @httpHeader("If-Modified-Since")
    IfModifiedSince: Timestamp
    @httpHeader("If-None-Match")
    IfNoneMatch: String
    @httpHeader("If-Unmodified-Since")
    IfUnmodifiedSince: Timestamp
    @httpLabel
    @required
    Key: String
    @httpHeader("Range")
    Range: String
    @httpQuery("response-cache-control")
    ResponseCacheControl: String
    @httpQuery("response-content-disposition")
    ResponseContentDisposition: String
    @httpQuery("response-content-encoding")
    ResponseContentEncoding: String
    @httpQuery("response-content-language")
    ResponseContentLanguage: String
    @httpQuery("response-content-type")
    ResponseContentType: String
    @httpQuery("response-expires")
    ResponseExpires: Timestamp
    @httpQuery("versionId")
    VersionId: String
    @httpHeader("x-amz-server-side-encryption-customer-algorithm")
    SSECustomerAlgorithm: String
    @httpHeader("x-amz-server-side-encryption-customer-key")
    SSECustomerKey: String
    @httpHeader("x-amz-server-side-encryption-customer-key-MD5")
    SSECustomerKeyMD5: String
    @httpHeader("x-amz-request-payer")
    RequestPayer: String
    @httpQuery("partNumber")
    PartNumber: Integer
    @httpHeader("x-amz-expected-bucket-owner")
    ExpectedBucketOwner: String
    @httpHeader("x-amz-checksum-mode")
    ChecksumMode: String
}

structure HeadObjectOutput {
    @httpHeader("x-amz-delete-marker")
    DeleteMarker: Boolean
    @httpHeader("accept-ranges")
    AcceptRanges: String
    @httpHeader("x-amz-expiration")
    Expiration: String
    @httpHeader("x-amz-restore")
    Restore: String
    @httpHeader("x-amz-archive-status")
    ArchiveStatus: String
    @httpHeader("Last-Modified")
    LastModified: Timestamp
    @httpHeader("Content-Length")
    ContentLength: Long
    @httpHeader("x-amz-checksum-crc32")
    ChecksumCRC32: String
    @httpHeader("x-amz-checksum-crc32c")
    ChecksumCRC32C: String
    @httpHeader("x-amz-checksum-crc64nvme")
    ChecksumCRC64NVME: String
    @httpHeader("x-amz-checksum-sha1")
    ChecksumSHA1: String
    @httpHeader("x-amz-checksum-sha256")
    ChecksumSHA256: String
    @httpHeader("x-amz-checksum-type")
    ChecksumType: String
    @httpHeader("ETag")
    ETag: String
    @httpHeader("x-amz-missing-meta")
    MissingMeta: Integer
    @httpHeader("x-amz-version-id")
    VersionId: String
    @httpHeader("Cache-Control")
    CacheControl: String
    @httpHeader("Content-Disposition")
    ContentDisposition: String
    @httpHeader("Content-Encoding")
    ContentEncoding: String
    @httpHeader("Content-Language")
    ContentLanguage: String
    @httpHeader("Content-Type")
    ContentType: String
    @httpHeader("Content-Range")
    ContentRange: String
    @httpHeader("Expires")
    Expires: Timestamp
    @httpHeader("x-amz-website-redirect-location")
    WebsiteRedirectLocation: String
    @httpHeader("x-amz-server-side-encryption")
    ServerSideEncryption: String
    @httpHeader("x-amz-server-side-encryption-customer-algorithm")
    SSECustomerAlgorithm: String
    @httpHeader("x-amz-server-side-encryption-customer-key-MD5")
    SSECustomerKeyMD5: String
    @httpHeader("x-amz-server-side-encryption-aws-kms-key-id")
    SSEKMSKeyId: String
    @httpHeader("x-amz-server-side-encryption-bucket-key-enabled")
    BucketKeyEnabled: Boolean
    @httpHeader("x-amz-storage-class")
    StorageClass: String
    @httpHeader("x-amz-request-charged")
    RequestCharged: String
    @httpHeader("x-amz-replication-status")
    ReplicationStatus: String
    @httpHeader("x-amz-mp-parts-count")
    PartsCount: Integer
    @httpHeader("x-amz-tagging-count")
    TagCount: Integer
    @httpHeader("x-amz-object-lock-mode")
    ObjectLockMode: String
    @httpHeader("x-amz-object-lock-retain-until-date")
    ObjectLockRetainUntilDate: Timestamp
    @httpHeader("x-amz-object-lock-legal-hold")
    ObjectLockLegalHoldStatus: String
}