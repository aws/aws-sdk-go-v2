Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Features
---

SDK Enhancements
---

SDK Bugs
---
* `service/dynamodb/expression: fix empty expression returned when unset ([#562](https://github.com/aws/aws-sdk-go-v2/pull/562))
  * Fixes a big in the expression builder that returns an empty string expression value when the expression has not been set.
  * Fixes [#554](https://github.com/aws/aws-sdk-go-v2/issues/554)
