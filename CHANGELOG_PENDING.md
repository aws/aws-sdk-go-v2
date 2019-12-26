Services
---

SDK Features
---

SDK Enhancements
---
* `aws/ec2metadata`: Adds support for EC2Metadata client to use secure tokens provided by the IMDS ([#453](https://github.com/aws/aws-sdk-go-v2/pull/453)) 
  * Modifies and adds tests to verify the behavior of the EC2Metadata client.

SDK Bugs
--
* `service/dynamodb/dynamodbattribute`: Fixes a panic when decoding into a map with a key string type alias. ([#465](https://github.com/aws/aws-sdk-go/pull/465))
