package enhancedclient

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var _ Client = (*mockClient)(nil)

type mockClient struct {
	TableDescriptions map[string]types.TableDescription
	Items             []map[string]types.AttributeValue

	SetupFns []mockClientSetupFn
	Expects  []expectFn

	CreateTableCalls    []ddbCall[dynamodb.CreateTableInput, dynamodb.CreateTableOutput]
	DescribeTableCalls  []ddbCall[dynamodb.DescribeTableInput, dynamodb.DescribeTableOutput]
	DeleteTableCalls    []ddbCall[dynamodb.DeleteTableInput, dynamodb.DeleteTableOutput]
	GetItemCalls        []ddbCall[dynamodb.GetItemInput, dynamodb.GetItemOutput]
	PutItemCalls        []ddbCall[dynamodb.PutItemInput, dynamodb.PutItemOutput]
	DeleteItemCalls     []ddbCall[dynamodb.DeleteItemInput, dynamodb.DeleteItemOutput]
	UpdateItemCalls     []ddbCall[dynamodb.UpdateItemInput, dynamodb.UpdateItemOutput]
	BatchGetItemCalls   []ddbCall[dynamodb.BatchGetItemInput, dynamodb.BatchGetItemOutput]
	BatchWriteItemCalls []ddbCall[dynamodb.BatchWriteItemInput, dynamodb.BatchWriteItemOutput]
	ScanCalls           []ddbCall[dynamodb.ScanInput, dynamodb.ScanOutput]
	QueryCalls          []ddbCall[dynamodb.QueryInput, dynamodb.QueryOutput]
}

func newMockClient(fns ...mockClientSetupFn) *mockClient {
	out := &mockClient{}

	for _, fn := range fns {
		fn(out)
	}

	return out
}

type ddbCall[I, O any] func(Client, context.Context, *I, ...func(*dynamodb.Options)) (*O, error)

type ddbCallAsert[C, I, O any] func(*C, context.Context, *I, ...func(*dynamodb.Options)) (*O, error)

var _ ddbCallAsert[dynamodb.Client, dynamodb.CreateTableInput, dynamodb.CreateTableOutput] = (*dynamodb.Client).CreateTable
var _ ddbCallAsert[dynamodb.Client, dynamodb.DescribeTableInput, dynamodb.DescribeTableOutput] = (*dynamodb.Client).DescribeTable
var _ ddbCallAsert[dynamodb.Client, dynamodb.DeleteTableInput, dynamodb.DeleteTableOutput] = (*dynamodb.Client).DeleteTable
var _ ddbCallAsert[dynamodb.Client, dynamodb.GetItemInput, dynamodb.GetItemOutput] = (*dynamodb.Client).GetItem
var _ ddbCallAsert[dynamodb.Client, dynamodb.PutItemInput, dynamodb.PutItemOutput] = (*dynamodb.Client).PutItem
var _ ddbCallAsert[dynamodb.Client, dynamodb.UpdateItemInput, dynamodb.UpdateItemOutput] = (*dynamodb.Client).UpdateItem
var _ ddbCallAsert[dynamodb.Client, dynamodb.DeleteItemInput, dynamodb.DeleteItemOutput] = (*dynamodb.Client).DeleteItem
var _ ddbCallAsert[dynamodb.Client, dynamodb.QueryInput, dynamodb.QueryOutput] = (*dynamodb.Client).Query
var _ ddbCallAsert[dynamodb.Client, dynamodb.ScanInput, dynamodb.ScanOutput] = (*dynamodb.Client).Scan

func doDdbCall[I, O any](
	ctx context.Context,
	client *mockClient,
	callStack *[]ddbCall[I, O],
	input *I,
	optFns ...func(*dynamodb.Options),
) (*O, error) {
	if callStack == nil || len(*callStack) == 0 {
		return nil, fmt.Errorf(`unexpected call for %T`, callStack)
	}

	call := (*callStack)[0]
	*callStack = (*callStack)[1:]

	return call(client, ctx, input, optFns...)
}

func (c *mockClient) CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	return doDdbCall(ctx, c, &c.CreateTableCalls, input, optFns...)
}

func (c *mockClient) DescribeTable(ctx context.Context, input *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return doDdbCall(ctx, c, &c.DescribeTableCalls, input, optFns...)
}

