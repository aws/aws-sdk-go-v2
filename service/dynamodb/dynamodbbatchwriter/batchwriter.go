package dynamodbbatchwriter

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

	tableName        string
	client           dynamodbiface.DynamoDBAPI
	flushAmount      int
	requestBuffer    []dynamodb.WriteRequest
	sendRequestItems func(
		dynamodbiface.DynamoDBAPI, map[string][]dynamodb.WriteRequest,
	) (*dynamodb.BatchWriteItemOutput, error)
}

// New creates a new BatchWriter that will write to table `tableName`
// using `client`.
func New(tableName string, client dynamodbiface.DynamoDBAPI) *BatchWriter {
	requestBuffer := make(
		[]dynamodb.WriteRequest, 0, defaultRequestBufferCap,
	)
	return &BatchWriter{
		tableName:        tableName,
		client:           client,
		flushAmount:      defaultFlushAmount,
		requestBuffer:    requestBuffer,
		sendRequestItems: sendRequestItems,
	}
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
		output, err := b.sendRequestItems(b.client, requestItems)
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

// This function is logically a part of BatchWriter.Flush().
// I am pointing to it through BatchWriter.sendRequestItems so that I can
// implement tests for Flush() that don't actually make any requests.
// It is an ugly hack, but I don't see a clear alternative and Flush might be
// the most important method to test in BatchWriter.
// TODO: half-decent formatting.
func sendRequestItems(
	client dynamodbiface.DynamoDBAPI,
	requestItems map[string][]dynamodb.WriteRequest,
) (
	*dynamodb.BatchWriteItemOutput, error,
) {
	batchInput := dynamodb.BatchWriteItemInput{RequestItems: requestItems}
	batchRequest := client.BatchWriteItemRequest(&batchInput)
	output, err := batchRequest.Send()
	return output, err
}

// Empty returns whether or not the request buffer is empty.
func (b *BatchWriter) Empty() bool {
	return len(b.requestBuffer) == 0
}
