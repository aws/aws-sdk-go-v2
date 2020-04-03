Breaking Change
---

Services
---

SDK Features
---

SDK Enhancements
---
* `service/dynamodb/expression`: Add IsSet helper for ConditionBuilder and KeyConditionBuilder ([#494](https://github.com/aws/aws-sdk-go-v2/pull/494))
  * Adds a IsSet helper for ConditionBuilder and KeyConditionBuilder to make it easier to determine if the condition builders have any conditions added to them.
  * Implements [#493](https://github.com/aws/aws-sdk-go-v2/issues/493).

SDK Bugs
---
