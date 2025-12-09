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

// PutItem writes the item without checking for collisions
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

// UpdateItem writes the item with additional checks (version checks, etc)
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

func (t *Table[T]) DeleteItemByKey(ctx context.Context, m Map, optFns ...func(*dynamodb.Options)) error {
	_, err := t.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: t.options.Schema.TableName(),
		Key:       m,
	}, optFns...)

	return err
}

func (t *Table[T]) Scan(ctx context.Context, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return func(yield func(ItemResult[T]) bool) {
		var lastEvaluatedKey map[string]types.AttributeValue

		for {
			scanInput := &dynamodb.ScanInput{
				TableName:                 t.options.Schema.TableName(),
				ConsistentRead:            aws.Bool(true),
				ExclusiveStartKey:         lastEvaluatedKey,
				Select:                    types.SelectAllAttributes,
				FilterExpression:          expr.Filter(),
				ProjectionExpression:      expr.Projection(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
			}

			res, err := t.client.Scan(ctx, scanInput, optFns...)
			if err != nil {
				if !yield(ItemResult[T]{err: err}) {
					return
				}

				return
			}

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

				if !yield(ItemResult[T]{item: *i}) {
					return
				}
			}

			lastEvaluatedKey = res.LastEvaluatedKey
			if lastEvaluatedKey == nil {
				return
			}
		}
	}
}

func (t *Table[T]) Query(ctx context.Context, expr expression.Expression, optFns ...func(*dynamodb.Options)) iter.Seq[ItemResult[T]] {
	return func(yield func(ItemResult[T]) bool) {
		var lastEvaluatedKey map[string]types.AttributeValue

		for {
			res, err := t.client.Query(ctx, &dynamodb.QueryInput{
				TableName:                 t.options.Schema.TableName(),
				ConsistentRead:            aws.Bool(true),
				ExclusiveStartKey:         lastEvaluatedKey,
				KeyConditionExpression:    expr.KeyCondition(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				FilterExpression:          expr.Filter(),
				ProjectionExpression:      expr.Projection(),
				Select:                    types.SelectAllAttributes,
			}, optFns...)
			if err != nil {
				if !yield(ItemResult[T]{err: err}) {
					return
				}

				return
			}

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

				if !yield(ItemResult[T]{item: *i}) {
					return
				}
			}

			lastEvaluatedKey = res.LastEvaluatedKey
			if lastEvaluatedKey == nil {
				return
			}
		}
	}
}
