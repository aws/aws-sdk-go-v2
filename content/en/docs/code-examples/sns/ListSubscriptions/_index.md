---
title: "ListSubscriptionsv2"
---
## ListSubscriptionsv2.go

This example lists the topic and subscription Amazon Resource Names (ARNs) for your Amazon SNS subscriptions.

`go run ListSubscriptionsv2.go`

## Source code

```
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSListSubscriptionsAPI defines the interface for the ListSubscriptions function.
// We use this interface to test the function using a mocked service.
type SNSListSubscriptionsAPI interface {
	ListSubscriptions(ctx context.Context,
		params *sns.ListSubscriptionsInput,
		optFns ...func(*sns.Options)) (*sns.ListSubscriptionsOutput, error)
}

// GetSubscriptions retrieves a list of your Amazon Simple Notification Service (Amazon SNS) subscriptions.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListSubscriptionsOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListSubscriptions.
func GetSubscriptions(c context.Context, api SNSListSubscriptionsAPI, input *sns.ListSubscriptionsInput) (*sns.ListSubscriptionsOutput, error) {
	return api.ListSubscriptions(c, input)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.ListSubscriptionsInput{}

	result, err := GetSubscriptions(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving the subscriptions:")
		fmt.Println(err)
		return
	}

	fmt.Println("Topic ARN")
	fmt.Println("Subscription ARN")
	fmt.Println("-------------------------")
	for _, s := range result.Subscriptions {
		fmt.Println(*s.TopicArn)
		fmt.Println(*s.SubscriptionArn)
		fmt.Println("")
	}
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/sns/ListSubscriptions/ListSubscriptionsv2.go).