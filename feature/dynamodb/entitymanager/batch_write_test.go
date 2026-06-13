package entitymanager

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestTableBatchWriteItem(t *testing.T) {
	cases := []struct {
		client        Client
		isDelete      bool
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 9}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 8}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 7}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 6}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0}),
				withExpectFns(expectItemsCount("order", 62)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 9}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 8}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 7}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 6}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0}),
				withExpectFns(expectItemsCount("order", 2)),
			),
			isDelete: true,
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchWriteItemCall(errors.New("1"), map[string]uint{"order": 0}),
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

			err = bgwo.Execute(context.Background())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestTableMultiBatchWriteItem(t *testing.T) {
	cases := []struct {
		client        Client
		isDelete      bool
		expectedError bool
	}{
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withItems("order_backup", makeItem[order], 32),
				// as the pool of items for order table is diminished, the request of order_backup will increase
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 9, "order_backup": 0}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 8, "order_backup": 2}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 7, "order_backup": 9}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 6, "order_backup": 15}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0, "order_backup": 4}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0, "order_backup": 0}),
				withExpectFns(expectItemsCount("order", 62)),
				withExpectFns(expectItemsCount("order_backup", 62)),
			),
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withItems("order_backup", makeItem[order], 32),
				// as the pool of items for order table is diminished, the request of order_backup will increase
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 9, "order_backup": 0}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 8, "order_backup": 2}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 7, "order_backup": 9}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 6, "order_backup": 15}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0, "order_backup": 4}),
				withDefaultBatchWriteItemCall(nil, map[string]uint{"order": 0, "order_backup": 0}),
				withExpectFns(expectItemsCount("order", 2)),
				withExpectFns(expectItemsCount("order_backup", 2)),
			),
			isDelete: true,
		},
		{
			client: newMockClient(
				withItems("order", makeItem[order], 32),
				withDefaultBatchWriteItemCall(errors.New("1"), map[string]uint{"order": 0}),
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

			bgwo := table.CreateBatchWriteOperation()
			bgwo2 := table2.CreateBatchWriteOperation()

			for range 32 {
				if c.isDelete {
					bgwo.AddRawDelete(makeItem[order]())
					bgwo2.AddRawDelete(makeItem[order]())
				} else {
					bgwo.AddRawPut(makeItem[order]())
					bgwo2.AddRawPut(makeItem[order]())
				}
			}

			err = bgwo.Merge(bgwo2).Execute(context.Background())
			if c.expectedError && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !c.expectedError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
