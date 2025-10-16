package enhancedclient

import (
	"context"
	"fmt"

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

	err = t.options.Schema.applyAfterReadExtensions(item)
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

	err = t.options.Schema.applyAfterReadExtensions(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// PutItem writes the item without checking for collisions
func (t *Table[T]) PutItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) (*T, error) {
	err := t.options.Schema.applyBeforeWriteExtensions(item)
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

	return t.options.Schema.Decode(itemMap)
}

// UpdateItem writes the item with additional checks (version checks, etc)
func (t *Table[T]) UpdateItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) (*T, error) {
	err := t.options.Schema.applyBeforeWriteExtensions(item)
	if err != nil {
		return nil, err
	}

	m, err := t.options.Schema.createKeyMap(item)
	if err != nil {
		return nil, err
	}

	expr, err := t.options.Schema.createUpdateExpression(item)
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

	return t.options.Schema.Decode(res.Attributes)
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

func (t *Table[T]) Scan(ctx context.Context, expr expression.Expression, optFns ...func(*dynamodb.Options)) ([]*T, error) {
	var exclusiveStartKey map[string]types.AttributeValue

	out := make([]*T, 0)
	done := false

	for !done {
		res, err := t.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                 t.options.Schema.TableName(),
			ConsistentRead:            aws.Bool(true),
			ExclusiveStartKey:         exclusiveStartKey,
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
			Select:                    types.SelectAllAttributes,
			//ReturnConsumedCapacity:    "",
			//TotalSegments:          nil,
			//Segment:                   nil,
			//Limit:                     aws.Int32(5),
		}, optFns...)

		if err != nil {
			return nil, err
		}

		exclusiveStartKey = res.LastEvaluatedKey
		done = len(exclusiveStartKey) == 0

		for idx := range res.Items {
			item, err := t.options.Schema.Decode(res.Items[idx])
			if err != nil {
				return nil, err
			}

			err = t.options.Schema.applyAfterReadExtensions(item)
			if err != nil {
				return nil, err
			}

			out = append(out, item)
		}
	}

	return out, nil
}
