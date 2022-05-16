---
title: "DeleteQueuev2"
---
## DeleteQueuev2.go

This example deletes an Amazon SQS queue.

`go run DeleteQueuev2.go -q QUEUE-NAME`

- _QUEUE-NAME_ is the name of the queue to delete.

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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSDeleteQueueAPI defines the interface for the GetQueueUrl and DeleteQueue functions.
// We use this interface to test the functions using a mocked service.
type SQSDeleteQueueAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteQueue(ctx context.Context,
		params *sqs.DeleteQueueInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteQueueOutput, error)
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetQueueUrl.
func GetQueueURL(c context.Context, api SQSDeleteQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// DeleteQueue deletes an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteQueueOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteQueue.
func DeleteQueue(c context.Context, api SQSDeleteQueueAPI, input *sqs.DeleteQueueInput) (*sqs.DeleteQueueOutput, error) {
	return api.DeleteQueue(c, input)
}

func main() {
	queue := flag.String("q", "", "The name of the queue")
	flag.Parse()

	if *queue == "" {
		fmt.Println("You must supply a queue name (-q QUEUE")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	qInput := &sqs.GetQueueUrlInput{
		QueueName: queue,
	}

	// Get the URL for the queue
	result, err := GetQueueURL(context.TODO(), client, qInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	dqInput := &sqs.DeleteQueueInput{
		QueueUrl: queueURL,
	}

	_, err = DeleteQueue(context.TODO(), client, dqInput)
	if err != nil {
		fmt.Println("Got an error deleting the queue:")
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted queue with URL " + *queueURL)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/sqs/DeleteQueue/DeleteQueuev2.go).