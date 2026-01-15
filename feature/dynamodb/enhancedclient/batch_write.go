package enhancedclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type batchWriteOpType int

const (
	batchWriteOpPut    batchWriteOpType = 0
	batchWriteOpDelete batchWriteOpType = 1

	maxItemsInBatchWrite = 25
)

type batchWriteQueueItem struct {
	typ  batchWriteOpType
	item map[string]types.AttributeValue
}

// BatchWriteOperation provides a batched write (BatchWriteItem) operation for a DynamoDB table.
// It allows adding put and delete operations to a queue and executes them in batches.
type BatchWriteOperation[T any] struct {
	client Client
	table  *Table[T]
	schema *Schema[T]

	queue []batchWriteQueueItem
}

// AddPut adds a put (insert/update) operation to the batch write queue.
// The item is encoded using the table's schema and extensions are applied before writing.
func (b *BatchWriteOperation[T]) AddPut(item *T) error {
	if err := b.table.applyBeforeWriteExtensions(item); err != nil {
		return fmt.Errorf("error calling table.applyBeforeWriteExtensions(): %w", err)
	}

	m, err := b.schema.Encode(item)
	if err != nil {
		return fmt.Errorf("error calling schema.Encode(): %w", err)
	}

	b.queue = append(b.queue, batchWriteQueueItem{
		typ:  batchWriteOpPut,
		item: m,
	})

	return nil
}

// AddRawPut adds a put operation to the batch write queue using a raw attribute value map.
// The map should represent the full item to be written.
func (b *BatchWriteOperation[T]) AddRawPut(i map[string]types.AttributeValue) error {
	if len(i) == 0 {
		return fmt.Errorf("input map is empty")
	}

	b.queue = append(b.queue, batchWriteQueueItem{
		typ:  batchWriteOpPut,
		item: i,
	})

	return nil
}

// AddDelete adds a delete operation to the batch write queue for the given item.
// The item's key is extracted using the schema.
func (b *BatchWriteOperation[T]) AddDelete(item *T) error {
	m, err := b.schema.createKeyMap(item)
	if err != nil {
		return fmt.Errorf("error calling schema.createKeyMap(): %w", err)
	}

	b.queue = append(b.queue, batchWriteQueueItem{
		typ:  batchWriteOpDelete,
		item: m,
	})

	return nil
}

// AddRawDelete adds a delete operation to the batch write queue using a raw key map.
// The map should represent the primary key attributes of the item to delete.
func (b *BatchWriteOperation[T]) AddRawDelete(i map[string]types.AttributeValue) error {
	if len(i) == 0 {
		return fmt.Errorf("input map is empty")
	}

	b.queue = append(b.queue, batchWriteQueueItem{
		typ:  batchWriteOpDelete,
		item: i,
	})

	return nil
}

// Execute performs the batch write operation for all queued put and delete requests.
// It sends requests in batches of up to 25 items, and retries unprocessed items until all are written.
// Returns an error if the table name is not set or if a request fails.
func (b *BatchWriteOperation[T]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) error {
	var consecutiveErrors uint = 0
	var maxConsecutiveErrors = b.table.options.MaxConsecutiveErrors
	if maxConsecutiveErrors == 0 {
		maxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}

	tableName := b.schema.TableName()
	if tableName == nil {
		return fmt.Errorf("empty table name, did you forget Schema[T].WithName()?")
	}

	for len(b.queue) > 0 {
		pos := min(maxItemsInBatchWrite, len(b.queue))
		batch := make([]types.WriteRequest, 0, pos)

		ops := b.queue[0:pos]
		b.queue = b.queue[pos:]

		for _, op := range ops {
			switch op.typ {
			case batchWriteOpPut:
				batch = append(batch, types.WriteRequest{
					PutRequest: &types.PutRequest{
						Item: op.item,
					},
				})
			case batchWriteOpDelete:
				batch = append(batch, types.WriteRequest{
					DeleteRequest: &types.DeleteRequest{
						Key: op.item,
					},
				})
			default:
				return fmt.Errorf("unknown batchWriteOpType: %d", op.typ)
			}
		}

		bwii := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				*tableName: batch,
			},
		}

		res, err := b.client.BatchWriteItem(ctx, bwii, optFns...)
		if err != nil {
			consecutiveErrors++

			if consecutiveErrors >= maxConsecutiveErrors {
				return err
			}
		}

		var unprocessedItems []types.WriteRequest
		if res != nil && res.UnprocessedItems != nil {
			unprocessedItems = res.UnprocessedItems[*tableName]
		} else if err != nil {
			unprocessedItems = bwii.RequestItems[*tableName]
		}

		for _, ui := range unprocessedItems {
			if ui.PutRequest != nil {
				b.queue = append(b.queue, batchWriteQueueItem{
					typ:  batchWriteOpPut,
					item: ui.PutRequest.Item,
				})
			} else if ui.DeleteRequest != nil {
				b.queue = append(b.queue, batchWriteQueueItem{
					typ:  batchWriteOpDelete,
					item: ui.DeleteRequest.Key,
				})
			}
		}
	}

	return nil
}

// NewBatchWriteOperation creates a new BatchWriteOperation for the given table.
// Use this to perform batched put and delete operations (BatchWriteItem) for the table's items.
func NewBatchWriteOperation[T any](table *Table[T]) *BatchWriteOperation[T] {
	return &BatchWriteOperation[T]{
		client: table.client,
		table:  table,
		schema: table.options.Schema,
	}
}
