package batchwriter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

const testTableName = "testtable"
const hashKey = "numb"

// Global var holds cases that will be used for many tests.
var sharedCases = []struct {
	put, delete bool
	item        map[string]dynamodb.AttributeValue
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
func getBatchWriter() (*BatchWriter, *dynamodb.DynamoDB) {
	config := unit.Config()
	dynamoClient := dynamodb.New(config)
	dynamoClient.Handlers.Send.Clear()
	addResponse(dynamoClient, 200, `{
		"Table": {
			"KeySchema": [
				{
					"AttributeName": "`+hashKey+`",
					"KeyType": "HASH"
				}
			]
		}
	}`)
	tableName := testTableName

	batchWriter, _ := New(tableName, dynamoClient)
	// Clear the handler before returning the client.
	dynamoClient.Handlers.Send.Clear()
	return batchWriter, dynamoClient
}

// Convenience type alias
type itemmap map[string]interface{}

// Add a dummy response to a client.
func addResponse(client *dynamodb.DynamoDB, statusCode int, response string) {
	client.Handlers.Send.PushBack(func(r *aws.Request) {
		reader := ioutil.NopCloser(bytes.NewReader([]byte(response)))
		r.HTTPResponse = &http.Response{StatusCode: statusCode, Body: reader}
	})
}

func TestNewBatchWriter(t *testing.T) {
	batchWriter, _ := getBatchWriter()
	if batchWriter.tableName != testTableName {
		t.Errorf(`batchWriter.tableName set to "%s" when it should be "%s".`,
			batchWriter.tableName, testTableName)
	}
}

func TestNewWithPrimaryKeys(t *testing.T) {
	config := unit.Config()
	dynamoClient := dynamodb.New(config)
	dynamoClient.Handlers.Send.Clear()

	testPKeys := []string{"transaction", "n_orders"}
	batchWriter := NewWithPrimaryKeys(testTableName, dynamoClient, testPKeys)
	if batchWriter.tableName != testTableName {
		t.Errorf(`batchWriter.tableName set to "%s" when it should be "%s".`,
			batchWriter.tableName, testTableName)
	}

	if len(batchWriter.primaryKeys) != len(testPKeys) {
		t.Errorf(`batchWriter.primaryKeys has length %d, should have %d.`,
			len(batchWriter.primaryKeys), len(testPKeys))
	}
	for i := 0; i < len(batchWriter.primaryKeys); i++ {
		bwKey := batchWriter.primaryKeys[i]
		testKey := testPKeys[i]
		if bwKey != testKey {
			t.Errorf(`Primary key number %d set to %s. Should be %s.`,
				i, bwKey, testKey)
		}
	}
}

func TestNewError(t *testing.T) {
	config := unit.Config()
	dynamoClient := dynamodb.New(config)
	dynamoClient.Handlers.Send.Clear()
	addResponse(dynamoClient, 404, "ERR")
	_, err := New(testTableName, dynamoClient)
	if err == nil {
		t.Error("Did not propagate BatchWriter creation error correctly.")
	}
}

func TestPutOrDeleteItem(t *testing.T) {
	batchWriter, _ := getBatchWriter()

	cases := sharedCases
	// Make sure the flush amount is larger than the number of items to add.
	batchWriter.FlushAmount = len(cases) * 2
	for i := 0; i < len(cases); i++ {
		c := cases[i]
		if c.put {
			batchWriter.PutItem(&dynamodb.PutRequest{
				Item: c.item,
			})
		} else {
			batchWriter.DeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			})
		}
		bufferLen := len(batchWriter.requestBuffer)
		if bufferLen != (i + 1) {
			t.Errorf("Length of requestBuffer is %d when it should be %d.",
				len(cases), i+1)
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
	batchWriter, _ := getBatchWriter()
	if !batchWriter.Empty() { // BatchWriters should start empty.
		t.Error("batchWriter was initialized not empty.")
	}
	cases := sharedCases
	// flushAmount should be higher than the number of cases, so that we know
	// Empty() should return false.
	batchWriter.FlushAmount = len(cases) * 2
	for i, c := range cases {
		if c.put {
			batchWriter.PutItem(&dynamodb.PutRequest{
				Item: c.item,
			})
		} else {
			batchWriter.DeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			})
		}
		if batchWriter.Empty() {
			t.Errorf("Empty() returned a fase positive in iteration %d.", i)
		}
	}
}

func TestFlushError(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter()
	// Add a dummy response that should yield an error.
	addResponse(dynamoClient, 404, "ERR")
	cases := sharedCases
	for i, c := range cases {
		if c.put {
			batchWriter.PutItem(&dynamodb.PutRequest{
				Item: c.item,
			})
		} else {
			batchWriter.DeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			})
		}
		err := batchWriter.Flush()
		if err == nil {
			t.Errorf("Failed to return an error in iteration %d.", i)
		}
	}
}

func TestFlushUnprocessed(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter()
	// Number of unprocessed items to return.
	numUnpItems := 5
	// Generate a dummy body.
	body := map[string](map[string][]interface{}){
		"UnprocessedItems": map[string][]interface{}{
			testTableName: []interface{}{},
		},
	}
	// Add unprocessed items.
	for i := 0; i < numUnpItems; i++ {
		body["UnprocessedItems"][testTableName] = append(
			body["UnprocessedItems"][testTableName], map[string]interface{}{
				"PutRequest": map[string]interface{}{
					"Item": map[string]int{"a": 1, "numb": 89},
				},
			},
		)
	}
	bodyBytes, _ := json.Marshal(body)
	bodyStr := bytes.NewBuffer(bodyBytes).String()
	addResponse(dynamoClient, 200, bodyStr)
	cases := sharedCases
	// Make sure we won't flush while still adding items.
	batchWriter.FlushAmount = len(cases) * 2

	for _, c := range cases {
		if c.put {
			batchWriter.PutItem(&dynamodb.PutRequest{
				Item: c.item,
			})
		} else {
			batchWriter.DeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			})
		}
	}
	if len(batchWriter.requestBuffer) != len(cases) {
		t.Error("Wrong length for requestBuffer.")
	}
	batchWriter.Flush()
	// Note: this works because flushAmount is guaranteed to be higher than
	// the size of the requestBuffer. So all items will be flushed.
	expectedLength := numUnpItems
	if len(batchWriter.requestBuffer) != expectedLength {
		t.Errorf("Wrong number of items after flushing. Should be %d but is"+
			" %d.", expectedLength, len(batchWriter.requestBuffer))
	}
}

func TestFlushAutomatically(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter()
	dynamoClient.Handlers.Send.PushBack(func(r *aws.Request) {
		reader := ioutil.NopCloser(bytes.NewReader([]byte(
			`{}`,
		)))
		r.HTTPResponse = &http.Response{StatusCode: 200, Body: reader}
	})

	flushAmount := 5
	batchWriter.FlushAmount = flushAmount
	// We only want enough cases to make the BatchWriter flush automatically.
	cases := sharedCases[:flushAmount]
	for i, c := range cases {
		if c.put {
			batchWriter.PutItem(&dynamodb.PutRequest{
				Item: c.item,
			})
		} else {
			batchWriter.DeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			})
		}
		if i == flushAmount-1 {
			if !batchWriter.Empty() {
				t.Error(
					"BatchWriter not empty after reaching enough requests.")
			}
		} else if len(batchWriter.requestBuffer) != i+1 {
			t.Errorf("Wrong size for the requestBuffer on iteration %d.", i)
		}
	}
}
