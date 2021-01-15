---
title: "CreateTopicv2"
---
## CreateTopicv2.go

This example creates an Amazon SNS topic.

`go run CreateTopicv2.go -t TOPIC`

- _TOPIC_ is the name of the topic to create.

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
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSCreateTopicAPI defines the interface for the CreateTopic function.
// We use this interface to test the function using a mocked service.
type SNSCreateTopicAPI interface {
	CreateTopic(ctx context.Context,
		params *sns.CreateTopicInput,
		optFns ...func(*sns.Options)) (*sns.CreateTopicOutput, error)
}

// MakeTopic creates an Amazon Simple Notification Service (Amazon SNS) topic.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateTopicOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateTopic.
func MakeTopic(c context.Context, api SNSCreateTopicAPI, input *sns.CreateTopicInput) (*sns.CreateTopicOutput, error) {
	return api.CreateTopic(c, input)
}

func main() {
	topic := flag.String("t", "", "The name of the topic")
	flag.Parse()

	if *topic == "" {
		fmt.Println("You must supply the name of the topic")
		fmt.Println("-t TOPIC")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.CreateTopicInput{
		Name: topic,
	}

	results, err := MakeTopic(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(*results.TopicArn)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/sns/CreateTopic/CreateTopicv2.go).