### SDK Features

### SDK Enhancements
* `private/protocol`: Add support for TimestampFormat in protocols ([#358](https://github.com/aws/aws-sdk-go-v2/pull/358))
    * Adds support for the timestampForamt API model trait to the V2 SDK. The SDK will now generate API client parameters with the correct time format for APIs modeled with custom time stamp formats specified.
    * Fixes [#202](https://github.com/aws/aws-sdk-go-v2/issues/202)
    * Fix [#286](https://github.com/aws/aws-sdk-go-v2/issues/286)

### SDK Bugs
* `private/model/api`: Fix API doc being generated with wrong value ([#359](https://github.com/aws/aws-sdk-go-v2/pull/359))
  * Fixes the SDK's generated API documentation for structure member being generated with the wrong documentation value when the member was included multiple times in the model doc-2.json file, but under different types.
  * V2 port of to v1 [aws/aws-sdk-go#2748](https://github.com/aws/aws-sdk-go/issues/2748)
* `aws/ec2rolecreds`: Fix security creds path to include trailing slash ([#356](https://github.com/aws/aws-sdk-go-v2/pull/356))
  * Fixes the iamSecurityCredsPath var to include a trailing slash preventing redirects when making requests to the EC2 Instance Metadata service.
  * Fixes [#351](https://github.com/aws/aws-sdk-go-v2/issues/351)



