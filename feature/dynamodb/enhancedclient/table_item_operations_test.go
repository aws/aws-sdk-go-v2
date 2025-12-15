package enhancedclient

import (
	"context"
	"errors"
	"fmt"
	"math"
	rand2 "math/rand/v2"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func makeField(name string, t reflect.Type) types.AttributeValue {
	k := t.Kind()
	switch k {
	case reflect.String:
		if strings.Contains(name, "version") {
			return &types.AttributeValueMemberS{
				Value: fmt.Sprintf("%d", rand2.Int32N(100)),
			}
		}

		return &types.AttributeValueMemberS{
			Value: strings.Repeat(string(byte(rand2.UintN(93)+33)), rand2.IntN(100)),
		}
	case reflect.Int64:
		return &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", rand2.Int64N(math.MaxInt)),
		}
	case reflect.Float64:
		return &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d.%d", rand2.Int64N(math.MaxInt), rand2.Int64N(math.MaxInt)),
		}
	case reflect.Map:
		m := map[string]types.AttributeValue{}
		for c := range 10 {
			s := fmt.Sprintf("%d", c)
			m[s] = makeField(s, reflect.TypeFor[string]())
		}
		return &types.AttributeValueMemberM{
			Value: m,
		}
	case reflect.Slice, reflect.Array:
		l := []types.AttributeValue{}
		for range 10 {
			l = append(l, makeField(name, t.Elem()))
		}
		return &types.AttributeValueMemberL{
			Value: l,
		}
	}
	return nil
}

func makeItem[T any]() map[string]types.AttributeValue {
	s, _ := NewSchema[T]()

	out := map[string]types.AttributeValue{}

	for _, f := range s.cachedFields.All() {
		out[f.Name] = makeField(f.Name, f.Type)
	}

	return out
}

func assertField[V any](t *testing.T, i map[string]types.AttributeValue, key string, value V) {
	var rv V
	err := NewDecoder[V]().Decode(i[key], &rv)
	if err != nil {
		t.Errorf(`unable to decode "%v"`, i[key])
		return
	}
	if diff := cmpDiff(rv, value); diff != "" {
		t.Errorf(`enexpected diff for "%s": %v`, key, diff)
	}
}

func assertItem(t *testing.T, i map[string]types.AttributeValue, o *order) {
	if o == nil {
		t.Error(`order is nil`)
		return
	}

	assertField(t, i, "order_id", o.OrderID)
	assertField(t, i, "customer_id", o.CustomerID)
	//assertField(t, i, "versionString", o.VersionString)
	assertField(t, i, "street", o.Street)
	assertField(t, i, "city", o.City)
	assertField(t, i, "zip", o.Zip)
	assertField(t, i, "note", o.customerNote)
	assertField(t, i, "first_name", o.CustomerFirstName)
	assertField(t, i, "last_name", o.CustomerLastName)
	//assertField(t, i, "created_at", o.CreatedAt)
	// float fields have garbage :(
	//assertField(t, i, "total", o.TotalAmount)
	//assertField(t, i, "version", o.Version)
	//assertField(t, i, "counter_up", o.CounterUp)
	//assertField(t, i, "counter_down", o.CounterDown)
}

