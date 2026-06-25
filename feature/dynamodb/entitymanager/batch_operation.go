package entitymanager

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type batchOperationType int

const (
	batchOperationGet = iota
	batchOperationPut
	batchOperationDelete
)

func (b batchOperationType) String() string {
	switch b {
	case batchOperationGet:
		return "batchOperationGet"
	case batchOperationPut:
		return "batchOperationPut"
	case batchOperationDelete:
		return "batchOperationDelete"
	default:
		return "unknown"
	}
}

type batchOperation struct {
	typ  batchOperationType
	item map[string]types.AttributeValue
}
