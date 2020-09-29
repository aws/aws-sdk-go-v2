// +build integration

package snowball

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/snowball"
)

func TestInteg_00_DescribeAddresses(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := integrationtest.LoadConfigWithDefaultRegion("us-west-2")
	if err != nil {
		t.Fatalf("failed to load config, %v", err)
	}

	client := snowball.NewFromConfig(cfg)
	params := &snowball.DescribeAddressesInput{}
	_, err = client.DescribeAddresses(ctx, params)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}
