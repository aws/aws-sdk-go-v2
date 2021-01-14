---
title: "Subscribev2"
---
## Subscribev2.go

This example subscribes a user, by email address, to an Amazon SNS topic.

`go run Subscribev2.go -m EMAIL-ADDRESS -t TOPIC-ARN`

- _EMAIL-ADDRESS_ is the email address of the user subscribing to the topic.
- _TOPIC-ARN_ is the ARN of the topic.

The unit test accepts a similar value in _config.json_.

## Source code

```
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSSubscribeAPI defines the interface for the Subscribe function.
// We use this interface to test the function using a mocked service.
type SNSSubscribeAPI interface {
	Subscribe(ctx context.Context,
		params *sns.SubscribeInput,
		optFns ...func(*sns.Options)) (*sns.SubscribeOutput, error)
}

// SubscribeTopic subscribes a user to an Amazon Simple Notification Service (Amazon SNS) topic by their email address
// Inputs:
//     c is the context of the method call, which includes the Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a SubscribeOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to Subscribe
func SubscribeTopic(c context.Context, api SNSSubscribeAPI, input *sns.SubscribeInput) (*sns.SubscribeOutput, error) {
	return api.Subscribe(c, input)
}

func main() {
	email := flag.String("e", "", "The email address of the user subscribing to the topic")
	topicARN := flag.String("t", "", "The ARN of the topic to which the user subscribes")

	flag.Parse()

	if *email == "" || *topicARN == "" {
		fmt.Println("You must supply an email address and topic ARN")
		fmt.Println("-e EMAIL -t TOPIC-ARN")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.SubscribeInput{
		Endpoint:              email,
		Protocol:              aws.String("email"),
		ReturnSubscriptionArn: true, // Return the ARN, even if user has yet to confirm
		TopicArn:              topicARN,
	}

	result, err := SubscribeTopic(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error subscribing to the topic:")
		fmt.Println(err)
		return
	}

	fmt.Println(*result.SubscriptionArn)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/sns/Subscribe/Subscribev2.go).