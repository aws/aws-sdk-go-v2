package enhancedclient

import (
	"context"
	"fmt"
	"iter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetItem retrieves a single item from the DynamoDB table by its key.
// Returns the decoded item or an error if not found or decoding fails.
func (t *Table[T]) GetItem(ctx context.Context, m Map, optFns ...func(*dynamodb.Options)) (*T, error) {
	res, err := t.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: t.options.Schema.TableName(),
		Key:       m,
	}, optFns...)
	if err != nil {
		return nil, err
	}

	if res == nil || res.Item == nil {
		return nil, fmt.Errorf("empty response or item in GetItem() call")
	}

	item, err := t.options.Schema.Decode(res.Item)
	if err != nil {
		return nil, err
	}

	err = t.applyAfterReadExtensions(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// GetItemWithProjection retrieves a single item from the DynamoDB table by its key, applying a projection to select specific attributes.
// Returns the decoded item or an error if not found or decoding fails.
func (t *Table[T]) GetItemWithProjection(ctx context.Context, m Map, proj expression.ProjectionBuilder, optFns ...func(*dynamodb.Options)) (*T, error) {
	b, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, err
	}

	res, err := t.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:                t.options.Schema.TableName(),
		Key:                      m,
		ExpressionAttributeNames: b.Names(),
		ProjectionExpression:     b.Projection(),
	}, optFns...)
	if err != nil {
		return nil, err
	}

	if res == nil || res.Item == nil {
		return nil, fmt.Errorf("empty response or item in GetItemWithProjection() call")
	}

	item, err := t.options.Schema.Decode(res.Item)
	if err != nil {
		return nil, err
	}

	err = t.applyAfterReadExtensions(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// PutItem writes the item to the DynamoDB table without checking for collisions.
// Returns the written item or an error if encoding or writing fails.
func (t *Table[T]) PutItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) (*T, error) {
	err := t.applyBeforeWriteExtensions(item)
	if err != nil {
		return nil, err
	}

	itemMap, err := t.options.Schema.Encode(item)
	if err != nil {
		return nil, err
	}

	res, err := t.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: t.options.Schema.TableName(),
		Item:      itemMap,
	}, optFns...)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response in PutItem() call")
	}

	out, err := t.options.Schema.Decode(itemMap)
	if err != nil {
		return nil, err
	}

	if err := t.applyAfterWriteExtensions(out); err != nil {
		return nil, err
	}

	return out, nil
}

// UpdateItem writes the item to the DynamoDB table with additional checks (e.g., version checks).
// Returns the updated item or an error if encoding or updating fails.
func (t *Table[T]) UpdateItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) (*T, error) {
	err := t.applyBeforeWriteExtensions(item)
	if err != nil {
		return nil, err
	}

	m, err := t.options.Schema.createKeyMap(item)
	if err != nil {
		return nil, err
	}

	expr, err := t.createUpdateExpression(item)
	if err != nil {
		return nil, err
	}

	res, err := t.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                           t.options.Schema.TableName(),
		Key:                                 m,
		ConditionExpression:                 expr.Condition(),
		ExpressionAttributeNames:            expr.Names(),
		ExpressionAttributeValues:           expr.Values(),
		UpdateExpression:                    expr.Update(),
		ReturnValues:                        types.ReturnValueAllNew,
		ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
	}, optFns...)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("empty response in UpdateItem() call")
	}

	out, err := t.options.Schema.Decode(res.Attributes)
	if err != nil {
		return nil, err
	}

	if err := t.applyAfterWriteExtensions(out); err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteItem deletes an item from the DynamoDB table by its struct value.
// Returns an error if the key cannot be created or the delete fails.
func (t *Table[T]) DeleteItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) error {
	m, err := t.options.Schema.createKeyMap(item)
	if err != nil {
		return err
	}

	_, err = t.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: t.options.Schema.TableName(),
		Key:       m,
	}, optFns...)

	return err
}

// DeleteItemByKey deletes an item from the DynamoDB table by its key map.
// Returns an error if the delete fails.
func (t *Table[T]) DeleteItemByKey(ctx context.Context, m Map, optFns ...func(*dynamodb.Options)) error {
	_, err := t.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: t.options.Schema.TableName(),
		Key:       m,
	}, optFns...)

	return err
}