func TestTableCreate(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultCreateTableCall(nil),
				withExpectFns(expectTablesCount(1)),
				withExpectFns(expectTable("order")),
			),
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(errors.New("1")),
				withExpectFns(expectTablesCount(0)),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, err = table.Create(context.TODO())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableDescribe(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultCreateTableCall(nil),
				withDefaultDescribeTableCall(nil),
				withExpectFns(expectTablesCount(1)),
				withExpectFns(expectTable("order")),
			),
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(errors.New("1")),
				withDefaultDescribeTableCall(errors.New("1")),
				withExpectFns(expectTablesCount(0)),
			),
			expectedError: true,
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(nil),
				withDefaultDescribeTableCall(errors.New("1")),
				withExpectFns(expectTablesCount(1)),
				withExpectFns(expectTable("order")),
			),
			expectedError: true,
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(errors.New("1")),
				withDefaultDescribeTableCall(nil),
				withExpectFns(expectTablesCount(0)),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, _ = table.Create(context.TODO())

			_, err = table.Describe(context.TODO())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableDelete(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultCreateTableCall(nil),
				withDefaultDeleteTableCall(nil),
				withExpectFns(expectTablesCount(0)),
			),
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(errors.New("1")),
				withDefaultDeleteTableCall(nil),
				withExpectFns(expectTablesCount(0)),
			),
			expectedError: true,
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(errors.New("1")),
				withDefaultDeleteTableCall(nil),
				withExpectFns(expectTablesCount(0)),
			),
			expectedError: true,
		},
		{
			client: newMockClient(
				withDefaultCreateTableCall(nil),
				withDefaultDeleteTableCall(errors.New("1")),
				withExpectFns(expectTablesCount(1)),
				withExpectFns(expectTable("order")),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, _ = table.Create(context.TODO())

			_, err = table.Delete(context.TODO())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableGetItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultGetItemCall(nil),
				withItem(makeItem[order]()),
			),
		},
		{
			client: newMockClient(
				withDefaultGetItemCall(errors.New("1")),
				withItem(makeItem[order]()),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, err = table.GetItem(context.TODO(), Map{})
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTablePutItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultPutItemCall(nil),
				withExpectFns(expectItemsCount(1)),
			),
		},
		{
			client: newMockClient(
				withDefaultPutItemCall(errors.New("1")),
				withExpectFns(expectItemsCount(0)),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, err = table.PutItem(context.TODO(), &order{})
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableUpdateItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withDefaultUpdateItemCall(nil),
				withExpectFns(expectItemsCount(1)),
			),
		},
		{
			client: newMockClient(
				withDefaultUpdateItemCall(errors.New("1")),
				withExpectFns(expectItemsCount(0)),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			_, err = table.UpdateItem(context.TODO(), &order{})
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableDeleteItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems(makeItem[order], 2),
				withDefaultDeleteItemCall(nil),
				withExpectFns(expectItemsCount(1)),
			),
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 2),
				withDefaultDeleteItemCall(errors.New("1")),
				withExpectFns(expectItemsCount(2)),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			err = table.DeleteItem(context.TODO(), &order{})
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableQuery(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultQueryCall(nil, 9),
				withDefaultQueryCall(nil, 8),
				withDefaultQueryCall(nil, 7),
				withDefaultQueryCall(nil, 6),
				withDefaultQueryCall(nil, 0),
				withExpectFns(expectItemsCount(2)),
			),
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultQueryCall(errors.New("1"), 0),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			for res := range table.Query(context.TODO(), expression.Expression{}) {
				if c.expectedError && res.Error() == nil {
					t.Fatalf("expected error but got none")
				}

				if !c.expectedError && res.Error() != nil {
					t.Fatalf("unexpected error: %v", res.Error())
				}
			}
		})
	}
}

func TestTableScan(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultScanCall(nil, 9),
				withDefaultScanCall(nil, 8),
				withDefaultScanCall(nil, 7),
				withDefaultScanCall(nil, 6),
				withDefaultScanCall(nil, 0),
				withExpectFns(expectItemsCount(2)),
			),
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultScanCall(errors.New("1"), 0),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			for res := range table.Scan(context.TODO(), expression.Expression{}) {
				if c.expectedError && res.Error() == nil {
					t.Fatalf("expected error but got none")
				}

				if !c.expectedError && res.Error() != nil {
					t.Fatalf("unexpected error: %v", res.Error())
				}
			}
		})
	}
}

func TestTableBatchGetItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultBatchGetItemCall(nil, 9, "order"),
				withDefaultBatchGetItemCall(nil, 8, "order"),
				withDefaultBatchGetItemCall(nil, 7, "order"),
				withDefaultBatchGetItemCall(nil, 6, "order"),
				withDefaultBatchGetItemCall(nil, 0, "order"),
				// even tho we request initally all the items, we expect 2 items to be left unprocessed
				// because we are forcing the UnprocessedKeys to be empty in last call
				withExpectFns(expectItemsCount(2)),
			),
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultBatchGetItemCall(errors.New("1"), 0, "order"),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			bgio := table.CreateBatchGetOperation()

			for _, item := range c.client.(*mockClient).Items {
				bgio.AddReadItemByMap(item)
			}

			for res := range bgio.Execute(context.TODO()) {
				if c.expectedError && res.Error() == nil {
					t.Fatalf("expected error but got none")
				}

				if !c.expectedError && res.Error() != nil {
					t.Fatalf("unexpected error: %v", res.Error())
				}
			}
		})
	}
}

func TestTableBatchWriteItem(t *testing.T) {

	cases := []struct {
		client        Client
		isDelete      bool
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultBatchWriteItemCall(nil, 9, "order"),
				withDefaultBatchWriteItemCall(nil, 8, "order"),
				withDefaultBatchWriteItemCall(nil, 7, "order"),
				withDefaultBatchWriteItemCall(nil, 6, "order"),
				withDefaultBatchWriteItemCall(nil, 0, "order"),
				withExpectFns(expectItemsCount(62)),
			),
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultBatchWriteItemCall(nil, 9, "order"),
				withDefaultBatchWriteItemCall(nil, 8, "order"),
				withDefaultBatchWriteItemCall(nil, 7, "order"),
				withDefaultBatchWriteItemCall(nil, 6, "order"),
				withDefaultBatchWriteItemCall(nil, 0, "order"),
				withExpectFns(expectItemsCount(2)),
			),
			isDelete: true,
		},
		{
			client: newMockClient(
				withItems(makeItem[order], 32),
				withDefaultBatchWriteItemCall(errors.New("1"), 0, "order"),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			bgwo := table.CreateBatchWriteOperation()

			for range 32 {
				if c.isDelete {
					bgwo.AddRawDelete(makeItem[order]())
				} else {
					bgwo.AddRawPut(makeItem[order]())
				}
			}

			err = bgwo.Execute(context.TODO())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