func (c *mockClient) DeleteTable(ctx context.Context, input *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	return doDdbCall(ctx, c, &c.DeleteTableCalls, input, optFns...)
}

func (c *mockClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return doDdbCall(ctx, c, &c.GetItemCalls, input, optFns...)
}

func (c *mockClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return doDdbCall(ctx, c, &c.PutItemCalls, input, optFns...)
}

func (c *mockClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return doDdbCall(ctx, c, &c.DeleteItemCalls, input, optFns...)
}

func (c *mockClient) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return doDdbCall(ctx, c, &c.UpdateItemCalls, input, optFns...)
}

func (c *mockClient) BatchGetItem(ctx context.Context, input *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	return doDdbCall(ctx, c, &c.BatchGetItemCalls, input, optFns...)
}

func (c *mockClient) BatchWriteItem(ctx context.Context, input *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
	return doDdbCall(ctx, c, &c.BatchWriteItemCalls, input, optFns...)
}

func (c *mockClient) Scan(ctx context.Context, input *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return doDdbCall(ctx, c, &c.ScanCalls, input, optFns...)
}

func (c *mockClient) Query(ctx context.Context, input *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return doDdbCall(ctx, c, &c.QueryCalls, input, optFns...)
}

func (c *mockClient) RunExpectations(t *testing.T) {
	for _, fn := range c.Expects {
		if err := fn(t, c); err != nil {
			t.Errorf("expectation failed: %v", err)
		}
	}
}

type mockClientSetupFn func(*mockClient)

func withDefaultCreateTableCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.CreateTableCalls = append(m.CreateTableCalls, defaultCreateTableCall(m, err))
	}
}

