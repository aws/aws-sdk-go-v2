package enhancedclient

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (t *Table[T]) Create(ctx context.Context) (*dynamodb.CreateTableOutput, error) {
	input, err := t.options.Schema.createTableInput()
	if err != nil {
		return nil, err
	}

	return t.client.CreateTable(ctx, input, t.options.DynamoDBOptions...)
}

func (t *Table[T]) CreateWithWait(ctx context.Context, maxWaitDur time.Duration) error {
	cto, err := t.Create(ctx)
	if err != nil {
		return err
	}

	waiter := dynamodb.NewTableExistsWaiter(t.client)

	return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: cto.TableDescription.TableName,
	}, maxWaitDur)
}

func (t *Table[T]) Describe(ctx context.Context) (*dynamodb.DescribeTableOutput, error) {
	describe, err := t.options.Schema.describeTableInput()
	if err != nil {
		return nil, err
	}

	return t.client.DescribeTable(ctx, describe, t.options.DynamoDBOptions...)
}

func (t *Table[T]) Delete(ctx context.Context) (*dynamodb.DeleteTableOutput, error) {
	dlt, err := t.options.Schema.deleteTableInput()
	if err != nil {
		return nil, err
	}

	return t.client.DeleteTable(ctx, dlt, t.options.DynamoDBOptions...)
}

func (t *Table[T]) DeleteWithWait(ctx context.Context, maxWaitDur time.Duration) error {
	dlt, err := t.Delete(ctx)
	if err != nil {
		return err
	}

	waiter := dynamodb.NewTableNotExistsWaiter(t.client)

	return waiter.Wait(ctx, &dynamodb.DescribeTableInput{
		TableName: dlt.TableDescription.TableName,
	}, maxWaitDur)
}

func (t *Table[T]) Exists(ctx context.Context) (bool, error) {
	_, err := t.Describe(ctx)

	if err != nil {
		var notFound *types.ResourceNotFoundException

		if ok := errors.As(err, &notFound); ok {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
