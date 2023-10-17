package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

type mockBatchGetItemClient struct {
	inputs  []*BatchGetItemInput
	outputs []*BatchGetItemOutput
	t       *testing.T
}

func (c *mockBatchGetItemClient) BatchGetItem(ctx context.Context, input *BatchGetItemInput, optFns ...func(*Options)) (*BatchGetItemOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	if len(c.outputs) < requestCnt {
		c.t.Errorf("Paginator calls client more than expected %d times", len(c.outputs))
	}
	return c.outputs[requestCnt-1], nil
}

type batchGetItemTestCase struct {
	requestCnt             int
	stopOnDuplicationToken bool
	outputs                []*BatchGetItemOutput
}

func TestBatchGetItemPaginator(t *testing.T) {
	cases := map[string]batchGetItemTestCase{
		"total count 3": {
			requestCnt: 3,
			outputs: []*BatchGetItemOutput{
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{
						"key1": {
							AttributesToGet: []string{
								"attr1",
							},
							ConsistentRead:       aws.Bool(true),
							ProjectionExpression: aws.String("attr2, attr3"),
						},
					},
				},
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{
						"key2": {
							AttributesToGet: []string{
								"attr4",
							},
							ConsistentRead:       aws.Bool(true),
							ProjectionExpression: aws.String("attr5, attr6"),
						},
					},
				},
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{},
				},
			},
		},
		"total count 2 due to duplicate token": {
			requestCnt:             2,
			stopOnDuplicationToken: true,
			outputs: []*BatchGetItemOutput{
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{
						"key1": {
							AttributesToGet: []string{
								"attr1",
							},
							ConsistentRead:       aws.Bool(true),
							ProjectionExpression: aws.String("attr2, attr3"),
						},
					},
				},
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{
						"key1": {
							AttributesToGet: []string{
								"attr1",
							},
							ConsistentRead:       aws.Bool(true),
							ProjectionExpression: aws.String("attr2, attr3"),
						},
					},
				},
				{
					UnprocessedKeys: map[string]types.KeysAndAttributes{
						"key2": {
							AttributesToGet: []string{
								"attr4",
							},
							ConsistentRead:       aws.Bool(true),
							ProjectionExpression: aws.String("attr5, attr6"),
						},
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := mockBatchGetItemClient{
				t:       t,
				outputs: c.outputs,
				inputs:  []*BatchGetItemInput{},
			}
			paginator := NewBatchGetItemPaginator(&client, &BatchGetItemInput{}, func(options *BatchGetItemPaginatorOptions) {
				options.StopOnDuplicateToken = c.stopOnDuplicationToken
			})

			for paginator.HasMorePages() {
				_, err := paginator.NextPage(context.TODO())
				if err != nil {
					t.Errorf("error: %v", err)
				}
			}

			inputLen := len(client.inputs)
			testTotalRequests(c.requestCnt, inputLen, t)
			for i := 1; i < inputLen; i++ {
				if !awsutil.DeepEqual(client.inputs[i].RequestItems, c.outputs[i-1].UnprocessedKeys) {
					t.Errorf("Expect next input's request items to be eaqul to %v, got %v",
						c.outputs[i-1].UnprocessedKeys, client.inputs[i].RequestItems)
				}
			}
		})
	}
}

func testTotalRequests(expect, actual int, t *testing.T) {
	if actual != expect {
		t.Errorf("Expect total request number to be %d, got %d", expect, actual)
	}
}
