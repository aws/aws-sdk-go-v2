### SDK Features

### SDK Enhancements

### SDK Bugs
* `aws/ec2rolecreds`: Fix security creds path to include trailing slash ([#356](https://github.com/aws/aws-sdk-go-v2/pull/356))
  * Fixes the iamSecurityCredsPath var to include a trailing slash preventing redirects when making requests to the EC2 Instance Metadata service.
  * Fixes [#351](https://github.com/aws/aws-sdk-go-v2/issues/351)
