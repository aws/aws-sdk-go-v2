// +build integration

package customizations_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

var (
	tableName string
)

func init() {
	flag.StringVar(&tableName, "table", "testTable",
		"The `name` of the table to test against")
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestInteg_ClientScan(t *testing.T) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.HTTPClient = smithyhttp.WrapLogClient(t, aws.NewBuildableHTTPClient(), false)
		o.Retryer = retry.NewStandard()

		o.EnableAcceptEncodingGzip = true
		o.DisableValidateResponseChecksum = false // default
	})

	_, err = client.Scan(context.Background(),
		&dynamodb.ScanInput{
			TableName:              &tableName,
			ReturnConsumedCapacity: ddbtypes.ReturnConsumedCapacityTotal,
		}, func(o *dynamodb.Options) {
			o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
				fmt.Println("Stack:", stack.String())
				return nil
			})
		})
	if err != nil {
		t.Fatal(err)
	}
}
