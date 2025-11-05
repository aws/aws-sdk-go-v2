package enhancedclient

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dummyExtension struct{}

func (*dummyExtension) IsExtension() {}

func (*dummyExtension) BeforeRead(context.Context, *dummyExtension) error  { return nil }
func (*dummyExtension) AfterRead(context.Context, *dummyExtension) error   { return nil }
func (*dummyExtension) BeforeWrite(context.Context, *dummyExtension) error { return nil }
func (*dummyExtension) AfterWrite(context.Context, *dummyExtension) error  { return nil }
func (*dummyExtension) BeforeQuery(context.Context, *dynamodb.QueryInput) error {
	return nil
}
func (*dummyExtension) AfterQuery(context.Context, []dummyExtension) error    { return nil }
func (*dummyExtension) BeforeScan(context.Context, *dynamodb.ScanInput) error { return nil }
func (*dummyExtension) AfterScan(context.Context, []dummyExtension) error     { return nil }

func TestExtensionRegistry(t *testing.T) {
	er := &ExtensionRegistry[dummyExtension]{}
	er.AddBeforeReader(&dummyExtension{})
	er.AddAfterReader(&dummyExtension{})
	er.AddBeforeWriter(&dummyExtension{})
	er.AddAfterWriter(&dummyExtension{})
	//er.AddBeforeScanner(&dummyExtension{})
	//er.AddAfterScanner(&dummyExtension{})
	//er.AddBeforeQuerier(&dummyExtension{})
	//er.AddAfterQuerier(&dummyExtension{})

	if len(er.beforeReaders) != 1 {
		t.Errorf("beforeReaders expected to be 1, got %d", len(er.beforeReaders))
	}
	if len(er.afterReaders) != 1 {
		t.Errorf("afterReaders expected to be 1, got %d", len(er.afterReaders))
	}
	if len(er.beforeWriters) != 1 {
		t.Errorf("beforeWriters expected to be 1, got %d", len(er.beforeWriters))
	}
	if len(er.afterWriters) != 1 {
		t.Errorf("afterWriters expected to be 1, got %d", len(er.afterWriters))
	}
	//if len(er.beforeQueriers) != 1 {
	//	t.Errorf("beforeQueriers expected to be 1, got %d", len(er.beforeQueriers))
	//}
	//if len(er.afterQueriers) != 1 {
	//	t.Errorf("afterQueriers expected to be 1, got %d", len(er.afterQueriers))
	//}
	//if len(er.beforeScanners) != 1 {
	//	t.Errorf("beforeScanners expected to be 1, got %d", len(er.beforeScanners))
	//}
	//if len(er.afterScanners) != 1 {
	//	t.Errorf("afterScanners expected to be 1, got %d", len(er.afterScanners))
	//}
}
