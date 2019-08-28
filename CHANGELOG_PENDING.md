### SDK Features

### SDK Enhancements

### SDK Bugs
* `aws`: Fixes bug where the throttled retry's math was off 
    * The `Retry-After` duration specified in the request is now added to the Retry delay for throttled exception
    * Adds test for retry delays for throttled exceptions [#370](https://github.com/aws/aws-sdk-go/pull/370)
    * Fixes [#45](https://github.com/aws/aws-sdk-go/issues/45)

