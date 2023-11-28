package readcloserwithmetrics

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/testutils"
	"io"
	"testing"
)

func TestNew(t *testing.T) {

	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, &testutils.NoopPublisher{})
	mctx := metrics.Context(ctx)

	expectedData := "testString"
	trc := testutils.TestReadCloser{Data: []byte(expectedData)}

	reader := New(mctx, &trc)
	readData, _ := io.ReadAll(reader)
	actualData := string(readData)

	if actualData != expectedData {
		t.Errorf("Unexpected Data, should be '%s' but was '%s'", expectedData, actualData)
	}
}

func TestReadCloserWithMetrics_Close(t *testing.T) {

	mdrp := &testutils.MetricDataRecorderPublisher{}
	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, mdrp)
	mctx := metrics.Context(ctx)

	expectedData := "testString"
	trc := testutils.TestReadCloser{Data: []byte(expectedData)}

	reader := New(mctx, &trc)

	err := reader.Close()

	if err != nil {
		t.Errorf("Unexpected Error in Close")
	}

	if mdrp.Data == nil {
		t.Errorf("Data should be set in publisher")
	}

	rb := mdrp.Data.Stream.ReadBytes
	if mdrp.Data.Stream.ReadBytes != 0 {
		t.Errorf("Unexpected ReadBytes, should be '%d' but was '%d'", 0, rb)
	}

}

func TestReadCloserWithMetrics_Read(t *testing.T) {
	mdrp := &testutils.MetricDataRecorderPublisher{}
	ctx := context.TODO()
	ctx = metrics.InitMetricContext(ctx, &metrics.SharedConnectionCounter{}, mdrp)
	mctx := metrics.Context(ctx)

	expectedData := "testString"
	trc := testutils.TestReadCloser{Data: []byte(expectedData)}

	reader := New(mctx, &trc)

	err := reader.Close()

	if err != nil {
		t.Errorf("Unexpected Error in Close")
	}

	if mdrp.Data == nil {
		t.Errorf("Data should be set in publisher")
	}

	rb := mdrp.Data.Stream.ReadBytes
	if mdrp.Data.Stream.ReadBytes != 0 {
		t.Errorf("Unexpected ReadBytes, should be '%d' but was '%d'", 0, rb)
	}
}
