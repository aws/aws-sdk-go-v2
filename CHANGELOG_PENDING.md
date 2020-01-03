Breaking Change
---
* Add generated service for wafregional and dynamodbstreams #463
  * Updates the wafregional and dynamodbstreams API clients to include all API operations, and types that were previously shared between waf and dynamodb API clients respectively. This update ensures that all API clients include all operations and types needed for that client, and shares no types with another client package.
  * To migrate your applications to use the updated wafregional and dynamodbstreams you'll need to update the package the impacted type is imported from to match the client the type is being used with.

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

