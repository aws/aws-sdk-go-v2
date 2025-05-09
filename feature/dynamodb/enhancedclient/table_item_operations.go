package enhancedclient

import (
	"context"
	"fmt"
	"strings"

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

	return t.options.Schema.Decode(res.Item)
}

func (t *Table[T]) PutItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) error {
	itemMap, err := t.options.Schema.Encode(item)
	if err != nil {
		return err
	}

	res, err := t.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: t.options.Schema.TableName(),
		Item:      itemMap,
	}, optFns...)
	if err != nil {
		return err
	}
	if res == nil {
		return fmt.Errorf("empty response in PutItem() call")
	}
	return nil
}

func (t *Table[T]) UpdateItem(ctx context.Context, item *T, optFns ...func(*dynamodb.Options)) error {
	m, err := t.options.Schema.createKeyMap(item)
	if err != nil {
		return err
	}

	itemMap, err := t.options.Schema.Encode(item)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(itemMap))
	values := map[string]types.AttributeValue{}
	for k, v := range itemMap {
		keys = append(keys, fmt.Sprintf("%s = :%s", k, k))
		values[fmt.Sprintf(":%s", k)] = v
	}

	_, err = t.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       m,
		TableName:                 t.options.Schema.TableName(),
		UpdateExpression:          pointer(fmt.Sprintf("SET %s", strings.Join(keys, ", "))),
		ExpressionAttributeValues: values,
	}, optFns...)

	return err
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
