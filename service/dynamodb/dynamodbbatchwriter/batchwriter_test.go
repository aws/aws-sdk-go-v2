package dynamodbbatchwriter

import (
	"reflect"
	"testing"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

type mockDynamoDBClient struct{
	dynamodbiface.DynamoDBAPI
}

func TestNewBatchWriter(t *testing.T) {
	dynamoClient := &mockDynamoDBClient{}
	tableName := "testtable"

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
	dynamoClient := &mockDynamoDBClient{}
	tableName := "testtable"

	batchWriter := New(tableName, dynamoClient)

	testValues := []int{20, 150, 1, 12, 112}
	for val := range(testValues) {
		batchWriter.SetFlushAmount(val)
		if batchWriter.flushAmount != val {
			t.Errorf("batchWriter.flushAmount is set to %v instead of %d",
				batchWriter.flushAmount, val)
		}
	}
}

type itemmap map[string]interface{}

func marshal(in interface{}) map[string]dynamodb.AttributeValue {
	out, _ := dynamodbattribute.MarshalMap(in)
	return out
}

func TestPutOrDeleteItem(t *testing.T) {
	dynamoClient := &mockDynamoDBClient{}
	tableName := "testtable"

	batchWriter := New(tableName, dynamoClient)

	cases := []struct{
		put, delete bool
		item map[string]dynamodb.AttributeValue
	}{
		{true, false, marshal(itemmap{"itemcount": []int{89, 91, 92}})},
		{true, false, marshal(itemmap{"dance": 1})},
		{false, true, marshal(itemmap{"word": "dance"})},
		{true, false, marshal(itemmap{"func": itemmap{"in": 1, "out": 2}})},
		{false, true, marshal(itemmap{"tel": 123, "pos": 2})},
	}
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
	dynamoClient := &mockDynamoDBClient{}
	tableName := "testtable"

	batchWriter := New(tableName, dynamoClient)
	if !batchWriter.Empty() { // BatchWriters should start empty.
		t.Error("batchWriter was initialized not empty.")
	}
	cases := []struct{
		put, delete bool
		item map[string]dynamodb.AttributeValue
	}{
		{true, false, marshal(itemmap{"itemcount": []int{89, 91, 92}})},
		{true, false, marshal(itemmap{"dance": 1})},
		{false, true, marshal(itemmap{"word": "dance"})},
		{true, false, marshal(itemmap{"func": itemmap{"in": 1, "out": 2}})},
		{false, true, marshal(itemmap{"tel": 123, "pos": 2})},
	}
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
