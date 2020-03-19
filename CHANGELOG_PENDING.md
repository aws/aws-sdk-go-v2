Breaking Change
---
* `service/s3/s3crypto`: Package and associated encryption/decryption clients have been removed from the SDK ([#511](https://github.com/aws/aws-sdk-go-v2/pull/511))
* `aws/endpoints`: Removes `DecodeModel` and `DecodeModelOptions` from the package ([#509](https://github.com/aws/aws-sdk-go-v2/pull/509))
* `aws/external`: Removes several export constants and types ([#508](https://github.com/aws/aws-sdk-go-v2/pull/508))
  * No longer exports AWS environment constants used by the external environment configuration loader
  * `DefaultSharedConfigProfile` is now defined an exported constant

Services
---

SDK Features
---

SDK Enhancements
---

SDK Bugs
---
