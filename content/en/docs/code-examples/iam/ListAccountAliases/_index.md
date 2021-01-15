---
title: "ListAccountAliasesv2"
---
404: Not Found

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMListAccountAliasesAPI defines the interface for the ListAccountAliases function.
// We use this interface to test the function using a mocked service.
type IAMListAccountAliasesAPI interface {
	ListAccountAliases(ctx context.Context,
		params *iam.ListAccountAliasesInput,
		optFns ...func(*iam.Options)) (*iam.ListAccountAliasesOutput, error)
}

// GetAccountAliases retrieves the aliases for your AWS Identity and Access Management (IAM) account.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a ListAccountAliasesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListAccountAliases.
func GetAccountAliases(c context.Context, api IAMListAccountAliasesAPI, input *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	return api.ListAccountAliases(c, input)
}

func main() {
	maxItems := flag.Int("m", 10, "Maximum number of aliases to list")
	flag.Parse()

	if *maxItems < 0 {
		*maxItems = 10
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.ListAccountAliasesInput{
		MaxItems: aws.Int32(int32(*maxItems)),
	}

	result, err := GetAccountAliases(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving account aliases")
		fmt.Println(err)
		return
	}

	for i, alias := range result.AccountAliases {
		fmt.Printf("Alias %d: %s\n", i, alias)
	}
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/iam/ListAccountAliases/ListAccountAliasesv2.go).