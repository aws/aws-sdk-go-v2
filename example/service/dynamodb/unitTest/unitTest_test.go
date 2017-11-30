// +build example

// Unit tests for package unitTest.
package unitTest

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbiface"
)

// A fakeDynamoDB instance. During testing, instatiate ItemGetter, then simply
// assign an instance of fakeDynamoDB to it.
type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	payload map[string]string // Store expected return values
	err     error
}

// Mock GetItem such that the output returned carries values identical to input.
func (fd *fakeDynamoDB) GetItemRequest(input *dynamodb.GetItemInput) dynamodb.GetItemRequest {
	output := &dynamodb.GetItemOutput{
		Item: map[string]dynamodb.AttributeValue{},
	}
	for key, value := range fd.payload {
		output.Item[key] = dynamodb.AttributeValue{
			S: aws.String(value),
		}
	}
	req := dynamodb.GetItemRequest{
		Request: &aws.Request{
			Data:  output,
			Error: fd.err,
		},
	}

	return req
}

func TestItemGetterGet(t *testing.T) {
	expectedKey := "expected key"
	expectedValue := "expected value"
	getter := new(ItemGetter)
	getter.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{"id": expectedKey, "value": expectedValue},
	}
	if actualValue := getter.Get(expectedKey); actualValue != expectedValue {
		t.Errorf("Expected %q but got %q", expectedValue, actualValue)
	}
}

// When DynamoDB.GetItem returns a non-nil error, expect an empty string.
func TestItemGetterGetFail(t *testing.T) {
	expectedKey := "expected key"
	expectedValue := "expected value"
	getter := new(ItemGetter)
	getter.DynamoDB = &fakeDynamoDB{
		payload: map[string]string{"id": expectedKey, "value": expectedValue},
		err:     errors.New("any error"),
	}
	if actualValue := getter.Get(expectedKey); len(actualValue) > 0 {
		t.Errorf("Expected %q but got %q", expectedValue, actualValue)
	}
}
