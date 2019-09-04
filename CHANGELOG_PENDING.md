### SDK Features

### SDK Enhancements

### SDK Bugs
* `aws` : Adds missing sdk error checking when seeking readers [#379](https://github.com/aws/aws-sdk-go-v2/pull/379).
  * Adds support for nonseekable io.Reader. Adds support for streamed payloads for unsigned body request. 
  * Fixes [#371](https://github.com/aws/aws-sdk-go-v2/issues/371)
  