func defaultCreateTableCall(client *mockClient, err error) ddbCall[dynamodb.CreateTableInput, dynamodb.CreateTableOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.CreateTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
		if err != nil {
			return nil, err
		}

		if client.TableDescriptions == nil {
			client.TableDescriptions = make(map[string]types.TableDescription)
		}

		tableName := aws.ToString(input.TableName)
		if desc, found := client.TableDescriptions[tableName]; found {
			return &dynamodb.CreateTableOutput{
				TableDescription: &desc,
			}, nil
		}

		desc := types.TableDescription{
			ArchivalSummary:      nil,
			AttributeDefinitions: input.AttributeDefinitions,
			BillingModeSummary: func() *types.BillingModeSummary {
				switch input.BillingMode {
				case types.BillingModePayPerRequest:
					return &types.BillingModeSummary{
						BillingMode:                       types.BillingModePayPerRequest,
						LastUpdateToPayPerRequestDateTime: aws.Time(time.Now()),
					}
				case types.BillingModeProvisioned:
					return &types.BillingModeSummary{
						BillingMode:                       types.BillingModeProvisioned,
						LastUpdateToPayPerRequestDateTime: aws.Time(time.Now()),
					}
				default:
					return nil
				}
			}(),
			CreationDateTime:          aws.Time(time.Now()),
			DeletionProtectionEnabled: input.DeletionProtectionEnabled,
			GlobalSecondaryIndexes: func() []types.GlobalSecondaryIndexDescription {
				var out []types.GlobalSecondaryIndexDescription
				for _, g := range input.GlobalSecondaryIndexes {
					out = append(out, types.GlobalSecondaryIndexDescription{
						Backfilling:        aws.Bool(false),
						IndexArn:           aws.String(fmt.Sprintf("arn:aws:dynamodb:eu-west-1:123456789012:table/%s/index/%s", tableName, *g.IndexName)),
						IndexName:          g.IndexName,
						IndexSizeBytes:     aws.Int64(0),
						IndexStatus:        types.IndexStatusActive,
						ItemCount:          aws.Int64(0),
						KeySchema:          g.KeySchema,
						OnDemandThroughput: g.OnDemandThroughput,
						Projection:         g.Projection,
						ProvisionedThroughput: func() *types.ProvisionedThroughputDescription {
							if g.ProvisionedThroughput == nil {
								return nil
							}

							return &types.ProvisionedThroughputDescription{
								LastDecreaseDateTime:   aws.Time(time.Now()),
								LastIncreaseDateTime:   aws.Time(time.Now()),
								NumberOfDecreasesToday: aws.Int64(0),
								ReadCapacityUnits:      g.ProvisionedThroughput.ReadCapacityUnits,
								WriteCapacityUnits:     g.ProvisionedThroughput.WriteCapacityUnits,
							}
						}(),
						WarmThroughput: func() *types.GlobalSecondaryIndexWarmThroughputDescription {
							if g.WarmThroughput == nil {
								return nil
							}

							return &types.GlobalSecondaryIndexWarmThroughputDescription{
								ReadUnitsPerSecond:  g.WarmThroughput.ReadUnitsPerSecond,
								WriteUnitsPerSecond: g.WarmThroughput.WriteUnitsPerSecond,
								Status:              types.IndexStatusActive,
							}
						}(),
					})
				}
				return out
			}(),
			GlobalTableVersion:   nil,
			GlobalTableWitnesses: nil,
			ItemCount:            aws.Int64(0),
			KeySchema:            input.KeySchema,
			LatestStreamArn:      nil,
			LatestStreamLabel:    nil,
			LocalSecondaryIndexes: func() []types.LocalSecondaryIndexDescription {
				var out []types.LocalSecondaryIndexDescription

				for _, l := range input.LocalSecondaryIndexes {
					out = append(out, types.LocalSecondaryIndexDescription{
						IndexArn:       aws.String(fmt.Sprintf("arn:aws:dynamodb:eu-west-1:123456789012:table/%s/index/%s", tableName, *l.IndexName)),
						IndexName:      l.IndexName,
						IndexSizeBytes: aws.Int64(0),
						ItemCount:      aws.Int64(0),
						KeySchema:      l.KeySchema,
						Projection:     l.Projection,
					})
				}

				return out
			}(),
			MultiRegionConsistency: types.MultiRegionConsistencyEventual,
			OnDemandThroughput:     input.OnDemandThroughput,
			ProvisionedThroughput: func() *types.ProvisionedThroughputDescription {
				if input.ProvisionedThroughput == nil {
					return nil
				}

				return &types.ProvisionedThroughputDescription{
					LastDecreaseDateTime:   aws.Time(time.Now()),
					LastIncreaseDateTime:   aws.Time(time.Now()),
					NumberOfDecreasesToday: aws.Int64(0),
					ReadCapacityUnits:      input.ProvisionedThroughput.ReadCapacityUnits,
					WriteCapacityUnits:     input.ProvisionedThroughput.WriteCapacityUnits,
				}
			}(),
			Replicas:            nil,
			RestoreSummary:      nil,
			SSEDescription:      nil,
			StreamSpecification: nil,
			TableArn:            aws.String(fmt.Sprintf("arn:aws:dynamodb:eu-west-1:123456789012:table/%s", tableName)),
			TableClassSummary:   nil,
			TableId:             nil,
			TableName:           input.TableName,
			TableSizeBytes:      aws.Int64(0),
			TableStatus:         types.TableStatusActive,
			WarmThroughput: func() *types.TableWarmThroughputDescription {
				if input.WarmThroughput == nil {
					return nil
				}

				return &types.TableWarmThroughputDescription{
					ReadUnitsPerSecond:  input.WarmThroughput.ReadUnitsPerSecond,
					WriteUnitsPerSecond: input.WarmThroughput.WriteUnitsPerSecond,
					Status:              types.TableStatusActive,
				}
			}(),
		}
		client.TableDescriptions[tableName] = desc

		return &dynamodb.CreateTableOutput{
			TableDescription: &desc,
		}, nil
	}
}

func withDefaultDescribeTableCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.DescribeTableCalls = append(m.DescribeTableCalls, defaultDescribeTableCall(m, err))
	}
}

func defaultDescribeTableCall(client *mockClient, err error) ddbCall[dynamodb.DescribeTableInput, dynamodb.DescribeTableOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
		if err != nil {
			return nil, err
		}

		desc, found := client.TableDescriptions[aws.ToString(input.TableName)]
		if !found {
			return nil, fmt.Errorf("table %q not found", aws.ToString(input.TableName))
		}

		return &dynamodb.DescribeTableOutput{
			Table: &desc,
		}, nil
	}
}

func withDefaultDeleteTableCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.DeleteTableCalls = append(m.DeleteTableCalls, defaultDeleteTableCall(m, err))
	}
}

