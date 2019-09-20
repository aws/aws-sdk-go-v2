Services
---

Deprecations
---
* `aws`: Removes plugin credential provider ([#391](https://github.com/aws/aws-sdk-go-v2/pull/391))
  * Removing plugin credential provider from the v2 SDK developer preview. This feature may be made available as a separate module.
* Removes support for deprecated Go versions ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
  * Removes support for Go version specific files from the SDK. Also removes irrelevant build tags, and updates the README.md file. 
  
SDK Features
---

SDK Enhancements
---
* `aws/request` : Retryer is now a named field on Request. ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
  
SDK Bugs
---
* `private/model/api`: Fixes broken test for code generation. ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
