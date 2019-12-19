Services
---

SDK Breaking Changes
---
  * `aws`: Context has been added to EC2Metadata operations and credential provider operations that use EC2Metadata.([#460](https://github.com/aws/aws-sdk-go-v2/pull/460)) 

SDK Features
---

SDK Enhancements
---
* `aws/ec2metadata`: Adds support for EC2Metadata client to use secure tokens provided by the IMDS ([#453](https://github.com/aws/aws-sdk-go-v2/pull/453)) 
  * The dialer timeout is set to 250 ms to reduce latency in fetching responses when running EC2Metadata client in a container ([#460](https://github.com/aws/aws-sdk-go-v2/pull/460))
  * Modifies and adds tests to verify the behavior of the EC2Metadata client

SDK Bugs
--