func defaultDeleteTableCall(client *mockClient, err error) ddbCall[dynamodb.DeleteTableInput, dynamodb.DeleteTableOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
		if err != nil {
			return nil, err
		}

		desc, found := client.TableDescriptions[aws.ToString(input.TableName)]
		if !found {
			return nil, fmt.Errorf("table %q not found", aws.ToString(input.TableName))
		}

		delete(client.TableDescriptions, aws.ToString(input.TableName))

		return &dynamodb.DeleteTableOutput{
			TableDescription: &desc,
		}, nil
	}
}

func withDefaultGetItemCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.GetItemCalls = append(m.GetItemCalls, defaultGetItemCall(m, err))
	}
}

func defaultGetItemCall(client *mockClient, err error) ddbCall[dynamodb.GetItemInput, dynamodb.GetItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		if err != nil {
			return nil, err
		}

		if len(client.Items) == 0 {
			return &dynamodb.GetItemOutput{
				Item: nil,
			}, nil
		}

		item := client.Items[0]
		client.Items = client.Items[1:]

		return &dynamodb.GetItemOutput{
			Item: item,
		}, nil
	}
}

func withDefaultPutItemCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.PutItemCalls = append(m.PutItemCalls, defaultPutItemCall(m, err))
	}
}

func defaultPutItemCall(client *mockClient, err error) ddbCall[dynamodb.PutItemInput, dynamodb.PutItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		if err != nil {
			return nil, err
		}

		client.Items = append(client.Items, input.Item)

		return &dynamodb.PutItemOutput{
			Attributes: input.Item,
		}, nil
	}
}

func withDefaultDeleteItemCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.DeleteItemCalls = append(m.DeleteItemCalls, defaultDeleteItemCall(m, err))
	}
}

func defaultDeleteItemCall(client *mockClient, err error) ddbCall[dynamodb.DeleteItemInput, dynamodb.DeleteItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
		if err != nil {
			return nil, err
		}

		if len(client.Items) == 0 {
			return &dynamodb.DeleteItemOutput{
				Attributes: nil,
			}, nil
		}

		item := client.Items[0]
		client.Items = client.Items[1:]

		return &dynamodb.DeleteItemOutput{
			Attributes: item,
		}, err
	}
}

func withDefaultUpdateItemCall(err error) mockClientSetupFn {
	return func(m *mockClient) {
		m.UpdateItemCalls = append(m.UpdateItemCalls, defaultUpdateItemCall(m, err))
	}
}

func defaultUpdateItemCall(client *mockClient, err error) ddbCall[dynamodb.UpdateItemInput, dynamodb.UpdateItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
		if err != nil {
			return nil, err
		}

		item := map[string]types.AttributeValue{}
		maps.Copy(item, input.Key)

		for k, v := range input.ExpressionAttributeValues {
			if nk, found := input.ExpressionAttributeNames[k]; found {
				k = nk
			}

			item[k] = v
		}

		client.Items = append(client.Items, item)

		return &dynamodb.UpdateItemOutput{
			Attributes: item,
		}, nil
	}
}

func withDefaultBatchGetItemCall(err error, retCount uint, tableName string) mockClientSetupFn {
	return func(m *mockClient) {
		m.BatchGetItemCalls = append(m.BatchGetItemCalls, defaultBatchGetItemCall(m, err, retCount, tableName))
	}
}

func defaultBatchGetItemCall(client *mockClient, err error, retCount uint, tableName string) ddbCall[dynamodb.BatchGetItemInput, dynamodb.BatchGetItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
		if err != nil {
			return nil, err
		}

		if retCount == 0 {
			return &dynamodb.BatchGetItemOutput{}, nil
		}

		if len(client.Items) == 0 {
			return nil, errors.New("items have already been exhausted")
		}

		items := client.Items[0:retCount]
		out := &dynamodb.BatchGetItemOutput{
			Responses: map[string][]map[string]types.AttributeValue{
				tableName: items,
			},
		}

		client.Items = client.Items[len(items):]

		if len(client.Items) > 0 {
			out.UnprocessedKeys = map[string]types.KeysAndAttributes{
				tableName: {
					Keys: input.RequestItems[tableName].Keys[len(items):],
				},
			}
		}

		return out, nil
	}
}

func withDefaultBatchWriteItemCall(err error, retCount uint, tableName string) mockClientSetupFn {
	return func(m *mockClient) {
		m.BatchWriteItemCalls = append(m.BatchWriteItemCalls, defaultBatchWriteItemCall(m, err, retCount, tableName))
	}
}

