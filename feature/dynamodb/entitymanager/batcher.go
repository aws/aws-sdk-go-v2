package entitymanager

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

// batcher is an internal interface implemented by batch get and batch write
// operations so their executors can treat them uniformly. It exposes access to
// queued operations, the target table name, error thresholds, and optional
// mapping from raw attribute maps to typed items.
type batcher interface {
	// queueItem returns the queued batch operation at the given offset, if any.
	queueItem(int) (batchOperation, bool)
	// tableName returns the DynamoDB table name associated with this batcher.
	tableName() string
	// maxConsecutiveErrors returns the maximum number of consecutive errors
	// allowed before the executor stops processing.
	maxConsecutiveErrors() uint
	// fromMap converts a DynamoDB attribute map into a typed item, when
	// applicable (read operations). For write operations it may be a no-op.
	fromMap(m map[string]types.AttributeValue) (any, error)
}
