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

Services
---

SDK Features
---

SDK Enhancements
---
* `internal/sdkio`: Adds RingBuffer data structure to the sdk [#417](https://github.com/aws/aws-sdk-go-v2/pull/417)
  * Adds an implementation of RingBuffer data structure which acts as a revolving buffer of a predefined length. The RingBuffer implements io.ReadWriter interface.
  * Adds unit tests to test the behavior of the ring buffer. 
* `aws/ec2metadata`: Adds support for EC2Metadata client to use secure tokens provided by the IMDS ([#453](https://github.com/aws/aws-sdk-go-v2/pull/453)) 
  * Modifies and adds tests to verify the behavior of the EC2Metadata client.
* `service/dynamodb/dynamodbattribute`: Adds clarifying docs on dynamodbattribute.UnixTime ([#464](https://github.com/aws/aws-sdk-go-v2/pull/464))
* `example/service/sts/assumeRole`: added sts assume role example ([#224](https://github.com/aws/aws-sdk-go-v2/pull/224))
  * Fixes [#157](https://github.com/aws/aws-sdk-go-v2/issues/157) by adding an example for Amazon STS assume role to retrieve credentials.

SDK Bugs
---
* `service/dynamodb/dynamodbattribute`: Fixes a panic when decoding into a map with a key string type alias. ([#465](https://github.com/aws/aws-sdk-go/pull/465))
  * Fixes [#410](https://github.com/aws/aws-sdk-go-v2/issues/410),  by adding support for keys that are string aliases.
