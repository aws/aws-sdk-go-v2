package enhancedclient

import (
	"context"
	"crypto/rand"
	"fmt"
	"iter"
	"math"
	"math/big"
	rand2 "math/rand/v2"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func testTableSetup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	t.Log("wuba luba dub dub")

	rnd, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt))

	tableName := fmt.Sprintf("test_e2e_%d", rnd.Int64())

	sch, err := NewSchema[order]()
	if err != nil {
		t.Fatalf("NewSchema() error: %v", err)
	}

	sch.WithTableName(pointer(tableName))

	{
		var tags []types.Tag
		for i := range 30 {
			tags = append(tags, types.Tag{
				Key:   pointer(fmt.Sprintf("key%d", i)),
				Value: pointer(fmt.Sprintf("value%d", i)),
			})
		}
		sch.WithTags(tags)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	c := dynamodb.NewFromConfig(cfg)
	tbl, err := NewTable[order](c, func(options *TableOptions[order]) {
		options.Schema = sch
	})
	if err != nil {
		t.Fatalf("NewTable() error: %v", err)
	}

	// create
	t.Logf("Table %s will be created", tableName)
	err = tbl.CreateWithWait(context.Background(), time.Minute*5)
	if err != nil {
		t.Fatalf("CreateWithWait() error: %v", err)
	}
	t.Logf("Table %s ready", tableName)

	// exists
	t.Logf("Table %s will be checked if it exists", tableName)
	exists, err := tbl.Exists(context.Background())
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists != true {
		t.Fatal("Expected table to exist")
	}
	t.Logf("Table %s exists", tableName)

	_ = tbl
}

func testTableTearDown(t *testing.T) {}

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

func TestTableGetItem(t *testing.T) {
	cases := []struct {
		items  []map[string]types.AttributeValue
		errors map[int]error
	}{
		{
			items: []map[string]types.AttributeValue{
				makeItem[order](),
				makeItem[order](),
				makeItem[order](),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := &mockClient{
				Items:  c.items,
				Errors: c.errors,
			}
			table, err := NewTable[order](client)
			if err != nil {
				t.Errorf("Unexpcted table error: %v", err)
			}

			for _, itm := range c.items {
				actual, err := table.GetItem(context.TODO(), Map{})
				assertItem(t, itm, actual)
				if err != nil {
					t.Errorf(err.Error())
				}
			}
		})
	}
}

func TestTablePutItem(t *testing.T) {
	cases := []struct {
		items  []map[string]types.AttributeValue
		errors map[int]error
	}{
		{
			items: []map[string]types.AttributeValue{
				makeItem[order](),
				makeItem[order](),
				makeItem[order](),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := &mockClient{
				Items:  c.items,
				Errors: c.errors,
			}
			table, err := NewTable[order](client)
			if err != nil {
				t.Errorf("Unexpcted table error: %v", err)
			}

			for _, itm := range c.items {
				o, _ := table.options.Schema.Decode(itm)
				actual, err := table.PutItem(context.TODO(), o)
				assertItem(t, itm, actual)
				if err != nil {
					t.Errorf(err.Error())
				}
			}
		})
	}
}

func TestTableUpdateItem(t *testing.T) {
	cases := []struct {
		items  []map[string]types.AttributeValue
		errors map[int]error
	}{
		{
			items: []map[string]types.AttributeValue{
				makeItem[order](),
				makeItem[order](),
				makeItem[order](),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := &mockClient{
				Items:  c.items,
				Errors: c.errors,
			}
			table, err := NewTable[order](client)
			if err != nil {
				t.Errorf("Unexpcted table error: %v", err)
			}

			for _, itm := range c.items {
				o, _ := table.options.Schema.Decode(itm)
				actual, err := table.UpdateItem(context.TODO(), o)
				assertItem(t, itm, actual)
				if err != nil {
					t.Errorf(err.Error())
				}
			}
		})
	}
}

func TestTableDeleteItem(t *testing.T) {
}

func TestTableQuery(t *testing.T) {
	cases := []struct {
		items  []map[string]types.AttributeValue
		errors map[int]error
	}{
		{
			items: []map[string]types.AttributeValue{
				makeItem[order](),
				makeItem[order](),
				makeItem[order](),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := &mockClient{
				Items:  c.items,
				Errors: c.errors,
			}
			table, err := NewTable[order](client)
			if err != nil {
				t.Errorf("Unexpcted table error: %v", err)
			}

			next, _ := iter.Pull(table.Query(context.TODO(), expression.Expression{}))

			for _, itm := range c.items {
				actual, _ := next()
				item := actual.Item()
				assertItem(t, itm, &item)
			}
		})
	}
}

func TestTableScan(t *testing.T) {
	cases := []struct {
		items  []map[string]types.AttributeValue
		errors map[int]error
	}{
		{
			items: []map[string]types.AttributeValue{
				makeItem[order](),
				makeItem[order](),
				makeItem[order](),
			},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := &mockClient{
				Items:  c.items,
				Errors: c.errors,
			}
			table, err := NewTable[order](client)
			if err != nil {
				t.Errorf("Unexpcted table error: %v", err)
			}

			next, _ := iter.Pull(table.Scan(context.TODO(), expression.Expression{}))

			for _, itm := range c.items {
				actual, _ := next()
				item := actual.Item()
				assertItem(t, itm, &item)
			}
		})
	}
}
