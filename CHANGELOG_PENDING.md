Services
---

Deprecations
---
* `service/s3/s3crypto`: Deprecates the crypto client from the SDK ([#394](https://github.com/aws/aws-sdk-go-v2/pull/394))
  * s3crypto client is now deprecated and may be removed from the future versions of the SDK. 
* `aws`: Removes plugin credential provider ([#391](https://github.com/aws/aws-sdk-go-v2/pull/391))
  * Removing plugin credential provider from the v2 SDK developer preview. This feature may be made available as a separate module.
* Removes support for deprecated Go versions ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
  * Removes support for Go version specific files from the SDK. Also removes irrelevant build tags, and updates the README.md file. 
  * Raises the minimum supported version to Go 1.11 for the SDK. Older versions may work, but are not actively supported
  
SDK Features
---
* `service/s3/s3manager`: Add Upload Buffer Provider ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Adds a new `BufferProvider` member for specifying how part data can be buffered in memory.
  * Windows platforms will now default to buffering 1MB per part to reduce contention when uploading files.
  * Non-Windows platforms will continue to employ a non-buffering behavior.
* `service/s3/s3manager`: Add Download Buffer Provider ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Adds a new `BufferProvider` member for specifying how part data can be buffered in memory when copying from the http response body.
  * Windows platforms will now default to buffering 1MB per part to reduce contention when downloading files.
  * Non-Windows platforms will continue to employ a non-buffering behavior.
* `service/dynamodb/dynamodbattribute`: New Encoder and Decoder Behavior for Empty Collections ([#401](https://github.com/aws/aws-sdk-go-v2/pull/401))
  * The `Encoder` and `Decoder` types have been enhanced to support the marshaling of empty structures, maps, and slices to and from their respective DynamoDB AttributeValues.
  * This change incorporates the behavior changes introduced via a marshal option in V1 ([#2834](https://github.com/aws/aws-sdk-go/pull/2834))

SDK Enhancements
---
* `internal/awsutil`: Add suppressing logging sensitive API parameters ([#398](https://github.com/aws/aws-sdk-go-v2/pull/398))
  * Adds suppressing logging sensitive API parameters marked with the `sensitive` trait. This prevents the API type's `String` method returning a string representation of the API type with sensitive fields printed such as keys and passwords.
  * Related to [aws/aws-sdk-go#2310](https://github.com/aws/aws-sdk-go/pull/2310)
  * Fixes [#251](https://github.com/aws/aws-sdk-go-v2/issues/251)
* `aws/request` : Retryer is now a named field on Request. ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
* `service/s3/s3manager`: Adds `sync.Pool` to allow reuse of part buffers for streaming payloads ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Fixes [#402](https://github.com/aws/aws-sdk-go-v2/issues/402) 
  * Uses the new behavior introduced in V1 [#2863](https://github.com/aws/aws-sdk-go/pull/2863) which allows the reuse of the sync.Pool across multiple Upload request that match part sizes.

SDK Bugs
---
* `service/s3/s3manager`: Fix index out of range when a streaming reader returns -1 ([#378](https://github.com/aws/aws-sdk-go-v2/pull/378))
  * Fixes the S3 Upload Manager's handling of an unbounded streaming reader that returns negative bytes read.
* `internal/ini`: Fix ini parser to handle empty values [#406](https://github.com/aws/aws-sdk-go-v2/pull/406)
  * Fixes incorrect modifications to the previous token value of the skipper. Adds checks for cases where a skipped statement should be marked as complete and not be ignored.
  * Adds tests for nested and empty field value parsing, along with tests suggested in [aws/aws-sdk-go#2801](https://github.com/aws/aws-sdk-go/pull/2801)
