package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.HTTPClient = smithyhttp.WrapLogClient(logger{}, aws.NewBuildableHTTPClient(), false)
		o.Retryer = retry.NewStandard()
		//o.DisableAcceptEncodingGzip = true
		//o.DisableValidateResponseChecksum = true
	})

	resp, err := client.Scan(context.Background(),
		&dynamodb.ScanInput{
			TableName:              ptr.String("RepoUsage"),
			ReturnConsumedCapacity: ddbtypes.ReturnConsumedCapacityTotal,
		}, func(o *dynamodb.Options) {
			o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
				fmt.Println("Deserialize:", stack.Deserialize.List())
				return nil
			})
		})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success", len(resp.Items))
}

type logger struct{}

func (logger) Logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
