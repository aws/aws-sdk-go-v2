Breaking Change
---
* `service`: Add generated service for wafregional and dynamodbstreams #463
  * Updates the wafregional and dynamodbstreams API clients to include all API operations, and types that were previously shared between waf and dynamodb API clients respectively. This update ensures that all API clients include all operations and types needed for that client, and shares no types with another client package.
  * To migrate your applications to use the updated wafregional and dynamodbstreams you'll need to update the package the impacted type is imported from to match the client the type is being used with.
* `aws`: Context has been added to EC2Metadata operations.([#461](https://github.com/aws/aws-sdk-go-v2/pull/461))
  * Also updates utilities that directly or indirectly depend on EC2Metadata client. Signer utilities, credential providers now take in context.
* `private/model`: Add utility for validating shape names for structs and enums for the service packages ([#471](https://github.com/aws/aws-sdk-go-v2/pull/471))
  * Fixes bug which allowed service package structs, enums to start with non alphabetic character 
  * Fixes the incorrect enum types in mediapackage service package, changing enum types __AdTriggersElement, __PeriodTriggersElement to AdTriggersElement, PeriodTriggersElement respectively.
* `aws`: Client, Metadata, and Request structures have been refactored to simplify the usage of resolved endpoints ([#473](https://github.com/aws/aws-sdk-go-v2/pull/473))
  * `aws.Client.Endpoint` struct member has been removed, and `aws.Request.Endpoint` struct member has been added of type `aws.Endpoint`
  * `aws.Client.Region` structure member has been removed

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Features
---
* `aws`: `PartitionID` has been added to `aws.Endpoint` structure, and is used by the endpoint resolver to indicate which AWS partition an endpoint was resolved for ([#473](https://github.com/aws/aws-sdk-go-v2/pull/473))
* `aws/endpoints`: Updated resolvers to populate `PartitionID` for a resolved `aws.Endpoint` ([#473](https://github.com/aws/aws-sdk-go-v2/pull/473))
* `service/s3`: Add support for Access Point resources
  * Adds support for using Access Point resource with Amazon S3 API operation calls. The Access Point resource are identified by an Amazon Resource Name (ARN).
  * To make operation calls to an S3 Access Point instead of a S3 Bucket, provide the Access Point ARN string as the value of the Bucket parameter. You can create an Access Point for your bucket with the Amazon S3 Control API. The Access Point ARN can be obtained from the S3 Control API. You should avoid building the ARN directly.

SDK Enhancements
---
* `internal/sdkio`: Adds RingBuffer data structure to the sdk [#417](https://github.com/aws/aws-sdk-go-v2/pull/417)
  * Adds an implementation of RingBuffer data structure which acts as a revolving buffer of a predefined length. The RingBuffer implements io.ReadWriter interface.
  * Adds unit tests to test the behavior of the ring buffer. 
* `aws/ec2metadata`: Adds support for EC2Metadata client to use secure tokens provided by the IMDS ([#453](https://github.com/aws/aws-sdk-go-v2/pull/453)) 
  * Modifies EC2Metadata client to use request context within its operations ([#462](https://github.com/aws/aws-sdk-go-v2/pull/462))
  * Reduces the default dialer timeout and response header timeout to help reduce latency for known issues with EC2Metadata client running inside a container
  * Modifies and adds tests to verify the behavior of the EC2Metadata client.
* `service/dynamodb/dynamodbattribute`: Adds clarifying docs on dynamodbattribute.UnixTime ([#464](https://github.com/aws/aws-sdk-go-v2/pull/464))
* `example/service/sts/assumeRole`: added sts assume role example ([#224](https://github.com/aws/aws-sdk-go-v2/pull/224))
  * Fixes [#157](https://github.com/aws/aws-sdk-go-v2/issues/157) by adding an example for Amazon STS assume role to retrieve credentials.

SDK Bugs
---
* `service/dynamodb/dynamodbattribute`: Fixes a panic when decoding into a map with a key string type alias. ([#465](https://github.com/aws/aws-sdk-go/pull/465))
  * Fixes [#410](https://github.com/aws/aws-sdk-go-v2/issues/410),  by adding support for keys that are string aliases.
