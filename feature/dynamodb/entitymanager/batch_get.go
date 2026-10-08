package entitymanager

import (
	"context"
	"fmt"
	"iter"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// BatchGetOperation provides a batched read (BatchGetItem) operation for a DynamoDB table.
// It allows adding items to a read queue and executes the batch operation, yielding results as an iterator.
type BatchGetOperation[T any] struct {
	client Client
	table  *Table[T]
	schema *Schema[T]

	queue []batchOperation
}

// AddReadItem adds an item to the batch read queue by extracting its key using the schema.
// The item must be a pointer to the struct type used by the table.
func (b *BatchGetOperation[T]) AddReadItem(item *T) error {
	m, err := b.schema.createKeyMap(item)
	if err != nil {
		return fmt.Errorf("error calling schema.createKeyMap: %w", err)
	}

	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationGet,
		item: m,
	})

	return nil
}

// AddReadItemByMap adds a key map directly to the batch read queue.
// The map should represent the primary key attributes for the table.
func (b *BatchGetOperation[T]) AddReadItemByMap(m Map) error {
	b.queue = append(b.queue, batchOperation{
		typ:  batchOperationGet,
		item: m,
	})

	return nil
}

// Execute performs the batch get operation for all queued items.
// It yields each result (or error) as an ItemResult[*T]. If the table name
// is not set, an error is yielded. Unprocessed keys are re-queued and
// retried until all are processed or the executor's error threshold is hit.
//
// Example usage:
//
//	seq := op.Execute(ctx)
//	for res := range iter.Chan(seq) { ... }
func (b *BatchGetOperation[T]) Execute(ctx context.Context, optFns ...func(options *dynamodb.Options)) iter.Seq[ItemResult[*T]] {
	executor := &BatchGetExecutor[*T]{
		client:   b.client,
		batchers: []batcher{b},
	}
	return executor.Execute(ctx, optFns...)
}

// tableName returns the DynamoDB table name associated with this batch get operation.
func (b *BatchGetOperation[T]) tableName() string {
	return *b.schema.TableName()
}

// queueItem returns the queued batch operation at the given offset, if any.
func (b *BatchGetOperation[T]) queueItem(offset int) (batchOperation, bool) {
	if offset >= len(b.queue) {
		return batchOperation{}, false
	}

	return b.queue[offset], true
}

// fromMap decodes a DynamoDB attribute map into a typed item and applies read extensions.
func (b *BatchGetOperation[T]) fromMap(m map[string]types.AttributeValue) (any, error) {
	i, err := b.schema.Decode(m)
	if err != nil {
		return nil, err
	}

	if err := b.table.applyAfterReadExtensions(i); err != nil {
		return nil, err
	}

	return i, err
}

// maxConsecutiveErrors returns the maximum number of allowed consecutive errors
// before the executor stops processing batch get results.
func (b *BatchGetOperation[T]) maxConsecutiveErrors() uint {
	return b.table.options.MaxConsecutiveErrors
}

// Merge creates a new BatchGetExecutor that combines this batch operation with
// additional batchers, allowing multiple tables or operations to be executed
// in a single BatchGetItem call.
func (b *BatchGetOperation[T]) Merge(bs ...batcher) *BatchGetExecutor[any] {
	return &BatchGetExecutor[any]{
		client:   b.client,
		batchers: append([]batcher{b}, bs...),
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
