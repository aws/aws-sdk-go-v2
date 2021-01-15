---
title: "ListTopicsv2"
---
## ListTopicsv2.go

This example lists the ARNs for your Amazon SNS topics.

`go run ListTopicsv2.go`

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSListTopicsAPI defines the interface for the ListTopics function.
// We use this interface to test the function using a mocked service.
type SNSListTopicsAPI interface {
	ListTopics(ctx context.Context,
		params *sns.ListTopicsInput,
		optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
}

// GetTopics retrieves information about the Amazon Simple Notification Service (Amazon SNS) topics
// Inputs:
//     c is the context of the method call, which includes the Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListTopicsOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to ListTopics
func GetTopics(c context.Context, api SNSListTopicsAPI, input *sns.ListTopicsInput) (*sns.ListTopicsOutput, error) {
	return api.ListTopics(c, input)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListTopicsInput{}

	results, err := GetTopics(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about the SNS topics:")
		fmt.Println(err)
		return
	}

	for _, t := range results.Topics {
		fmt.Println(*t.TopicArn)
	}
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/sns/ListTopics/ListTopicsv2.go).