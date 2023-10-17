package kinesis

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type mockDescribeStreamClient struct {
	outputs []*DescribeStreamOutput
	inputs  []*DescribeStreamInput
	t       *testing.T
	limit   int32
}

func (c *mockDescribeStreamClient) DescribeStream(ctx context.Context, input *DescribeStreamInput, optFns ...func(*Options)) (*DescribeStreamOutput, error) {
	c.inputs = append(c.inputs, input)
	requestCnt := len(c.inputs)
	testCurRequest(len(c.outputs), requestCnt, c.limit, *input.Limit, c.t)
	return c.outputs[requestCnt-1], nil
}

type describeStreamTestCase struct {
	limit                  int32
	requestCnt             int
	stopOnDuplicationToken bool
	outputs                []*DescribeStreamOutput
}

func TestDescribeStreamPaginator(t *testing.T) {
	cases := map[string]describeStreamTestCase{
		"page limit 3": {
			limit:      3,
			requestCnt: 3,
			outputs: []*DescribeStreamOutput{
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard1"),
							},
							{
								ShardId: aws.String("shard2"),
							},
							{
								ShardId: aws.String("shard3"),
							},
						},
						HasMoreShards: aws.Bool(true),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard4"),
							},
							{
								ShardId: aws.String("shard5"),
							},
							{
								ShardId: aws.String("shard6"),
							},
						},
						HasMoreShards: aws.Bool(true),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard7"),
							},
						},
						HasMoreShards: aws.Bool(false),
					},
				},
			},
		},
		"total count 2 due to no more shards marker": {
			limit:      3,
			requestCnt: 2,
			outputs: []*DescribeStreamOutput{
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard1"),
							},
							{
								ShardId: aws.String("shard2"),
							},
							{
								ShardId: aws.String("shard3"),
							},
						},
						HasMoreShards: aws.Bool(true),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard4"),
							},
							{
								ShardId: aws.String("shard5"),
							},
							{
								ShardId: aws.String("shard6"),
							},
						},
						HasMoreShards: aws.Bool(false),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard7"),
							},
						},
						HasMoreShards: aws.Bool(false),
					},
				},
			},
		},
		"total count 2 due to duplicate shard ID": {
			limit:                  3,
			requestCnt:             2,
			stopOnDuplicationToken: true,
			outputs: []*DescribeStreamOutput{
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard1"),
							},
							{
								ShardId: aws.String("shard2"),
							},
							{
								ShardId: aws.String("shard3"),
							},
						},
						HasMoreShards: aws.Bool(true),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard4"),
							},
							{
								ShardId: aws.String("shard5"),
							},
							{
								ShardId: aws.String("shard3"),
							},
						},
						HasMoreShards: aws.Bool(true),
					},
				},
				{
					StreamDescription: &types.StreamDescription{
						Shards: []types.Shard{
							{
								ShardId: aws.String("shard7"),
							},
						},
						HasMoreShards: aws.Bool(false),
					},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := mockDescribeStreamClient{
				t:       t,
				outputs: c.outputs,
				inputs:  []*DescribeStreamInput{},
				limit:   c.limit,
			}
			paginator := NewDescribeStreamPaginator(&client, &DescribeStreamInput{}, func(options *DescribeStreamPaginatorOptions) {
				options.Limit = &c.limit
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
				shardsLength := len(c.outputs[i-1].StreamDescription.Shards)
				if *client.inputs[i].ExclusiveStartShardId != *c.outputs[i-1].StreamDescription.Shards[shardsLength-1].ShardId {
					t.Errorf("Expect next input's exclusive start shard ID to be eaqul to %s, got %s",
						*c.outputs[i-1].StreamDescription.Shards[shardsLength-1].ShardId, *client.inputs[i].ExclusiveStartShardId)
				}
			}
		})
	}
}

func testCurRequest(maxReqCnt, actualReqCnt int, expectLimit, actualLimit int32, t *testing.T) {
	if actualReqCnt > maxReqCnt {
		t.Errorf("Paginator calls client more than expected %d times", maxReqCnt)
	}
	if expectLimit != actualLimit {
		t.Errorf("Expect page limit to be %d, got %d", expectLimit, actualLimit)
	}
}

func testTotalRequests(expect, actual int, t *testing.T) {
	if actual != expect {
		t.Errorf("Expect total request number to be %d, got %d", expect, actual)
	}
}
