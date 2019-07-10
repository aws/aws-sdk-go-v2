### SDK Features

### SDK Enhancements

### SDK Bugs
* `private/model`: Handles empty map vs unset map behavior in send request ([#337](https://github.com/aws/aws-sdk-go-v2/pull/337))
  * Updates shape marshal model to handle the empty map vs nil map behavior. Also adds a test case to assert behavior when a user sends an empty map vs unset map.