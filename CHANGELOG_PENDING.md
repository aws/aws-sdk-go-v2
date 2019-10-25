Services
---

SDK Features
---

SDK Enhancements
---

SDK Bugs
---
* `aws/endpoints`: aws/endpoints: Fix SDK resolving endpoint without region ([#420](https://github.com/aws/aws-sdk-go-v2/pull/420))
  * Fixes the SDK's endpoint resolve incorrectly resolving endpoints for a service when the region is empty. Also fixes the SDK attempting to resolve a service when the service value is empty..
  * Related to [aws/aws-sdk-go#2909](https://github.com/aws/aws-sdk-go/issues/2909)
