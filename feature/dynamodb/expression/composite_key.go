package expression

import (
	"fmt"
	"strings"
)

// compositeKeyCondition builds the composite key condition expression for DynamoDB queries.
// It counts partition and sort keys, validates the presence of at least one partition key,
// and formats the expression accordingly.
func compositeKeyCondition(keyConditionBuilder KeyConditionBuilder, node exprNode) (exprNode, error) {
	pks := 0
	sks := 0
	for _, kcb := range keyConditionBuilder.keyConditionList {
		switch kcb.compositeKeyMode {
		case compositeKeyCondPartition:
			pks++
		case compositeKeyCondSort:
			sks++
		}
	}

	if pks == 0 {
		return exprNode{}, newInvalidParameterError("compositeKeyCondition", "KeyConditionBuilder")
	}

	node.fmtExpr = fmt.Sprintf(
		"(%s)",
		strings.Repeat(" AND ($c)", pks+sks)[5:],
	)

	return node, nil
}

// CompositeKey initializes a KeyConditionBuilder in composite key mode.
//
// Note: When using CompositeKey, you must use AddPartitionKey and AddSortKey in the desired order.
// The order in which you call these methods determines the order of keys in the resulting expression.
// The And() method is not supported for composite keys.
//
// Example:
//
//	// Create a composite key condition with partition and sort keys
//	partitionKey := expression.Key("TeamName").Equal(expression.Value("Wildcats"))
//	sortKey := expression.Key("Number").Equal(expression.Value(1))
//	keyCondition := expression.CompositeKey().AddPartitionKey(partitionKey).AddSortKey(sortKey)
//	builder := expression.NewBuilder().WithKeyCondition(keyCondition)
func CompositeKey() KeyConditionBuilder {
	return KeyConditionBuilder{
		mode: compositeKeyCond,
	}
}

// AddPartitionKey adds a partition key condition to the composite key builder.
//
// Note: The order in which AddPartitionKey and AddSortKey are called determines
// the order of keys in the resulting expression. Do not change the order after building.
//
// Example:
//
//	partitionKey := expression.Key("TeamName").Equal(expression.Value("Wildcats"))
//	keyCondition := expression.CompositeKey().AddPartitionKey(partitionKey)
func (kcb KeyConditionBuilder) AddPartitionKey(pk KeyConditionBuilder) KeyConditionBuilder {
	pk.compositeKeyMode = compositeKeyCondPartition
	kcb.keyConditionList = append(kcb.keyConditionList, pk)

	return kcb
}

// AddSortKey adds a sort key condition to the composite key builder.
//
// Note: The order in which AddPartitionKey and AddSortKey are called determines
// the order of keys in the resulting expression. Do not change the order after building.
//
// IMPORTANT: You cannot use AddSortKey alone; you must also add at least one partition key
// using AddPartitionKey, or the query will fail.
//
// Example:
//
//	partitionKey := expression.Key("TeamName").Equal(expression.Value("Wildcats"))
//	sortKey := expression.Key("Number").Equal(expression.Value(1))
//	keyCondition := expression.CompositeKey().AddPartitionKey(partitionKey).AddSortKey(sortKey)
func (kcb KeyConditionBuilder) AddSortKey(sk KeyConditionBuilder) KeyConditionBuilder {
	sk.compositeKeyMode = compositeKeyCondSort
	kcb.keyConditionList = append(kcb.keyConditionList, sk)

	return kcb
}
