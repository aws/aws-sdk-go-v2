package enhancedclient

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

type batchReadQueueItem struct {
	item map[string]types.AttributeValue
}

// BatchGetOperation provides a batched read (BatchGetItem) operation for a DynamoDB table.
// It allows adding items to a read queue and executes the batch operation, yielding results as an iterator.
type BatchGetOperation[T any] struct {
	client Client
	table  *Table[T]
	schema *Schema[T]

	queue []batchReadQueueItem
}

// AddReadItem adds an item to the batch read queue by extracting its key using the schema.
// The item must be a pointer to the struct type used by the table.
func (b *BatchGetOperation[T]) AddReadItem(item *T) error {
	m, err := b.schema.createKeyMap(item)
	if err != nil {
		return fmt.Errorf("error calling schema.createKeyMap: %w", err)
	}

	b.queue = append(b.queue, batchReadQueueItem{
		item: m,
	})

	return nil
}

// AddReadItemByMap adds a key map directly to the batch read queue.
// The map should represent the primary key attributes for the table.
func (b *BatchGetOperation[T]) AddReadItemByMap(m Map) error {
	b.queue = append(b.queue, batchReadQueueItem{
		item: m,
	})

	return nil
}

// Execute performs the batch get operation for all queued items.
// It yields each result (or error) using the provided iterator pattern.
// If the table name is not set, an error is yielded.
// Unprocessed keys are re-queued and retried until all are processed.
//
// Example usage:
//
//	seq := op.Execute(ctx)
//	for res := range iter.Chan(seq) { ... }
func (b *BatchGetOperation[T]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) iter.Seq[ItemResult[T]] {
	var consecutiveErrors uint = 0
	var maxConsecutiveErrors = b.table.options.MaxConsecutiveErrors
	if maxConsecutiveErrors == 0 {
		maxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}

	return func(yield func(ItemResult[T]) bool) {
		tableName := b.schema.TableName()
		if tableName == nil {
			yield(ItemResult[T]{err: fmt.Errorf("empty table name, did you forget Schema[T].WithName()?")})
			return
		}

		for len(b.queue) > 0 {
			pos := min(len(b.queue), maxItemsInBatchGet)
			items := b.queue[0:pos]
			keys := make([]map[string]types.AttributeValue, 0, pos)
			for _, item := range items {
				keys = append(keys, item.item)
			}
			b.queue = b.queue[pos:]
			bgii := &dynamodb.BatchGetItemInput{
				RequestItems: map[string]types.KeysAndAttributes{
					*tableName: {
						Keys: keys,
					},
				},
			}

			res, err := b.client.BatchGetItem(ctx, bgii, optFns...)
			if err != nil {
				if !yield(ItemResult[T]{err: err}) {
					return
				}

				if consecutiveErrors >= maxConsecutiveErrors {
					return
				}

				continue
			}
			consecutiveErrors = 0

			if res != nil && res.Responses != nil {
				for _, item := range res.Responses[*tableName] {
					i, err := b.schema.Decode(item)
					if err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if err := b.table.applyAfterReadExtensions(i); err != nil {
						if !yield(ItemResult[T]{err: err}) {
							return
						}

						continue
					}

					if !yield(ItemResult[T]{item: i}) {
						return
					}
				}
			}

			if res != nil && res.UnprocessedKeys != nil {
				for _, item := range res.UnprocessedKeys[*tableName].Keys {
					b.queue = append(b.queue, batchReadQueueItem{
						item: item,
					})
				}
			}
		}
	}
}

// NewBatchGetOperation creates a new BatchGetOperation for the given table.
// Use this to perform batched reads (BatchGetItem) for the table's items.
func NewBatchGetOperation[T any](table *Table[T]) *BatchGetOperation[T] {
	return &BatchGetOperation[T]{
		client: table.client,
		table:  table,
		schema: table.options.Schema,
	}
}
