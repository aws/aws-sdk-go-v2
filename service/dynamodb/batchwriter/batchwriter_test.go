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
const (
	hashKey = "word"
	sortKey = "number"
)

type testCases []struct {
	put, delete bool
	item        map[string]dynamodb.AttributeValue
}

// Global var holds cases that will be used for many tests.
var sharedCases = testCases{
	{true, false, marshal(itemmap{
		hashKey: "far", sortKey: 0,
		"itemcount": []int{89, 91, 92}, "key": "stf"},
	)},
	{true, false, marshal(itemmap{hashKey: "near", sortKey: 0, "dance": 1})},
	{false, true, marshal(itemmap{hashKey: "dance", sortKey: 2})},
	{true, false, marshal(itemmap{
		hashKey: "far", sortKey: 1,
		"func": itemmap{"in": 1, "out": 2}, "id": 142},
	)},
	{false, true, marshal(itemmap{hashKey: "1", sortKey: 9})},
	{false, true, marshal(itemmap{hashKey: "9", sortKey: 1})},
	{true, false, marshal(itemmap{
		hashKey: "me", sortKey: 55555,
		"pd": "three", "func": itemmap{"in": 1, "out": 2, "us": 5}})},
}

// Convenience wrapper.
func marshal(in interface{}) map[string]dynamodb.AttributeValue {
	out, _ := dynamodbattribute.MarshalMap(in)
	return out
}

// Convenience wrapper.
func getBatchWriter(addSortKey bool) (*BatchWriter, *dynamodb.DynamoDB) {
	config := unit.Config()
	dynamoClient := dynamodb.New(config)
	dynamoClient.Handlers.Send.Clear()
	if addSortKey {
		addResponse(dynamoClient, 200, `{
			"Table": {
				"KeySchema": [
					{
						"AttributeName": "`+hashKey+`",
						"KeyType": "HASH"
					},
					{
						"AttributeName": "`+sortKey+`",
						"KeyType": "SORT"
					}
				]
			}
		}`)
	} else {
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

	}
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
	batchWriter, _ := getBatchWriter(true)
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
	wrapperTestPutOrDeleteItem := func(batchWriter *BatchWriter, cases testCases) {
		// Make sure the flush amount is larger than the number of items to add.
		batchWriter.FlushAmount = len(cases) * 2
		for i := 0; i < len(cases); i++ {
			c := cases[i]
			if c.put {
				batchWriter.Add(WrapPutItem(&dynamodb.PutRequest{
					Item: c.item,
				}))
			} else {
				batchWriter.Add(WrapDeleteItem(&dynamodb.DeleteRequest{
					Key: c.item,
				}))
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
		// Repeat an item to check we're overriding by primary keys.
		// This should work regardless of whether the first case is supposed to be
		// a DeleteItem or a PutItem.
		err := batchWriter.Add(WrapPutItem(
			&dynamodb.PutRequest{Item: cases[0].item},
		))
		if err != nil || len(batchWriter.requestBuffer) != len(cases) {
			t.Error("Failed when removing duplicated items.")
		}
		// Only consider the first primary key now.
		batchWriter.primaryKeys = batchWriter.primaryKeys[:1]
		err = batchWriter.Add(WrapPutItem(
			&dynamodb.PutRequest{Item: cases[1].item},
		))
		if err != nil || len(batchWriter.requestBuffer) != len(cases) {
			t.Error("Failed when removing duplicated items with one primary key.")
		}

	}
	batchWriterWithSort, _ := getBatchWriter(true)
	wrapperTestPutOrDeleteItem(batchWriterWithSort, sharedCases)
	batchWriterOnlyHash, _ := getBatchWriter(false)
	wrapperTestPutOrDeleteItem(batchWriterOnlyHash, testCases{
		{true, false, marshal(itemmap{hashKey: "patented", "id": 218})},
		{true, false, marshal(itemmap{hashKey: "charger", "favourite": 5})},
		{true, false, marshal(itemmap{hashKey: "bottle", "favourite": 5})},
		{true, false, marshal(itemmap{hashKey: "dancing", "name": "internal"})},
		{true, false, marshal(itemmap{
			hashKey:     "fortunate",
			"favourite": 5,
			"message":   "abcdefghijklmnopqrst1234567890",
			"url":       "https://golang.org/",
			"part":      18,
		})},
		{true, false, marshal(itemmap{hashKey: "traffic", "favourite": 5})},
	})
}

func TestInvalidRequestError(t *testing.T) {
	batchWriter, _ := getBatchWriter(true)
	err := batchWriter.Add(dynamodb.WriteRequest{})
	if err == nil {
		t.Error("Should return error on empty WriteRequest.")
	}
	if !batchWriter.Empty() {
		t.Error("Should not append to the buffer when there is an error.")
	}
}

func TestBigBuffer(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter(true)
	cases := sharedCases

	expectedRemainingAfterFlush := len(cases)/2 - 1
	// Testing the case where the buffer is bigger than the FlushAmount.
	batchWriter.FlushAmount = len(cases) - expectedRemainingAfterFlush
	for _, c := range cases {
		if c.put {
			batchWriter.requestBuffer = append(
				batchWriter.requestBuffer, dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: c.item,
					},
				},
			)
		} else {
			batchWriter.requestBuffer = append(
				batchWriter.requestBuffer, dynamodb.WriteRequest{
					DeleteRequest: &dynamodb.DeleteRequest{
						Key: c.item,
					},
				},
			)
		}
	}
	addResponse(dynamoClient, 200, "{}")
	err := batchWriter.Flush()
	if err != nil {
		t.Error("Flush errored when it should not have.")
	}
	if len(batchWriter.requestBuffer) != expectedRemainingAfterFlush {
		t.Errorf("Flushed the wrong number of items. Have %d items after"+
			" flushing when there should be %d.",
			len(batchWriter.requestBuffer), expectedRemainingAfterFlush)
	}
	addResponse(dynamoClient, 200, "{}")
	err = batchWriter.Flush()
	if err != nil {
		t.Error("Flush errored on the second call when it should not have.")
	}
	if !batchWriter.Empty() {
		t.Error("requestBuffer should be empty after flushing twice.")
	}
}

