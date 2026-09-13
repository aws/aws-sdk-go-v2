package entitymanager

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestTableBatchGetItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 9}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 8}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 7}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 6}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 0}),
				// even tho we request initally all the items, we expect 2 items to be left unprocessed
				// because we are forcing the UnprocessedKeys to be empty in last call
				withExpectFns(expectItemsCount("order", 2)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchGetItemCall(errors.New("1"), map[string]uint{"order": 0}),
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

			for _, item := range c.client.(*mockClient).Items["order"] {
				bgio.AddReadItemByMap(item)
			}

			for res := range bgio.Execute(context.Background()) {
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

func TestTableMultiBatchGetItem(t *testing.T) {
	cases := []struct {
		client        Client
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withItems("order_backup", makeItem[order], 32),
				// as the pool of items for order table is diminished, the request of order_backup will increase
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 9, "order_backup": 0}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 8, "order_backup": 2}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 7, "order_backup": 9}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 6, "order_backup": 15}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 0, "order_backup": 4}),
				withDefaultBatchGetItemCall(nil, map[string]uint{"order": 0, "order_backup": 0}),
				// even tho we request initally all the items, we expect 2 items to be left unprocessed
				// because we are forcing the UnprocessedKeys to be empty in last call
				withExpectFns(expectItemsCount("order", 2)),
				withExpectFns(expectItemsCount("order_backup", 2)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchGetItemCall(errors.New("1"), map[string]uint{"order": 0}),
			),
			expectedError: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			defer c.client.(*mockClient).RunExpectations(t)

			table, err := NewTable[order](c.client)
			table2, err2 := NewTable[order](c.client, func(options *TableOptions[order]) {
				sch, _ := NewSchema[order]()
				sch.WithTableName(aws.String("order_backup"))
				options.Schema = sch
			})
			if err != nil {
				t.Errorf("unexpcted table error: %v", err)
			}
			if err2 != nil {
				t.Errorf("unexpcted table error: %v", err)
			}

			bgio := table.CreateBatchGetOperation()
			bgio2 := table2.CreateBatchGetOperation()

			for _, item := range c.client.(*mockClient).Items["order"] {
				bgio.AddReadItemByMap(item)
			}
			for _, item := range c.client.(*mockClient).Items["order_backup"] {
				bgio2.AddReadItemByMap(item)
			}

			executor := bgio.Merge(bgio2)

			for res := range executor.Execute(context.Background()) {
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
