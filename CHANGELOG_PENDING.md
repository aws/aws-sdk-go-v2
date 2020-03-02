Breaking Change
---

Services
---

SDK Features
---
* `aws`: Add Support for additional credential providers and credential configuration chaining ([#488](https://github.com/aws/aws-sdk-go-v2/pull/488))
  * `aws/processcreds`: Adds Support for the Process Credential Provider
    * Fixes [#249](https://github.com/aws/aws-sdk-go-v2/issues/249)
  * `aws/stscreds`: Adds Support for the Web Identity Credential Provider
    * Fixes [#475](https://github.com/aws/aws-sdk-go-v2/issues/475)
    * Fixes [#338](https://github.com/aws/aws-sdk-go-v2/issues/338)
  * Adds Support for `credential_source`
    * Fixes [#274](https://github.com/aws/aws-sdk-go-v2/issues/274)

SDK Enhancements
---

SDK Bugs
---
