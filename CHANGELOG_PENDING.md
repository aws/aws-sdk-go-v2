Services
---

SDK Features
---

SDK Enhancements
---
* `aws/ec2metadata`: Adds support for EC2Metadata client to use secure tokens provided by the IMDS ([#453](https://github.com/aws/aws-sdk-go-v2/pull/453)) 
  * Modifies and adds tests to verify the behavior of the EC2Metadata client.
* `service/dynamodb/dynamodbattribute`: Adds clarifying docs on dynamodbattribute.UnixTime ([#464](https://github.com/aws/aws-sdk-go-v2/pull/464))
* `example/service/sts/assumeRole`: added sts assume role example ([#224](https://github.com/aws/aws-sdk-go-v2/pull/224))
  * Fixes [#157](https://github.com/aws/aws-sdk-go-v2/issues/157) by adding an example for Amazon STS assume role to retrieve credentials.

SDK Bugs
---
* `service/dynamodb/dynamodbattribute`: Fixes a panic when decoding into a map with a key string type alias. ([#465](https://github.com/aws/aws-sdk-go/pull/465))
  * Fixes [#410](https://github.com/aws/aws-sdk-go-v2/issues/410),  by adding support for keys that are string aliases.
