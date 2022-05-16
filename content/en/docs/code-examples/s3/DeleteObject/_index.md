---
title: "DeleteObjectv2"
---
## DeleteObjectv2.go

This example deletes an Amazon S3 bucket object.

`go run DeleteObjectv2.go -b BUCKET -o OBJECT`

- _BUCKET_ is the name of the bucket containing the item to delete.
- _OBJECT_ is the name of the object to delete.

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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3DeleteObjectAPI defines the interface for the DeleteObject function.
// We use this interface to test the function using a mocked service.
type S3DeleteObjectAPI interface {
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// DeleteItem deletes an object from an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteObjectOutput object containing the result of the service call and nil
//     Otherwise, an error from the call to DeleteObject
func DeleteItem(c context.Context, api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return api.DeleteObject(c, input)
}

func main() {
	bucket := flag.String("b", "", "The bucket from which the object is deleted")
	objectName := flag.String("o", "", "The object to delete")
	flag.Parse()

	if *bucket == "" || *objectName == "" {
		fmt.Println("You must supply the bucket (-b BUCKET), and object to delete (-o OBJECT")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    objectName,
	}

	_, err = DeleteItem(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error deleting item:")
		fmt.Println(err)
		return
	}

	fmt.Println("Deleted " + *objectName + " from " + *bucket)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/s3/DeleteObject/DeleteObjectv2.go).