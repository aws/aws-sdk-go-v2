package entitymanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	maxItemsInBatchWrite = 25 // maximum items allowed in a single BatchWriteItem call
)

// BatchWriteExecutor coordinates executing one or more batch write operations
// against DynamoDB, handling chunking, retries, and error thresholds.
type BatchWriteExecutor[R any] struct {
	client      Client
	batchers    []batcher
	fromMappers map[string]func(m map[string]types.AttributeValue) (any, error)
}

// Merge creates a new BatchWriteExecutor that combines the current batchers
// with additional ones, allowing multi-table batch write execution.
func (b *BatchWriteExecutor[R]) Merge(br ...batcher) *BatchWriteExecutor[any] {
	return &BatchWriteExecutor[any]{
		client:   b.client,
		batchers: append(b.batchers, br...),
	}
}

// Execute runs the batch write requests for all configured batchers. It sends
// requests in batches of up to maxItemsInBatchWrite items, and retries
// unprocessed items until they are written or the maximum consecutive error
// threshold is reached. Returns the last error encountered when the threshold
// is exceeded, or nil on success.
func (b *BatchWriteExecutor[R]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) error {
	// holds the starting point for each table
	batchersOffsets := map[string]int{}

	var consecutiveErrors uint = 0
	var maxConsecutiveErrors uint = 0

	if len(b.batchers) > 0 {
		maxConsecutiveErrors = b.batchers[0].maxConsecutiveErrors()
	}

	if maxConsecutiveErrors == 0 {
		maxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}

	remainder := make(map[string][]types.WriteRequest)

	for {
		bwii := &dynamodb.BatchWriteItemInput{
			RequestItems: remainder,
		}
		done := 0

		for _, items := range remainder {
			done += len(items)
		}

		for _, br := range b.batchers {
			for ; done < maxItemsInBatchWrite; done++ {
				offset := batchersOffsets[br.tableName()]
				if item, ok := br.queueItem(offset); ok {
					ri := bwii.RequestItems[br.tableName()]
					switch item.typ {
					case batchOperationPut:
						ri = append(ri, types.WriteRequest{
							PutRequest: &types.PutRequest{
								Item: item.item,
							},
						})
					case batchOperationDelete:
						ri = append(ri, types.WriteRequest{
							DeleteRequest: &types.DeleteRequest{
								Key: item.item,
							},
						})
					default:
						return fmt.Errorf(`unsupported operation type found: "%s"`, item.typ)
					}
					bwii.RequestItems[br.tableName()] = ri
				} else {
					break
				}

				batchersOffsets[br.tableName()] = offset + 1
			}
		}

		if done == 0 {
			break
		}

		res, err := b.client.BatchWriteItem(ctx, bwii, optFns...)
		if err != nil {
			consecutiveErrors++
			if consecutiveErrors >= maxConsecutiveErrors {
				return err
			}
		}

		consecutiveErrors = 0

		if res != nil && res.UnprocessedItems != nil {
			remainder = res.UnprocessedItems
		} else {
			remainder = make(map[string][]types.WriteRequest)
		}
	}

	return nil
}
