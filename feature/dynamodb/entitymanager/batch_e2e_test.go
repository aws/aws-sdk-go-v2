package entitymanager

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func getTables(t *testing.T, count int) []*Table[order] {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	c := dynamodb.NewFromConfig(cfg)

	tables := make([]*Table[order], count)

	wg := sync.WaitGroup{}
	errChan := make(chan error, len(tables))
	defer close(errChan)

	for i := range len(tables) {
		sch, err := NewSchema[order]()
		if err != nil {
			t.Fatalf("NewTable() error: %v", err)
		}

		tableName := fmt.Sprintf("test_batch_e2e_%s_%d", time.Now().Format("2006_01_02_15_04_05.000000000"), i)

		sch.WithTableName(&tableName)

		tables[i], err = NewTable[order](
			c,
			WithSchema(sch),
		)
		if err != nil {
			t.Fatalf("NewTable() error: %v", err)
		}

		wg.Add(1)

		go func(table *Table[order]) {
			defer wg.Done()

			if err := table.CreateWithWait(context.Background(), time.Minute); err != nil {
				errChan <- err
			}
		}(tables[i])
	}

	wg.Wait()

	if len(errChan) > 0 {
		for err := range errChan {
			if err != nil {
				t.Fatalf("CreateWithWait() error: %v", err)
			}
		}
	}

	return tables
}

func TestTableBatchE2E(t *testing.T) {
	t.Parallel()

	numItems := 32

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	c := dynamodb.NewFromConfig(cfg)

	sch, err := NewSchema[order]()
	if err != nil {
		t.Fatalf("NewTable() error: %v", err)
	}

	tableName := fmt.Sprintf("test_batch_e2e_%s", time.Now().Format("2006_01_02_15_04_05.000000000"))
	sch.WithTableName(&tableName)

	table, err := NewTable[order](
		c,
		WithSchema(sch),
	)
	if err != nil {
		t.Fatalf("NewTable() error: %v", err)
	}

	if err := table.CreateWithWait(context.Background(), time.Minute); err != nil {
		t.Fatalf("Error during CreateWithWait(): %v", err)
	}

	batchWrite := table.CreateBatchWriteOperation()
	now := time.Now()
	for c := range numItems {
		batchWrite.AddPut(&order{
			OrderID:      fmt.Sprintf("order:%d", c),
			CreatedAt:    now.Unix(),
			CustomerID:   fmt.Sprintf("customer:%d", c),
			TotalAmount:  42.1337,
			customerNote: fmt.Sprintf("note:%d", c),
			address: address{
				Street: fmt.Sprintf("steet:%d", c),
				City:   fmt.Sprintf("city:%d", c),
				Zip:    fmt.Sprintf("zip:%d", c),
			},
		})
	}

	if err := batchWrite.Execute(context.Background()); err != nil {
		t.Fatalf("Error during Execute(): %v", err)
	}

	batchGet := table.CreateBatchGetOperation()
	for c := range numItems {
		batchGet.AddReadItemByMap(Map{}.With("order_id", fmt.Sprintf("order:%d", c)).With("created_at", now.Unix()))
	}

	items := make([]*order, 0, 32)
	for res := range batchGet.Execute(context.Background()) {
		if res.Error() != nil {
			t.Errorf("Error during get: %v", res.Error())

			continue
		}

		items = append(items, res.Item())
	}

	if len(items) != numItems {
		t.Errorf("Expected to fetch %d number, got %d", numItems, len(items))
	}

	defer func() {
		if err := table.DeleteWithWait(context.Background(), time.Minute); err != nil {
			t.Fatalf("Error during DeleteWithWait(): %v", err)
		}
	}()
}

func TestTableMultiBatchE2E(t *testing.T) {
	t.Parallel()

	numItems := 32
	numTables := 3

	tables := getTables(t, numTables) // must be higher than 2
	t.Cleanup(func() {
		errChan := make(chan error, len(tables))
		defer close(errChan)
		wg := sync.WaitGroup{}

		for _, table := range tables {
			wg.Add(1)

			go func(table *Table[order]) {
				defer wg.Done()

				if err := table.DeleteWithWait(context.Background(), time.Minute); err != nil {
					errChan <- err
				}
			}(table)
		}

		wg.Wait()

		if len(errChan) > 0 {
			for err := range errChan {
				if err != nil {
					t.Fatalf("DeleteWithWait() error: %v", err)
				}
			}
		}
	})

	now := time.Now()

	// put items
	batchWrites := make([]*BatchWriteOperation[order], len(tables))
	for i := range tables {
		batchWrites[i] = tables[i].CreateBatchWriteOperation()
		for c := range numItems {
			batchWrites[i].AddPut(&order{
				OrderID:      fmt.Sprintf("order:%d", c),
				CreatedAt:    now.Unix(),
				CustomerID:   fmt.Sprintf("customer:%d", c),
				TotalAmount:  42.1337,
				customerNote: fmt.Sprintf("note:%d", c),
				address: address{
					Street: fmt.Sprintf("steet:%d", c),
					City:   fmt.Sprintf("city:%d", c),
					Zip:    fmt.Sprintf("zip:%d", c),
				},
			})
		}
	}

	writeExecutor := batchWrites[0].Merge(batchWrites[1])
	for c := 2; c < len(batchWrites); c++ {
		writeExecutor = writeExecutor.Merge(batchWrites[c])
	}
	if err := writeExecutor.Execute(context.Background()); err != nil {
		t.Errorf("Execute() error: %v", err)
	}

	// read items
	batchGets := make([]*BatchGetOperation[order], len(tables))
	for i := range tables {
		batchGets[i] = tables[i].CreateBatchGetOperation()
		for c := range numItems {
			batchGets[i].AddReadItemByMap(Map{}.With("order_id", fmt.Sprintf("order:%d", c)).With("created_at", now.Unix()))
		}
	}

	getExecutor := batchGets[0].Merge(batchGets[1])
	for c := 2; c < len(batchGets); c++ {
		getExecutor = getExecutor.Merge(batchGets[c])
	}

	found := make(map[string][]string)
	for res := range getExecutor.Execute(context.Background()) {
		if res.Error() != nil {
			t.Errorf("Error during get: %v", res.Error())

			continue
		}

		if item, ok := res.Item().(*order); ok {
			f := found[item.OrderID]
			f = append(f, res.Table())
			found[item.OrderID] = f
		}
	}

	for orderId, tableNames := range found {
		if len(tables) != len(tableNames) {
			t.Logf(`Order ID "%s" was not found in all tables: %v`, orderId, tableNames)
		}
	}
}
