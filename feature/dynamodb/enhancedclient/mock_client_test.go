package enhancedclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ Client = (*mockClient)(nil)

type mockClient struct {
	Items  []map[string]types.AttributeValue
	Errors map[int]error
	//CreateTableCalls   int //[]func(context.Context, *dynamodb.CreateTableInput, ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	//DescribeTableCalls int //[]func(context.Context, *dynamodb.DescribeTableInput, ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	//DeleteTableCalls   int //[]func(context.Context, *dynamodb.DeleteTableInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
	//GetItemCalls       int //[]func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	//PutItemCalls       int //[]func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	//UpdateItemCalls    int //[]func(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	//DeleteItemCalls    int //[]func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	//BatchGetItemCalls  int //[]func(context.Context, *dynamodb.BatchGetItemInput, ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
	//ScanCalls          int //[]func(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	//QueryCalls         int //[]func(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

func (c *mockClient) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	//if len(c.CreateTableCalls) == 0 {
	//	panic("unexpected CreateTable() call")
	//}
	//
	//f := c.CreateTableCalls[0]
	//
	//c.CreateTableCalls = c.CreateTableCalls[1:]
	//
	//return f(ctx, input, optFns...)
	return nil, nil
}

func (c *mockClient) DescribeTable(ctx context.Context, input *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	//if len(c.DescribeTableCalls) == 0 {
	//	panic("unexpected DescribeTable() call")
	//}
	//
	//f := c.DescribeTableCalls[0]
	//
	//c.DescribeTableCalls = c.DescribeTableCalls[1:]
	//
	//return f(ctx, input, optFns...)
	return nil, nil
}

func (c *mockClient) DeleteTable(ctx context.Context, input *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	//if len(c.DeleteTableCalls) == 0 {
	//	panic("unexpected DeleteTable() call")
	//}
	//
	//f := c.DeleteTableCalls[0]
	//
	//c.DeleteTableCalls = c.DeleteTableCalls[1:]
	//
	//return f(ctx, input, optFns...)
	return nil, nil
}

func (c *mockClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	defer func() {
		if len(c.Items) > 0 {
			c.Items = c.Items[1:]
		}
	}()
	if len(c.Items) == 0 {
		panic("unexpected GetItem() call")
	}

	return &dynamodb.GetItemOutput{
		Item: c.Items[0],
	}, c.Errors[0]
}

func (c *mockClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	defer func() {
		if len(c.Items) > 0 {
			c.Items = c.Items[1:]
		}
	}()
	if len(c.Items) == 0 {
		panic("unexpected UpdateItem() call")
	}

	return &dynamodb.PutItemOutput{
		Attributes: c.Items[0],
	}, nil
}

func (c *mockClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return &dynamodb.DeleteItemOutput{}, nil
}

func (c *mockClient) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	defer func() {
		if len(c.Items) > 0 {
			c.Items = c.Items[1:]
		}
	}()
	if len(c.Items) == 0 {
		panic("unexpected UpdateItem() call")
	}

	return &dynamodb.UpdateItemOutput{
		Attributes: c.Items[0],
	}, nil
}

func (c *mockClient) BatchGetItem(ctx context.Context, input *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	//if len(c.BatchGetItemCalls) == 0 {
	//	panic("unexpected BatchGetItem() call")
	//}
	//
	//f := c.BatchGetItemCalls[0]
	//
	//c.BatchGetItemCalls = c.BatchGetItemCalls[1:]
	//
	//return f(ctx, input, optFns...)
	return nil, nil
}

func (c *mockClient) Scan(ctx context.Context, input *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return &dynamodb.ScanOutput{
		Items: c.Items,
	}, nil
}

func (c *mockClient) Query(ctx context.Context, input *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return &dynamodb.QueryOutput{
		Items: c.Items,
	}, nil
}
