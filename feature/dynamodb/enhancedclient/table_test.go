package enhancedclient

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MyAuditExtension struct{}

func (a *MyAuditExtension) BeforeWrite(ctx context.Context, v *order) error {
	log.Printf("Audit: about to write item: %+v", v.OrderID)
	return nil
}

func (a *MyAuditExtension) AfterRead(ctx context.Context, v *order) error {
	log.Printf("Audit: read item: %+v", v.OrderID)
	return nil
}

func TestTableE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	tableName := fmt.Sprintf("test_e2e_%s", time.Now().Format("2006_01_02_15_04_05"))

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

	ext := &MyAuditExtension{}
	registry := DefaultExtensionRegistry[order]().Clone()
	registry.AddBeforeWriter(ext)
	registry.AddAfterReader(ext)

	tbl, err := NewTable[order](
		c,
		WithSchema(sch),
		WithExtensionRegistry(registry),
	)
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

	// defer table delete
	t.Cleanup(func() {
		t.Logf("Table %s will be deleted", tableName)
		err = tbl.DeleteWithWait(context.Background(), time.Minute)
		if err != nil {
			t.Fatalf("DeleteWithWait() error: %v", err)
		}
		t.Logf("Table %s deleted", tableName)
	})

	itemsToManage := 128
	orderIds := make([]string, itemsToManage)
	createdAts := make([]int64, itemsToManage)
	// Put()
	for i := range itemsToManage {
		o := &order{
			CustomerID:    fmt.Sprintf("CustomerID%d", i),
			TotalAmount:   float64(i),
			IgnoredField:  fmt.Sprintf("IgnoredField%d", i),
			Version:       0,
			VersionString: "0",
			CounterUp:     0,
			CounterDown:   0,
			Metadata: map[string]string{
				"test": "test",
			},
			address: address{
				Street: fmt.Sprintf("Street%d", i),
				City:   fmt.Sprintf("City%d", i),
				Zip:    fmt.Sprintf("Zip%d", i),
			},
			Notes:             []string{fmt.Sprintf("Notes%d", i)},
			customerNote:      fmt.Sprintf("customerNote%d", i),
			CustomerFirstName: fmt.Sprintf("CustomerFirstName%d", i),
			CustomerLastName:  fmt.Sprintf("CustomerLastName%d", i),
		}

		item, err := tbl.PutItem(context.Background(), o)
		if err != nil {
			t.Errorf("Unable to PutItem() [%d]: %v", i, err)
			continue
		}

		orderIds[i] = item.OrderID
		createdAts[i] = item.CreatedAt

		t.Logf("PutItem: %s - %d", item.OrderID, item.CreatedAt)
		t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
		t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)
	}

	// Get() + Update()
	for i := range itemsToManage {
		m := Map{}.
			With("order_id", orderIds[i]).
			With("created_at", createdAts[i])
		item, err := tbl.GetItem(context.Background(), m)
		if err != nil {
			t.Errorf("Unable to GetItem() [%s]: %v", m, err)
			continue
		}
		t.Logf("GetItem: %s - %d", item.OrderID, item.CreatedAt)
		t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
		t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)

		item.TotalAmount *= 2

		item, err = tbl.UpdateItem(context.Background(), item)
		if err != nil {
			t.Errorf("Unable to UpdateItem() [%s]: %v", m, err)
			continue
		}
		t.Logf("UpdateItem: %s - %d", item.OrderID, item.CreatedAt)
		t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
		t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)
	}

	{
		t.Log("Scan()")
		scanExpr := expression.Expression{}
		items := tbl.Scan(context.Background(), scanExpr)
		scannedItems := 0
		for res := range items {
			if res.Error() != nil {
				t.Errorf("Error during Scan(): %v", res.Error())
			}
			item := res.Item()
			t.Logf("Scan: %s - %d", item.OrderID, item.CreatedAt)
			t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
			t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)

			scannedItems++
		}
		if scannedItems != itemsToManage {
			t.Errorf("Scanned %d item(s), expected %d", scannedItems, itemsToManage)
		}
	}

	{
		t.Log("ScanIndex()")
		scanExpr := expression.Expression{}
		items := tbl.ScanIndex(context.Background(), "CustomerIndex", scanExpr)
		scannedItems := 0
		for res := range items {
			if res.Error() != nil {
				t.Errorf("Error during ScanIndex(): %v", res.Error())

				continue
			}

			item := res.Item()
			if item != nil {
				t.Logf("Scan: %s - %d", item.OrderID, item.CreatedAt)
				t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
				t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)
			} else {
				t.Log("no error and item was nil :(")
			}

			scannedItems++
		}
		if scannedItems != itemsToManage {
			t.Errorf("Scanned %d item(s), expected %d", scannedItems, itemsToManage)
		}
	}

	knowVersions := map[string]int64{}
	{
		t.Log("Query()")
		queriedItems := 0
		for i := range itemsToManage {
			queryExprBuilder := expression.NewBuilder()
			queryExprBuilder = queryExprBuilder.WithKeyCondition(
				expression.Key("order_id").Equal(expression.Value(orderIds[i])).And(
					expression.Key("created_at").Equal(expression.Value(createdAts[i])),
				),
			)
			queryExpr, err := queryExprBuilder.Build()
			if err != nil {
				t.Errorf("Unable to build query: %v", err)

				return
			}

			items := tbl.Query(context.Background(), queryExpr)
			for res := range items {
				if res.Error() != nil {
					t.Errorf("Error during Query(): %v", res.Error())

					continue
				}

				item := res.Item()
				t.Logf("Query: %s - %d", item.OrderID, item.CreatedAt)
				t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
				t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)

				knowVersions[item.OrderID] = item.Version

				queriedItems++
			}
		}
		if queriedItems != itemsToManage {
			t.Errorf("Queried %d item(s), expected %d", queriedItems, itemsToManage)
		}
	}

	{
		t.Log("QueryIndex()")
		queriedItems := 0
		for orderId, version := range knowVersions {
			queryExprBuilder := expression.NewBuilder()
			queryExprBuilder = queryExprBuilder.WithKeyCondition(
				expression.Key("order_id").Equal(expression.Value(orderId)).And(
					expression.Key("version").Equal(expression.Value(version)),
				),
			)
			queryExpr, err := queryExprBuilder.Build()
			if err != nil {
				t.Errorf("Unable to build query: %v", err)

				return
			}

			items := tbl.QueryIndex(context.Background(), "OrderVersionIndex", queryExpr)
			for res := range items {
				if res.Error() != nil {
					t.Errorf("Error during QueryIndex(): %v", res.Error())

					continue
				}

				item := res.Item()
				t.Logf("Query: %s - %d", item.OrderID, item.CreatedAt)
				t.Logf("\tVersion: %d - %s", item.Version, item.VersionString)
				t.Logf("\tCounter (Up/Down): %d/%d", item.CounterUp, item.CounterDown)

				queriedItems++
			}
		}

		if queriedItems != itemsToManage {
			t.Errorf("Queried %d item(s), expected %d", queriedItems, itemsToManage)
		}
	}

	// batch
	{
		bwo := tbl.CreateBatchWriteOperation()
		batchItems := make([]order, 128)
		for i := range 128 {
			batchItems[i] = order{
				CustomerID:    fmt.Sprintf("CustomerID%d", i),
				TotalAmount:   float64(i),
				IgnoredField:  fmt.Sprintf("IgnoredField%d", i),
				Version:       0,
				VersionString: "0",
				CounterUp:     0,
				CounterDown:   0,
				Metadata: map[string]string{
					"test": "test",
				},
				address: address{
					Street: fmt.Sprintf("Street%d", i),
					City:   fmt.Sprintf("City%d", i),
					Zip:    fmt.Sprintf("Zip%d", i),
				},
				Notes:             []string{fmt.Sprintf("Notes%d", i)},
				customerNote:      fmt.Sprintf("customerNote%d", i),
				CustomerFirstName: fmt.Sprintf("CustomerFirstName%d", i),
				CustomerLastName:  fmt.Sprintf("CustomerLastName%d", i),
			}
			err = bwo.AddPut(&batchItems[i])
			if err != nil {
				t.Error(err.Error())
			}
		}

		err = bwo.Execute(context.TODO())
		if err != nil {
			t.Error(err.Error())
		} else {
			t.Log("BatchWritePut done")
		}
		for _, batchItem := range batchItems {
			t.Logf("OrderID: %s", batchItem.OrderID)
		}

		// get
		bgo := tbl.CreateBatchGetOperation()
		for i := range batchItems {
			if err = bgo.AddReadItem(&batchItems[i]); err != nil {
				t.Error(err.Error())
			}
		}

		for item := range bgo.Execute(context.TODO()) {
			if item.Error() != nil {
				t.Errorf(`error during BatchGetOperation iteration: %v`, err)
				continue
			}

			if item.Item() == nil {
				t.Error("nil item returned")
			}

			found := false
			for i := range batchItems {
				if batchItems[i].OrderID == item.Item().OrderID {
					found = true
				}
			}
			if !found {
				t.Errorf(`item not in initial query returned: %s`, item.Item().OrderID)
			}
		}

		// delete
		bwod := tbl.CreateBatchWriteOperation()
		for i := range batchItems {
			if err = bwod.AddDelete(&batchItems[i]); err != nil {
				t.Error(err.Error())
			}
		}

		if err = bwod.Execute(context.TODO()); err != nil {
			t.Error(err.Error())
		}
	}
}
