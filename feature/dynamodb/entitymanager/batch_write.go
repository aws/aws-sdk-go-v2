package entitymanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// BatchWriteOperation provides a batched write (BatchWriteItem) operation for a DynamoDB table.
// It allows adding put and delete operations to a queue and executes them in batches.
type BatchWriteOperation[T any] struct {
	client Client
	table  *Table[T]
	schema *Schema[T]

	queue []batchOperation
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

	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationPut,
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

	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationPut,
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

	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationDelete,
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

	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationDelete,
		item: i,
	})

	return nil
}

// tableName returns the DynamoDB table name associated with this batch write operation.
func (b *BatchWriteOperation[T]) tableName() string {
	return *b.schema.TableName()
}

// queueItem returns the queued batch operation at the given offset, if any.
func (b *BatchWriteOperation[T]) queueItem(offset int) (batchOperation, bool) {
	if offset >= len(b.queue) {
		return batchOperation{}, false
	}

	return b.queue[offset], true
}

// fromMap satisfies the batcher interface for write operations. It returns nil
// because BatchWriteItem does not return items that need decoding.
func (b *BatchWriteOperation[T]) fromMap(_ map[string]types.AttributeValue) (any, error) {
	return nil, nil
}

// maxConsecutiveErrors returns the maximum number of allowed consecutive errors
// before the batch write executor stops processing requests.
func (b *BatchWriteOperation[T]) maxConsecutiveErrors() uint {
	return b.table.options.MaxConsecutiveErrors
}

// Merge creates a new BatchWriteExecutor that combines this batch write
// operation with additional batchers, allowing multiple tables or queues to be
// written in a single BatchWriteItem workflow.
func (b *BatchWriteOperation[T]) Merge(bs ...batcher) *BatchWriteExecutor[any] {
	return &BatchWriteExecutor[any]{
		client:   b.client,
		batchers: append([]batcher{b}, bs...),
	}
}

// Execute performs the batch write operation for all queued put and delete
// requests. It sends requests in batches of up to the maximum BatchWriteItem
// size and retries unprocessed items until they are written or the
// executor's maximum consecutive error threshold is reached.
func (b *BatchWriteOperation[T]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) error {
	executor := &BatchWriteExecutor[T]{
		client:   b.client,
		batchers: []batcher{b},
	}

	return executor.Execute(ctx, optFns...)
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
