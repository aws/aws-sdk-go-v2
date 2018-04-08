package batchwriter

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

const defaultFlushAmount = 25
const defaultRequestBufferCap = 50

// BatchWriter wraps a dynamodb client to expose a simple PutItem/DeleteItem
// API that buffers requests and takes advantage of BatchWriteItem behind
// the scenes.
type BatchWriter struct {
	_ struct{}

	tableName     string
	client        dynamodbiface.DynamoDBAPI
	primaryKeys   []string
	flushAmount   int
	requestBuffer []dynamodb.WriteRequest
}

// New creates a new BatchWriter that will write to table `tableName`
// using `client`.
//
// New will return an error if it fails to access the table information with a
// DescribeTableRequest.
func New(tableName string, client dynamodbiface.DynamoDBAPI) (*BatchWriter, error) {
	describeTableReq := client.DescribeTableRequest(&dynamodb.DescribeTableInput{
		TableName: &tableName,
	})
	describeTableOut, err := describeTableReq.Send()
	if err != nil {
		return &BatchWriter{}, err
	}
	// List of primary keys of the table. We will get them from a
	// DescribeTable request to DynamoDB.
	pKeys := []string{}
	for _, key := range describeTableOut.Table.KeySchema {
		pKeys = append(pKeys, *key.AttributeName)
	}
	requestBuffer := make(
		[]dynamodb.WriteRequest, 0, defaultRequestBufferCap,
	)
	batchWriter := &BatchWriter{
		tableName:     tableName,
		client:        client,
		primaryKeys:   pKeys,
		flushAmount:   defaultFlushAmount,
		requestBuffer: requestBuffer,
	}
	return batchWriter, nil
}

// SetFlushAmount changes the size at which the requestBuffer gets `Flush`ed
// automatically.
func (b *BatchWriter) SetFlushAmount(flushAmount int) {
	b.flushAmount = flushAmount
}

func (b *BatchWriter) flushIfNeeded() {
	if len(b.requestBuffer) >= b.flushAmount {
		b.Flush()
	}
}

func (b *BatchWriter) addWriteRequest(wr dynamodb.WriteRequest) {
	b.requestBuffer = append(b.requestBuffer, wr)
	b.flushIfNeeded()
}

// PutItem adds a PutRequest operation to the requestBuffer.
func (b *BatchWriter) PutItem(putRequest *dynamodb.PutRequest) {
	writeRequest := dynamodb.WriteRequest{PutRequest: putRequest}
	b.addWriteRequest(writeRequest)
}

// DeleteItem adds a DeleteRequest operation to the requestBuffer.
//
// The key argument should have the form of the Key argument the normal
// DeleteItem API call takes.
func (b *BatchWriter) DeleteItem(deleteRequest *dynamodb.DeleteRequest) {
	writeRequest := dynamodb.WriteRequest{DeleteRequest: deleteRequest}
	b.addWriteRequest(writeRequest)
}

// Flush makes a BatchWriteItem with some or all requests that have been added
// so far to the buffer by means of PutItem and DeleteIem.
//
// Any unprocessed items sent in the response will be added to the end of the
// buffer, to be sent at a later date.
func (b *BatchWriter) Flush() error {
	if len(b.requestBuffer) > 0 {
		flushBound := b.flushAmount
		if flushBound > len(b.requestBuffer) {
			flushBound = len(b.requestBuffer)
		}
		itemsToSend := b.requestBuffer[:flushBound]
		b.requestBuffer = b.requestBuffer[flushBound:]
		requestItems := make(map[string][]dynamodb.WriteRequest)
		requestItems[b.tableName] = itemsToSend
		output, err := b.sendRequestItems(requestItems)
		if err != nil {
			return err
		}
		unpItems, ok := output.UnprocessedItems[b.tableName]
		if ok {
			b.requestBuffer = append(b.requestBuffer, unpItems...)
		}
	}
	return nil
}

func (b *BatchWriter) sendRequestItems(
	requestItems map[string][]dynamodb.WriteRequest,
) (
	*dynamodb.BatchWriteItemOutput, error,
) {
	batchInput := dynamodb.BatchWriteItemInput{RequestItems: requestItems}
	batchRequest := b.client.BatchWriteItemRequest(&batchInput)
	output, err := batchRequest.Send()
	return output, err
}

// Empty returns whether or not the request buffer is empty.
func (b *BatchWriter) Empty() bool {
	return len(b.requestBuffer) == 0
}
