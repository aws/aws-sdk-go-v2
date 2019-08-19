### SDK Features

### SDK Enhancements

### SDK Bugs
* `private/model/api`: Fix API doc being generated with wrong value ([#359](https://github.com/aws/aws-sdk-go-v2/pull/359))
  * Fixes the SDK's generated API documentation for structure member being generated with the wrong documentation value when the member was included multiple times in the model doc-2.json file, but under different types.
