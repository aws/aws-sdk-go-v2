package enhancedclient

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestTableE2E(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	c := dynamodb.NewFromConfig(cfg)

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

	<-time.After(time.Second * 5)

	_, err = tbl.Describe(context.Background())
	if err != nil {
		t.Fatalf("Error describing table: %v", err)
	}

	itemsToManage := 10
	sk := time.Now().Unix()
	t.Logf("Will now write %d items", itemsToManage)
	for i := range itemsToManage {
		o := order{
			OrderID:      fmt.Sprintf("%d", i),
			CreatedAt:    sk,
			CustomerID:   fmt.Sprintf("CustomerID%d", i),
			TotalAmount:  float64(i),
			IgnoredField: fmt.Sprintf("IgnoredField%d", i),
			Version:      int64(i),
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

		err = tbl.PutItem(context.Background(), &o)
		if err != nil {
			t.Logf("Unable to PutItem() [%d]: %v", i, err)
		}
	}
	t.Logf("Will now fetch %d items", itemsToManage)
	for i := range itemsToManage {
		m := Map{}.
			With("order_id", fmt.Sprintf("%d", i)).
			With("created_at", sk)
		item, err := tbl.GetItem(context.Background(), m)
		if err != nil {
			t.Logf("Unable to GetItem() [%d]: %v", i, err)
		}
		t.Logf("Item: %s - %d", item.OrderID, item.CreatedAt)
	}
	t.Logf("Will now delete %d items", itemsToManage)
	for i := range itemsToManage {
		if i&1 == 0 {
			m := Map{}.
				With("order_id", fmt.Sprintf("%d", i)).
				With("created_at", sk)
			err = tbl.DeleteItemByKey(context.Background(), m)
		} else {
			err = tbl.DeleteItem(context.Background(), &order{
				OrderID:   fmt.Sprintf("%d", i),
				CreatedAt: sk,
			})
		}
		if err != nil {
			t.Logf("Unable to GetItem() [%d]: %v", i, err)
		}
	}

	// delete
	t.Logf("Table %s will be deleted", tableName)
	err = tbl.DeleteWithWait(context.Background(), time.Minute)
	if err != nil {
		t.Fatalf("DeleteWithWait() error: %v", err)
	}
	t.Logf("Table %s deleted", tableName)
}
