Breaking Change
---
* `aws/endpoints`: Several functions and types have been removed
  * Removes `DecodeModel` and `DecodeModelOptions` from the package ([#509](https://github.com/aws/aws-sdk-go-v2/pull/509))
  * Remove Region Constants, Partition Constants, and types use for exploring the endpoint data model ([#512](https://github.com/aws/aws-sdk-go-v2/pull/512))
* `service/s3/s3crypto`: Package and associated encryption/decryption clients have been removed from the SDK ([#511](https://github.com/aws/aws-sdk-go-v2/pull/511))
* `aws/external`: Removes several export constants and types ([#508](https://github.com/aws/aws-sdk-go-v2/pull/508))
  * No longer exports AWS environment constants used by the external environment configuration loader
  * `DefaultSharedConfigProfile` is now defined an exported constant
* `aws`: `ErrMissingRegion`, `ErrMissingEndpoint`, `ErrStaticCredentialsEmpty` are now concrete error types ([#510](https://github.com/aws/aws-sdk-go-v2/pull/510))

Services
---

SDK Features
---
* `aws/signer/v4`: New methods `SignHTTP` and `PresignHTTP` have been added ([#519](https://github.com/aws/aws-sdk-go-v2/pull/519))
  * `SignHTTP` replaces `Sign`, and usage of `Sign` should be migrated before it's removal at a later date
  * `PresignHTTP` replaces `Presign`, and usage of `Presign` should be migrated before it's removal at a later date
  * `DisableRequestBodyOverwrite` and `UnsignedPayload` are now deprecated options and have no effect on `SignHTTP` or `PresignHTTP`. These options will be removed at a later date.
* `aws/external`: Add Support for setting a default fallback region and resolving region from EC2 IMDS ([#523](https://github.com/aws/aws-sdk-go-v2/pull/523))
  * `WithDefaultRegion` helper has been added which can be passed to `LoadDefaultAWSConfig`
    * This helper can be used to configure a default fallback region in the event a region fails to be resolved from other sources
  * Support has been added to resolve region using EC2 IMDS when available
    * The IMDS region will be used if region as not found configured in either the shared config or the process environment.
  * Fixes [#244](https://github.com/aws/aws-sdk-go-v2/issues/244)
  * Fixes [#515](https://github.com/aws/aws-sdk-go-v2/issues/515)
SDK Enhancements
---
* `service/dynamodb/expression`: Add IsSet helper for ConditionBuilder and KeyConditionBuilder ([#494](https://github.com/aws/aws-sdk-go-v2/pull/494))
  * Adds a IsSet helper for ConditionBuilder and KeyConditionBuilder to make it easier to determine if the condition builders have any conditions added to them.
  * Implements [#493](https://github.com/aws/aws-sdk-go-v2/issues/493).
* `internal/ini`: Normalize Section keys to lowercase ([#495](https://github.com/aws/aws-sdk-go-v2/pull/495))
  * Update's SDK's ini utility to store all keys as lowercase. This brings the SDK inline with the AWS CLI's behavior.


SDK Bugs
---
