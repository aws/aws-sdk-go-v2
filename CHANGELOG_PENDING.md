### SDK Features

### SDK Enhancements
* `aws`: Provides more customization options for retryer by adding a constructor for default Retryer which accepts functional options. Adds NoOpRetryer to support no retry behavior. Exposes members of default retryer. [#375](https://github.com/aws/aws-sdk-go-v2/pull/375)
    * Updates the underlying logic used by the default retryer to calculate jittered delay for retry. 
    * Handles int overflow for retry delay. Fixes [#370](https://github.com/aws/aws-sdk-go-v2/issues/370)
* `service/ec2`: Adds custom retryer implementation for service/ec2.
  * Adds test case to test custom retryer. 
  
### SDK Bugs

  