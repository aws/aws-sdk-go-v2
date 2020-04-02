Breaking Change
---

Services
---

SDK Features
---

SDK Enhancements
---
* `service/dynamodb/expression`: Add IsSet method for ConditionBuilder and KeyConditionBuilder ([#494](https://github.com/aws/aws-sdk-go-v2/pull/494))
  * Adds IsSet method to the ConditionBuilder and KeyConditionBuilder types. This methods makes it easier to discover if the condition builders have any conditions added to them.
  * Implements [#493](https://github.com/aws/aws-sdk-go-v2/issues/493).

SDK Bugs
---
