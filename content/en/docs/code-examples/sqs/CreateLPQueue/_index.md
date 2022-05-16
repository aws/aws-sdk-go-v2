---
title: "CreateLPQueuev2"
---
## CreateLPQueuev2.go

This example creates a long-polling Amazon SQS queue.

`go run CreateLPQueuev2.go -q QUEUE-NAME [-w WAIT-TIME]`

- _QUEUE-NAME_ is the name of the queue to create.
- _WAIT-TIME_ is how long, in seconds, to wait.
  The example ensures the value is between 1 and 20;
  the default is 10.

The unit test accepts similar values in _config.json_.

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSCreateQueueAPI defines the interface for the CreateQueue function.
// We use this interface to test the function using a mocked service.
type SQSCreateQueueAPI interface {
	CreateQueue(ctx context.Context,
		params *sqs.CreateQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.CreateQueueOutput, error)
}

// CreateLPQueue creates an Amazon SQS queue with long polling enabled.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateQueueOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateQueue.
func CreateLPQueue(c context.Context, api SQSCreateQueueAPI, input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	return api.CreateQueue(c, input)
}

func main() {
	queue := flag.String("q", "", "The name of the queue")
	waitTime := flag.Int("w", 10, "How long, in seconds, to wait for long polling")
	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply a queue name (-q QUEUE")
		return
	}

	if *waitTime < 1 {
		*waitTime = 1
	}

	if *waitTime > 20 {
		*waitTime = 20
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.CreateQueueInput{
		QueueName: queue,
		Attributes: map[string]string{
			"ReceiveMessageWaitTimeSeconds": strconv.Itoa(*waitTime),
		},
	}

	resp, err := CreateLPQueue(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error creating the long polling queue:")
		fmt.Println(err)
		return
	}

	fmt.Println("URL for long polling queue " + *queue + ": " + *resp.QueueUrl)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/sqs/CreateLPQueue/CreateLPQueuev2.go).