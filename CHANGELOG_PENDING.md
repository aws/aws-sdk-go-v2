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

SDK Enhancements
---
* `internal/awsutil`: Add suppressing logging sensitive API parameters ([#398](https://github.com/aws/aws-sdk-go-v2/pull/398))
  * Adds suppressing logging sensitive API parameters marked with the `sensitive` trait. This prevents the API type's `String` method returning a string representation of the API type with sensitive fields printed such as keys and passwords.
  * Related to [aws/aws-sdk-go#2310](https://github.com/aws/aws-sdk-go/pull/2310)
  * Fixes [#251](https://github.com/aws/aws-sdk-go-v2/issues/251)
* `aws/request` : Retryer is now a named field on Request. ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))

SDK Bugs
---
* `service/s3/s3manager`: Fix index out of range when a streaming reader returns -1 ([#378](https://github.com/aws/aws-sdk-go-v2/pull/378))
  * Fixes the S3 Upload Manager's handling of an unbounded streaming reader that returns negative bytes read.
