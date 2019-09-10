### SDK Features

### SDK Enhancements
* `aws/endpoints`: Expose DNSSuffix for partitions ([#368](https://github.com/aws/aws-sdk-go/pull/368))
  * Exposes the underlying partition metadata's DNSSuffix value via the `DNSSuffix` method on the endpoint's `Partition` type. This allows access to the partition's DNS suffix, e.g. "amazon.com".
  * Fixes [#347](https://github.com/aws/aws-sdk-go/issues/347)
* `private/protocol`: Add support for parsing fractional timestamp ([#367](https://github.com/aws/aws-sdk-go-v2/pull/367))
  * Fixes the SDK's ability to parse fractional unix timestamp values and adds tests.
  * Fixes [#365](https://github.com/aws/aws-sdk-go-v2/issues/365)
* `aws/ec2metadata`: Add marketplaceProductCodes to EC2 Instance Identity Document
  * Adds `MarketplaceProductCodes` to the EC2 Instance Metadata's Identity Document. The ec2metadata client will now retrieve these values if they are available.
  * Related to: [aws/aws-sdk-go#2781](https://github.com/aws/aws-sdk-go/issues/2781)

### SDK Bugs
* `aws`: Fixes bug in calculating throttled retry delay ([#373](https://github.com/aws/aws-sdk-go-v2/pull/373))
  * The `Retry-After` duration specified in the request is now added to the Retry delay for throttled exception. Adds test for retry delays for throttled exceptions. Fixes bug where the throttled retry's math was off.
  * Fixes [#45](https://github.com/aws/aws-sdk-go-v2/issues/45)
* `aws` : Adds missing sdk error checking when seeking readers [#379](https://github.com/aws/aws-sdk-go-v2/pull/379).
  * Adds support for nonseekable io.Reader. Adds support for streamed payloads for unsigned body request. 
  * Fixes [#371](https://github.com/aws/aws-sdk-go-v2/issues/371)
  
