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
	// Size of the buffer in which it will be flushed.
	// Do not set to a number above 25 as DynamoDB rejects BatchWrites with
	// more than 25 items.
	FlushAmount   int
	tableName     string
	client        dynamodbiface.DynamoDBAPI
	primaryKeys   []string
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
		FlushAmount:   defaultFlushAmount,
		tableName:     tableName,
		client:        client,
		primaryKeys:   pKeys,
		requestBuffer: requestBuffer,
	}
	return batchWriter, nil
}

func (b *BatchWriter) flushIfNeeded() error {
	if len(b.requestBuffer) < b.FlushAmount {
		return nil
	}
	err := b.Flush()
	return err
}

func (b *BatchWriter) addWriteRequest(wr dynamodb.WriteRequest) error {
	b.requestBuffer = append(b.requestBuffer, wr)
	err := b.flushIfNeeded()
	return err
}

// PutItem adds a PutRequest operation to the requestBuffer.
func (b *BatchWriter) PutItem(putRequest *dynamodb.PutRequest) error {
	writeRequest := dynamodb.WriteRequest{PutRequest: putRequest}
	err := b.addWriteRequest(writeRequest)
	return err
}

// DeleteItem adds a DeleteRequest operation to the requestBuffer.
//
// The key argument should have the form of the Key argument the normal
// DeleteItem API call takes.
func (b *BatchWriter) DeleteItem(deleteRequest *dynamodb.DeleteRequest) error {
	writeRequest := dynamodb.WriteRequest{DeleteRequest: deleteRequest}
	err := b.addWriteRequest(writeRequest)
	return err
}

// Flush makes a BatchWriteItem with some or all requests that have been added
// so far to the buffer by means of PutItem and DeleteIem.
//
// Any unprocessed items sent in the response will be added to the end of the
// buffer, to be sent at a later date.
func (b *BatchWriter) Flush() error {
	if b.Empty() {
		return nil
	}
	flushBound := b.FlushAmount
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
