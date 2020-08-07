Services
---

SDK Features
---

SDK Enhancements
---

SDK Bugs
---
* `aws/external`: Fixed a bug where an error from `EC2RoleCredentialProviderOptions` would not be returned during credential resolution. ([#659](https://github.com/aws/aws-sdk-go-v2/pull/659))
  * Additionally, fixes a bug that would prevent config sources that implement `EndpointCredentialProviderOptions`  from being used when resolving credentials from an HTTP provider.
