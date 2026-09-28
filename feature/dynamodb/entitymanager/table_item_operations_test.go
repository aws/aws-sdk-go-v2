package entitymanager

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

			_, err = table.Create(context.Background())
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

			_, _ = table.Create(context.Background())

			_, err = table.Describe(context.Background())
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

			_, _ = table.Create(context.Background())

			_, err = table.Delete(context.Background())
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
				withItem("order", makeItem[order]()),
			),
		},
		{
			client: newMockClient(
				withDefaultGetItemCall(errors.New("1")),
				withItem("order", makeItem[order]()),
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

			_, err = table.GetItem(context.Background(), Map{})
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
				withExpectFns(expectItemsCount("order", 1)),
			),
		},
		{
			client: newMockClient(
				withDefaultPutItemCall(errors.New("1")),
				withExpectFns(expectItemsCount("order", 0)),
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

			_, err = table.PutItem(context.Background(), &order{})
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
				withExpectFns(expectItemsCount("order", 1)),
			),
		},
		{
			client: newMockClient(
				withDefaultUpdateItemCall(errors.New("1")),
				withExpectFns(expectItemsCount("order", 0)),
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

			_, err = table.UpdateItem(context.Background(), &order{})
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
				withItems("order", makeItem[order], 2),
				withDefaultDeleteItemCall(nil),
				withExpectFns(expectItemsCount("order", 1)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 2),
				withDefaultDeleteItemCall(errors.New("1")),
				withExpectFns(expectItemsCount("order", 2)),
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

			err = table.DeleteItem(context.Background(), &order{})
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
				withItems("order", makeItem[order], 32),
				withDefaultQueryCall(nil, 9),
				withDefaultQueryCall(nil, 8),
				withDefaultQueryCall(nil, 7),
				withDefaultQueryCall(nil, 6),
				withDefaultQueryCall(nil, 0),
				withExpectFns(expectItemsCount("order", 2)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
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

			for res := range table.Query(context.Background(), expression.Expression{}) {
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
				withItems("order", makeItem[order], 32),
				withDefaultScanCall(nil, 9),
				withDefaultScanCall(nil, 8),
				withDefaultScanCall(nil, 7),
				withDefaultScanCall(nil, 6),
				withDefaultScanCall(nil, 0),
				withExpectFns(expectItemsCount("order", 2)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
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

			for res := range table.Scan(context.Background(), expression.Expression{}) {
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
