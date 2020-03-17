Breaking Change
---

* Update SDK retry behavior 
  * Significant updates were made the SDK's retry behavior. The SDK will now retry all connections error. In addition, to changing what errors are retried the SDK's retry behavior not distinguish the difference between throttling errors, and regular retryable errors. All errors will be retried with the same backoff jitter delay scaling.
  * The SDK will attempt an operation request 3 times by default. This is one less than the previous initial request with 3 retires.
  * New helper functions in the new `aws/retry` package allow wrapping a `Retrier` with custom behavior, overriding the base `Retrier`, (e.g. `AddWithErrorCodes`, and `AddWithMaxAttempts`)
* Update SDK error handling
  * Updates the SDK's handling of errors to take advantage of Go 1.13's new `errors.As`, `Is`, and `Unwrap`. The SDK's errors were updated to satisfy the `Unwrap` interface, returning the underlying error.
  * With this update, you can now more easily access the SDK's layered errors, and meaningful state such as, `Timeout`, `Temporary`, and other states added to the SDK such as `CanceledError`.
* Bump SDK minimum supported version from Go 1.12 to Go 1.13
  * The SDK's minimum supported version is bumped to take advantage of Go 1.13's updated `errors` package.

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
* `aws/awserr`: Adds support for Go 1.13's `errors.Unwrap` ([#487](https://github.com/aws/aws-sdk-go-v2/pull/487))
* `aws`: Updates SDK retry behavior ([#487](https://github.com/aws/aws-sdk-go-v2/pull/487))
  * `aws/retry`: New package defining logic to determine if a request should be retried, and backoff.
  * `aws/ratelimit`: New package defining rate limit logic such as token bucket used by the `retry.Standard` retrier.

SDK Enhancements
---
* `aws`: Add grouping of concurrent refresh of credentials ([#503](https://github.com/aws/aws-sdk-go-v2/pull/503)
  * Concurrent calls to `Retrieve` are now grouped in order to prevent numerous synchronous calls to refresh the credentials. Replacing the mutex with a singleflight reduces the overall amount of time request signatures need to wait while retrieving credentials. This is improvement becomes pronounced when many requests are made concurrently.
* `service/s3/s3manager`: Improve memory allocation behavior by replacing sync.Pool with custom pool implementation
  * Improves memory allocations that occur when the provided `io.Reader` to upload does not satisfy both the `io.ReaderAt` and `io.ReadSeeker` interfaces.

SDK Bugs
---
* `service/s3/s3manager`: Fix resource leaks when the following occurred:
  * Failed CreateMultipartUpload call
  * Failed UploadPart call