func TestEmptyFlush(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter(true)
	addResponse(dynamoClient, 404, "ERR")
	err := batchWriter.Flush()
	if err != nil {
		t.Error("Flush() should not error when the requestBuffer is empty.")
	}
}

func TestEmpty(t *testing.T) {
	batchWriter, _ := getBatchWriter(true)
	if !batchWriter.Empty() { // BatchWriters should start empty.
		t.Error("batchWriter was initialized not empty.")
	}
	cases := sharedCases
	// flushAmount should be higher than the number of cases, so that we know
	// Empty() should return false.
	batchWriter.FlushAmount = len(cases) * 2
	for i, c := range cases {
		if c.put {
			batchWriter.Add(WrapPutItem(&dynamodb.PutRequest{
				Item: c.item,
			}))
		} else {
			batchWriter.Add(WrapDeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			}))
		}
		if batchWriter.Empty() {
			t.Errorf("Empty() returned a fase positive in iteration %d.", i)
		}
	}
}

func TestFlushError(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter(true)
	// Add a dummy response that should yield an error.
	addResponse(dynamoClient, 404, "ERR")
	cases := sharedCases
	for i, c := range cases {
		if c.put {
			batchWriter.Add(WrapPutItem(&dynamodb.PutRequest{
				Item: c.item,
			}))
		} else {
			batchWriter.Add(WrapDeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			}))
		}
		err := batchWriter.Flush()
		if err == nil {
			t.Errorf("Failed to return an error in iteration %d.", i)
		}
	}
}

func TestFlushUnprocessed(t *testing.T) {
	batchWriter, dynamoClient := getBatchWriter(true)
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
			batchWriter.Add(WrapPutItem(&dynamodb.PutRequest{
				Item: c.item,
			}))
		} else {
			batchWriter.Add(WrapDeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			}))
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
	batchWriter, dynamoClient := getBatchWriter(true)
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
			batchWriter.Add(WrapPutItem(&dynamodb.PutRequest{
				Item: c.item,
			}))
		} else {
			batchWriter.Add(WrapDeleteItem(&dynamodb.DeleteRequest{
				Key: c.item,
			}))
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
