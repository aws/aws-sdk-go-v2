package dynamodbbatchwriter

import (
	"errors"
	"reflect"
	"testing"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

const testTableName = "testtable"

// Global var holds cases that will be used for many tests.
var sharedCases = []struct{
	put, delete bool
	item map[string]dynamodb.AttributeValue
}{
	{true, false, marshal(itemmap{"itemcount": []int{89, 91, 92}, "key": "stf"})},
	{true, false, marshal(itemmap{"dance": 1})},
	{false, true, marshal(itemmap{"word": "dance"})},
	{true, false, marshal(itemmap{"func": itemmap{"in": 1, "out": 2}, "id": 142})},
	{false, true, marshal(itemmap{"tel": 123, "pos": 2})},
	{false, true, marshal(itemmap{"customer": 7776555})},
	{true, false, marshal(itemmap{
		"pd": "three", "func": itemmap{"in": 1, "out": 2, "us": 5}})},
}

// Convenience wrapper.
func marshal(in interface{}) map[string]dynamodb.AttributeValue {
	out, _ := dynamodbattribute.MarshalMap(in)
	return out
}

// Convenience wrapper.
func getBatchWriter() *BatchWriter {
	dynamoClient := &mockDynamoDBClient{}
	tableName := testTableName

	return New(tableName, dynamoClient)
}

// Convenience type alias
type itemmap map[string]interface{}

// Mock client to avoid running into client-generation errors.
type mockDynamoDBClient struct{
	dynamodbiface.DynamoDBAPI
}

func TestNewBatchWriter(t *testing.T) {
	dynamoClient := &mockDynamoDBClient{}
	tableName := testTableName

	batchWriter := New(tableName, dynamoClient)
	if !(reflect.TypeOf(batchWriter.client) == reflect.TypeOf(dynamoClient)) {
		t.Error("batchWriter.client is set incorrectly.")
	}
	if batchWriter.tableName != tableName {
		t.Errorf(`batchWriter.tableName set to "%v" when it should be "%s"`,
			batchWriter.tableName, tableName)
	}
}

func TestSetFlushAmount(t *testing.T) {
	batchWriter := getBatchWriter()

	testValues := []int{20, 150, 1, 12, 112}
	for val := range(testValues) {
		batchWriter.SetFlushAmount(val)
		if batchWriter.flushAmount != val {
			t.Errorf("batchWriter.flushAmount is set to %v instead of %d",
				batchWriter.flushAmount, val)
		}
	}
}

func TestPutOrDeleteItem(t *testing.T) {
	batchWriter := getBatchWriter()

	cases := sharedCases
	// Make sure the flush amount is larger than the number of items to add.
	batchWriter.SetFlushAmount(len(cases) * 2)
	for i := 0; i < len(cases); i++ {
		c := cases[i]
		if c.put {
			batchWriter.PutItem(c.item)
		} else {
			batchWriter.DeleteItem(c.item)
		}
		bufferLen := len(batchWriter.requestBuffer)
		if bufferLen != (i + 1) {
			t.Errorf("Length of requestBuffer is %d when it should be %d",
				len(cases), i + 1)
		}
		lastItem := batchWriter.requestBuffer[bufferLen-1]
		if c.put && lastItem.PutRequest == nil {
			t.Errorf(
				"Case no. %d has PutRequest == nil when it should be set.",
				i,
			)
		}
		if c.delete && lastItem.DeleteRequest == nil {
			t.Errorf(
				"Case no. %d has DeleteRequest == nil when it should be set.",
				i,
			)
		}
	}
}

func TestEmpty(t *testing.T) {
	batchWriter := getBatchWriter()
	if !batchWriter.Empty() { // BatchWriters should start empty.
		t.Error("batchWriter was initialized not empty.")
	}
	cases := sharedCases
	// flushAmount should be higher than the number of cases, so that we know
	// Empty() should return false.
	batchWriter.SetFlushAmount(len(cases) * 2)
	for i, c := range(cases) {
		if c.put {
			batchWriter.PutItem(c.item)
		} else {
			batchWriter.DeleteItem(c.item)
		}
		if batchWriter.Empty() {
			t.Errorf("Empty() returned a fase positive in iteration %d", i)
		}
	}
}

func TestFlushError(t *testing.T) {
	batchWriter := getBatchWriter()
	// Substitute function that requests the API for a dummy function that
	// returns an error
	batchWriter.sendRequestItems = errorOnItem
	cases := sharedCases
	for i, c := range(cases) {
		if c.put {
			batchWriter.PutItem(c.item)
		} else {
			batchWriter.DeleteItem(c.item)
		}
		err := batchWriter.Flush()
		if err == nil {
			t.Errorf("Failed to return an error in iteration %d", i)
		}
	}
}

func TestFlushUnprocessed(t *testing.T) {
	batchWriter := getBatchWriter()
	batchWriter.sendRequestItems = unprocessPutItems
	cases := sharedCases
	// Make sure we won't flush while still adding items.
	batchWriter.SetFlushAmount(len(cases) * 2)
	numPutItems := 0
	for _, c := range(cases) {
		if c.put {
			batchWriter.PutItem(c.item)
			// Count PutItems, because all of them will be
			// unprocessed.
			numPutItems++
		} else {
			batchWriter.DeleteItem(c.item)
		}
	}
	if len(batchWriter.requestBuffer) != len(cases) {
		t.Error("Wrong length for requestBuffer.")
	}
	batchWriter.Flush()
	// Note: this works because flushAmount is guaranteed to be higher than
	// the size of the requestBuffer. So all items will be flushed.
	expectedLength := numPutItems
	if len(batchWriter.requestBuffer) != expectedLength {
		t.Errorf("Wrong number of items after flushing. Should be %d but is" +
			" %d.", len(cases) - numPutItems, len(batchWriter.requestBuffer))
	}
}

func TestFlushAutomatically(t *testing.T) {
	batchWriter := getBatchWriter()
	batchWriter.sendRequestItems = dummyProcessItems

	flushAmount := 5
	batchWriter.SetFlushAmount(flushAmount)
	// We only want enough cases to make the BatchWriter flush automatically.
	cases := sharedCases[:flushAmount]
	for i, c := range(cases) {
		if c.put {
			batchWriter.PutItem(c.item)
		} else {
			batchWriter.DeleteItem(c.item)
		}
		if i == flushAmount - 1 {
			if !batchWriter.Empty() {
				t.Error(
					"BatchWriter not empty after reaching enough requests.")
			}
		} else if len(batchWriter.requestBuffer) != i + 1 {
			t.Errorf("Wrong size for the requestBuffer on iteration %d", i)
		}
	}
}

// Dummy functions that substitute sendRequestItems.

// Return an error regardless of the input.
func errorOnItem(
	client dynamodbiface.DynamoDBAPI,
	requestItems map[string][]dynamodb.WriteRequest,
) (
	*dynamodb.BatchWriteItemOutput, error,
) {
	return &dynamodb.BatchWriteItemOutput{}, errors.New("Error")
}

// Send all PutRequests back as UnprocessedItems.
func unprocessPutItems(
	client dynamodbiface.DynamoDBAPI,
	requestItems map[string][]dynamodb.WriteRequest,
) (
	*dynamodb.BatchWriteItemOutput, error,
) {
	unpItems := make([]dynamodb.WriteRequest, 0, 10)
	for _, req := range(requestItems[testTableName]) {
		if req.PutRequest != nil {
			item := req.PutRequest.Item
			putReq := &dynamodb.PutRequest{Item: item}
			writeReq := dynamodb.WriteRequest{PutRequest: putReq}
			unpItems = append(unpItems, writeReq)
		}
	}
	output := &dynamodb.BatchWriteItemOutput{
		UnprocessedItems: map[string][]dynamodb.WriteRequest{
			testTableName: unpItems,
		},
	}
	return output, nil
}

// Return an empty BatchWriteItemOutput and no error.
func dummyProcessItems(
	client dynamodbiface.DynamoDBAPI,
	requestItems map[string][]dynamodb.WriteRequest,
) (
	*dynamodb.BatchWriteItemOutput, error,
) {
	return &dynamodb.BatchWriteItemOutput{}, nil
}