func defaultBatchWriteItemCall(client *mockClient, err error, retCount uint, tableName string) ddbCall[dynamodb.BatchWriteItemInput, dynamodb.BatchWriteItemOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error) {
		if err != nil {
			return nil, err
		}

		if retCount == 0 {
			return &dynamodb.BatchWriteItemOutput{}, nil
		}

		if len(client.Items) == 0 {
			return nil, errors.New("items have already been exhausted")
		}

		items := input.RequestItems[tableName][:retCount]
		for _, i := range items {
			if i.PutRequest != nil {
				client.Items = append(client.Items, i.PutRequest.Item)
			}
			if i.DeleteRequest != nil {
				client.Items = client.Items[1:]
			}
		}
		out := &dynamodb.BatchWriteItemOutput{}

		if len(client.Items) > 0 {
			out.UnprocessedItems = map[string][]types.WriteRequest{
				tableName: input.RequestItems[tableName][retCount:],
			}
		}

		return out, nil
	}
}

func withDefaultScanCall(err error, retCount uint) mockClientSetupFn {
	return func(m *mockClient) {
		m.ScanCalls = append(m.ScanCalls, defaultScanCall(m, err, retCount))
	}
}

func defaultScanCall(client *mockClient, err error, retCount uint) ddbCall[dynamodb.ScanInput, dynamodb.ScanOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
		if err != nil {
			return nil, err
		}

		if retCount == 0 {
			return &dynamodb.ScanOutput{}, nil
		}

		if len(client.Items) == 0 {
			return nil, errors.New("items have already been exhausted")
		}

		items := client.Items[0:retCount]
		out := &dynamodb.ScanOutput{
			Items: items,
		}

		client.Items = client.Items[len(items):]

		if len(client.Items) > 0 {
			out.LastEvaluatedKey = items[len(items)-1]
		}

		return out, nil
	}
}

func withDefaultQueryCall(err error, retCount uint) mockClientSetupFn {
	return func(m *mockClient) {
		m.QueryCalls = append(m.QueryCalls, defaultQueryCall(m, err, retCount))
	}
}

func defaultQueryCall(client *mockClient, err error, retCount uint) ddbCall[dynamodb.QueryInput, dynamodb.QueryOutput] {
	return func(_ Client, _ context.Context, input *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		if err != nil {
			return nil, err
		}

		if retCount == 0 {
			return &dynamodb.QueryOutput{}, nil
		}

		if len(client.Items) == 0 {
			return nil, errors.New("items have already been exhausted")
		}

		items := client.Items[0:retCount]
		out := &dynamodb.QueryOutput{
			Items: items,
		}

		client.Items = client.Items[len(items):]

		if len(client.Items) > 0 {
			out.LastEvaluatedKey = items[len(items)-1]
		}

		return out, nil
	}
}

func withItem(item map[string]types.AttributeValue) mockClientSetupFn {
	return func(m *mockClient) {
		m.Items = append(m.Items, item)
	}
}

func withItems(generator func() map[string]types.AttributeValue, count uint) mockClientSetupFn {
	return func(m *mockClient) {
		for i := count; i > 0; i-- {
			m.Items = append(m.Items, generator())
		}
	}
}

type expectFn func(*testing.T, *mockClient) error

func withExpectFns(fn expectFn) mockClientSetupFn {
	return func(m *mockClient) {
		m.Expects = append(m.Expects, fn)
	}
}

func expectTablesCount(c uint) expectFn {
	return func(t *testing.T, m *mockClient) error {
		if len(m.TableDescriptions) != int(c) {
			return fmt.Errorf("expected %d tables, but found %d", c, len(m.TableDescriptions))
		}

		return nil
	}
}

func expectTable(tableName string) expectFn {
	return func(t *testing.T, m *mockClient) error {
		if _, found := m.TableDescriptions[tableName]; !found {
			return fmt.Errorf("expected table %q not found", tableName)
		}

		return nil
	}
}

func expectItemsCount(c uint) expectFn {
	return func(t *testing.T, m *mockClient) error {
		if len(m.Items) != int(c) {
			return fmt.Errorf("expected %d items, but found %d", c, len(m.Items))
		}

		return nil
	}
}