// createScanIterator returns an iterator that scans a DynamoDB table or index and yields results as ItemResult[T].
// It automatically handles pagination and error thresholds using MaxConsecutiveErrors.
// If the number of consecutive errors reaches the threshold, iteration stops.
func (t Table[T]) createScanIterator(ctx context.Context, indexName *string, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	var consecutiveErrors uint = 0
	var maxConsecutiveErrors = t.options.MaxConsecutiveErrors
	if maxConsecutiveErrors == 0 {
		maxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}

	return func(yield func(ItemResult[T]) bool) {
		var lastEvaluatedKey map[string]types.AttributeValue

		for {
			scanInput := &dynamodb.ScanInput{
				TableName:                 t.options.Schema.TableName(),
				IndexName:                 indexName,
				ConsistentRead:            aws.Bool(indexName == nil),
				ExclusiveStartKey:         lastEvaluatedKey,
				Select:                    types.SelectAllAttributes,
				FilterExpression:          expr.Filter(),
				ProjectionExpression:      expr.Projection(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
			}

			res, err := t.client.Scan(ctx, scanInput, optFns...)
			if err != nil {
				consecutiveErrors++

				if !yield(ItemResult[T]{err: err}) {
					return
				}

				if consecutiveErrors >= maxConsecutiveErrors {
					return
				}

				continue
			}

			consecutiveErrors = 0

			if res != nil && res.Items != nil {
				for _, item := range res.Items {
					i, err := t.options.Schema.Decode(item)
					if err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if err := t.applyAfterReadExtensions(i); err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if !yield(ItemResult[T]{item: i}) {
						return
					}
				}

				lastEvaluatedKey = res.LastEvaluatedKey
			} else {
				lastEvaluatedKey = nil
			}

			if lastEvaluatedKey == nil {
				return
			}
		}
	}
}

// ScanIndex scans a DynamoDB index and returns an iterator of results.
// The scan uses the provided index name and expression.
func (t *Table[T]) ScanIndex(ctx context.Context, indexName string, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return t.createScanIterator(ctx, &indexName, expr, optFns...)
}

// Scan scans the DynamoDB table and returns an iterator of results.
// The scan uses the provided expression.
func (t *Table[T]) Scan(ctx context.Context, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return t.createScanIterator(ctx, nil, expr, optFns...)
}

// createQueryIterator returns an iterator that queries a DynamoDB table or index and yields results as ItemResult[T].
// It automatically handles pagination and error thresholds using MaxConsecutiveErrors.
// If the number of consecutive errors reaches the threshold, iteration stops.
func (t *Table[T]) createQueryIterator(ctx context.Context, indexName *string, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	var consecutiveErrors uint = 0
	var maxConsecutiveErrors = t.options.MaxConsecutiveErrors
	if maxConsecutiveErrors == 0 {
		maxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}

	return func(yield func(ItemResult[T]) bool) {
		var lastEvaluatedKey map[string]types.AttributeValue

		for {
			res, err := t.client.Query(ctx, &dynamodb.QueryInput{
				TableName:                 t.options.Schema.TableName(),
				IndexName:                 indexName,
				ConsistentRead:            aws.Bool(indexName == nil),
				ExclusiveStartKey:         lastEvaluatedKey,
				KeyConditionExpression:    expr.KeyCondition(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				FilterExpression:          expr.Filter(),
				ProjectionExpression:      expr.Projection(),
				Select:                    types.SelectAllAttributes,
			}, optFns...)

			if err != nil {
				consecutiveErrors++

				if !yield(ItemResult[T]{err: err}) {
					return
				}

				if consecutiveErrors >= maxConsecutiveErrors {
					return
				}

				continue
			}

			consecutiveErrors = 0

			if res == nil {
				return
			}

			if res != nil && res.Items != nil {
				for _, item := range res.Items {
					i, err := t.options.Schema.Decode(item)
					if err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if err := t.applyAfterReadExtensions(i); err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if !yield(ItemResult[T]{item: i}) {
						return
					}
				}

				lastEvaluatedKey = res.LastEvaluatedKey
			} else {
				lastEvaluatedKey = nil
			}

			if lastEvaluatedKey == nil {
				return
			}
		}
	}
}

// QueryIndex queries a DynamoDB index and returns an iterator of results.
// The query uses the provided index name and expression.
func (t *Table[T]) QueryIndex(ctx context.Context, indexName string, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return t.createQueryIterator(ctx, &indexName, expr, optFns...)
}

// Query queries the DynamoDB table and returns an iterator of results.
// The query uses the provided expression.
func (t *Table[T]) Query(ctx context.Context, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return t.createQueryIterator(ctx, nil, expr, optFns...)
}

// CreateBatchWriteOperation creates a new BatchWriteOperation for the table.
// Use this to perform batched put and delete operations for the table's items.
func (t *Table[T]) CreateBatchWriteOperation() *BatchWriteOperation[T] {
	return NewBatchWriteOperation(t)
}

// CreateBatchGetOperation creates a new BatchGetOperation for the table.
// Use this to perform batched reads for the table's items.
func (t *Table[T]) CreateBatchGetOperation() *BatchGetOperation[T] {
	return NewBatchGetOperation(t)
}
