---
title: "DeleteAccessKeyv2"
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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMDeleteAccessKeyAPI defines the interface for the DeleteAccessKey function.
// We use this interface to test the function using a mocked service.
type IAMDeleteAccessKeyAPI interface {
	DeleteAccessKey(ctx context.Context,
		params *iam.DeleteAccessKeyInput,
		optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error)
}

// RemoveAccessKey deletes an AWS Identity and Access Management (IAM) access key.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a DeleteAccessKeyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteAccessKey.
func RemoveAccessKey(c context.Context, api IAMDeleteAccessKeyAPI, input *iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return api.DeleteAccessKey(c, input)
}

func main() {
	keyID := flag.String("k", "", "The ID of the access key")
	userName := flag.String("u", "", "The name of the user")
	flag.Parse()

	if *keyID == "" || *userName == "" {
		fmt.Println("You must supply the key ID and user name (-k KEY-ID -u USER-NAME")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: keyID,
		UserName:    userName,
	}

	_, err = RemoveAccessKey(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Deleted key with ID " + *keyID + " from user " + *userName)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/iam/DeleteAccessKey/DeleteAccessKeyv2.go).