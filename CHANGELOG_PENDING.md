Services
---

Deprecations
---
* `aws`: Removes plugin credential provider ([#391](https://github.com/aws/aws-sdk-go-v2/pull/391))
  * Removing plugin credential provider from the v2 SDK developer preview. This feature may be made available as a separate module.

SDK Features
---

SDK Enhancements
---
* `internal/awsutil`: Add suppressing logging sensitive API parameters ([#398](https://github.com/aws/aws-sdk-go-v2/pull/398))
  * Adds suppressing logging sensitive API parameters marked with the `sensitive` trait. This prevents the API type's `String` method returning a string representation of the API type with sensitive fields printed such as keys and passwords.
  * Related to [aws/aws-sdk-go#2310](https://github.com/aws/aws-sdk-go/pull/2310)
  * Fixes [#251](https://github.com/aws/aws-sdk-go-v2/issues/251)

SDK Bugs
---
