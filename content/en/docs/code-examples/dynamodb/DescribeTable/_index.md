---
title: "DescribeTablev2"
---
## DescribeTablev2.go

This example lists the following properties of a DynamoDB table.

- Number of items
- Size, in bytes
- Status, such as Active

`go run DescribeTablev2.go -t TABLE`

- _TABLE_ is the name of the table.

The unit test accepts a similar value in _config.json_.

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBDescribeTableAPI defines the interface for the DescribeTable function.
// We use this interface to enable unit testing.
type DynamoDBDescribeTableAPI interface {
	DescribeTable(ctx context.Context,
		params *dynamodb.DescribeTableInput,
		optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

// GetTableInfo retrieves information about the table.
func GetTableInfo(c context.Context, api DynamoDBDescribeTableAPI, input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	return api.DescribeTable(c, input)
}

func main() {
	table := flag.String("t", "", "The name of the table")
	flag.Parse()

	if *table == "" {
		fmt.Println("You must specify a table name (-t TABLE)")
		return
	}

	// Use the SDK's default configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Create an Amazon DynamoDB client.
	client := dynamodb.NewFromConfig(cfg)

	// Build the input parameters for the request.
	input := &dynamodb.DescribeTableInput{
		TableName: table,
	}

	resp, err := GetTableInfo(context.TODO(), client, input)
	if err != nil {
		panic("failed to describe table, " + err.Error())
	}

	fmt.Println("Info about " + *table + ":")
	fmt.Println("  #items:     ", resp.Table.ItemCount)
	fmt.Println("  Size (bytes)", resp.Table.TableSizeBytes)
	fmt.Println("  Status:     ", string(resp.Table.TableStatus))
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/DescribeTable/DescribeTablev2.go).