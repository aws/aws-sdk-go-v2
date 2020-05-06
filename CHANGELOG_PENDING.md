Services
---

SDK Features
---

SDK Enhancements
---

SDK Bugs
---
* `aws/defaults`: Fix handling of unexpected Date response header value ([#](https://github.com/aws/aws-sdk-go-v2/pull/560))
  * Fixes the SDK's behavior to parse unexpected HTTP Date header received that was formated with single digit day of the month instead of two digit RFC822 datetime like defined in RFC 2616. This should prevent log messages about unable to compute clock skew.
  * Fixes [#](https://github.com/aws/aws-sdk-go-v2/issues/556)
`service/s3`: Fix S3 client behavior wrt 200 OK response with empty payload
