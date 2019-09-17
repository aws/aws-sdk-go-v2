Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Features
---

SDK Enhancements
---
* `aws/endpoints`: Expose DNSSuffix for partitions ([#369](https://github.com/aws/aws-sdk-go-v2/pull/369))
  * Exposes the underlying partition metadata's DNSSuffix value via the `DNSSuffix` method on the endpoint's `Partition` type. This allows access to the partition's DNS suffix, e.g. "amazon.com".
  * Fixes [#347](https://github.com/aws/aws-sdk-go-v2/issues/347)
* `private/protocol`: Add support for parsing fractional timestamp ([#367](https://github.com/aws/aws-sdk-go-v2/pull/367))
  * Fixes the SDK's ability to parse fractional unix timestamp values and adds tests.
  * Fixes [#365](https://github.com/aws/aws-sdk-go-v2/issues/365)
* `aws/ec2metadata`: Add marketplaceProductCodes to EC2 Instance Identity Document ([#374](https://github.com/aws/aws-sdk-go-v2/pull/374))
  * Adds `MarketplaceProductCodes` to the EC2 Instance Metadata's Identity Document. The ec2metadata client will now retrieve these values if they are available.
  * Related to: [aws/aws-sdk-go#2781](https://github.com/aws/aws-sdk-go/issues/2781)
* `aws`: Adds configurations to the default retryer ([#375](https://github.com/aws/aws-sdk-go-v2/pull/375))
  * Provides more customization options for retryer by adding a constructor for default Retryer which accepts functional options. Adds NoOpRetryer to support no retry behavior. Exposes members of default retryer.
  * Updates the underlying logic used by the default retryer to calculate jittered delay for retry. Handles int overflow for retry delay. 
  * Fixes [#370](https://github.com/aws/aws-sdk-go-v2/issues/370)
* `aws` : Refactors request retry behavior path logic ([#384](https://github.com/aws/aws-sdk-go-v2/pull/384))
  * Retry utilities now follow a consistent code path. aws.IsErrorRetryable is the primary entry point to determine if a request is retryable. 
  * Corrects sdk's behavior by not retrying errors with status code 501. Adds support for retrying the Kinesis API error, LimitExceededException.
  * Fixes [#372](https://github.com/aws/aws-sdk-go-v2/issues/372), [#145](https://github.com/aws/aws-sdk-go-v2/issues/145)
  
SDK Bugs
---
* `aws`: Fixes bug in calculating throttled retry delay ([#373](https://github.com/aws/aws-sdk-go-v2/pull/373))
  * The `Retry-After` duration specified in the request is now added to the Retry delay for throttled exception. Adds test for retry delays for throttled exceptions. Fixes bug where the throttled retry's math was off.
  * Fixes [#45](https://github.com/aws/aws-sdk-go-v2/issues/45)
* `aws` : Adds missing sdk error checking when seeking readers ([#379](https://github.com/aws/aws-sdk-go-v2/pull/379))
  * Adds support for nonseekable io.Reader. Adds support for streamed payloads for unsigned body request. 
  * Fixes [#371](https://github.com/aws/aws-sdk-go-v2/issues/371)
* `service/s3` : Fixes unexpected EOF error by s3manager ([#386](https://github.com/aws/aws-sdk-go-v2/pull/386))
  * Fixes bug which threw unexpected EOF error when s3 upload is performed for a file of maximum allowed size
  * Fixes [#316](https://github.com/aws/aws-sdk-go-v2/issues/316)
* `private/model` : Fixes generated API Reference docs links being invalid ([387](https://github.com/aws/aws-sdk-go-v2/pull/387))
  * Fixes [#327](https://github.com/aws/aws-sdk-go-v2/issues/327)
