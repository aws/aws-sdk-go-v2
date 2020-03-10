Breaking Change
---

Services
---

SDK Features
---
* `aws`: Add Support for additional credential providers and credential configuration chaining ([#488](https://github.com/aws/aws-sdk-go-v2/pull/488))
  * `aws/processcreds`: Adds Support for the Process Credential Provider
    * Fixes [#249](https://github.com/aws/aws-sdk-go-v2/issues/249)
  * `aws/stscreds`: Adds Support for the Web Identity Credential Provider
    * Fixes [#475](https://github.com/aws/aws-sdk-go-v2/issues/475)
    * Fixes [#338](https://github.com/aws/aws-sdk-go-v2/issues/338)
  * Adds Support for `credential_source`
    * Fixes [#274](https://github.com/aws/aws-sdk-go-v2/issues/274)

SDK Enhancements
---
* `service/s3/s3manager`: Improve memory allocation behavior by replacing sync.Pool with custom pool implementation
  * Improves memory allocations that occur when the provided `io.Reader` to upload does not satisfy both the `io.ReaderAt` and `io.ReadSeeker` interfaces.

SDK Bugs
---
* `service/s3/s3manager`: Fix resource leaks when the following occurred:
  * Failed CreateMultipartUpload call
  * Failed UploadPart call

