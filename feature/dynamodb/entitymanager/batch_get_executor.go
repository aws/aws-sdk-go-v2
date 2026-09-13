package entitymanager

import (
	"context"
	"fmt"
	"iter"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	maxItemsInBatchGet = 25
)

// BatchGetExecutor coordinates executing one or more batch get operations
// against DynamoDB, handling chunking, retries, and result decoding.
type BatchGetExecutor[R any] struct {
	client      Client
	batchers    []batcher
	fromMappers map[string]func(m map[string]types.AttributeValue) (any, error)
}

// fromMap finds the appropriate mapper for the given table and decodes the
// attribute map into a typed item.
func (b *BatchGetExecutor[R]) fromMap(tableName string, m map[string]types.AttributeValue) (any, error) {
	if b.fromMappers == nil {
		b.fromMappers = make(map[string]func(m map[string]types.AttributeValue) (any, error))
	}

	if r, ok := b.fromMappers[tableName]; ok {
		return r(m)
	}

	for _, br := range b.batchers {
		if br.tableName() == tableName {
			b.fromMappers[tableName] = br.fromMap

			break
		}
	}

	if r, ok := b.fromMappers[tableName]; ok {
		return r(m)
	} else {
		return nil, fmt.Errorf(`unable to find fromMapper() for table "%s"`, tableName)
	}
}

// Merge creates a new BatchGetExecutor that combines the current batchers with
// additional ones, enabling multi-table batch get execution.
func (b *BatchGetExecutor[R]) Merge(br ...batcher) *BatchGetExecutor[any] {
	return &BatchGetExecutor[any]{
		client:   b.client,
		batchers: append(b.batchers, br...),
	}
}

// Execute runs the batch get requests for all configured batchers, yielding
// each decoded item or error as an ItemResult[R]. It retries unprocessed keys
// until they are exhausted or the maximum consecutive error threshold is
// reached.
func (b *BatchGetExecutor[R]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) iter.Seq[ItemResult[R]] {
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

	return func(yield func(ItemResult[R]) bool) {
		remainder := map[string]types.KeysAndAttributes{}

		for {
			bgii := &dynamodb.BatchGetItemInput{
				RequestItems: remainder,
			}

			done := 0
			for _, items := range remainder {
				done += len(items.Keys)
			}

			for _, br := range b.batchers {
				for ; done < maxItemsInBatchGet; done++ {
					offset := batchersOffsets[br.tableName()]
					if item, ok := br.queueItem(offset); ok {
						ri := bgii.RequestItems[br.tableName()]
						ri.Keys = append(ri.Keys, item.item)
						bgii.RequestItems[br.tableName()] = ri
					} else {
						break
					}

					batchersOffsets[br.tableName()] = offset + 1
				}
			}

			if done == 0 {
				return
			}

			res, err := b.client.BatchGetItem(ctx, bgii, optFns...)
			if err != nil {
				if !yield(ItemResult[R]{err: err}) {
					return
				}

				consecutiveErrors++
				if consecutiveErrors >= maxConsecutiveErrors {
					return
				}

				continue
			}

			consecutiveErrors = 0

			for tableName, items := range res.Responses {
				for _, i := range items {
					item, err := b.fromMap(tableName, i)
					if !yield(ItemResult[R]{item: item.(R), table: tableName, err: err}) {
						return
					}
				}
			}

			if res != nil && res.UnprocessedKeys != nil {
				remainder = res.UnprocessedKeys
			} else {
				remainder = make(map[string]types.KeysAndAttributes)
			}
		}
	}
}
