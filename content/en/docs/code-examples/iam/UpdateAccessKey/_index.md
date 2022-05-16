---
title: "UpdateAccessKeyv2"
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
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// IAMUpdateAccessKeyAPI defines the interface for the UpdateAccessKey function.
// We use this interface to test the function using a mocked service.
type IAMUpdateAccessKeyAPI interface {
	UpdateAccessKey(ctx context.Context,
		params *iam.UpdateAccessKeyInput,
		optFns ...func(*iam.Options)) (*iam.UpdateAccessKeyOutput, error)
}

// ActivateKey sets the status of an AWS Identity and Access Management (IAM) access key to active.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a UpdateAccessKeyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to UpdateAccessKey.
func ActivateKey(c context.Context, api IAMUpdateAccessKeyAPI, input *iam.UpdateAccessKeyInput) (*iam.UpdateAccessKeyOutput, error) {
	return api.UpdateAccessKey(c, input)
}

func main() {
	keyID := flag.String("k", "", "The ID of the access key")
	userName := flag.String("u", "", "The name of the user")
	flag.Parse()

	if *keyID == "" || *userName == "" {
		fmt.Println("You must supply an access key ID and user name (-k KEY-ID -u USER-NAME)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.UpdateAccessKeyInput{
		AccessKeyId: keyID,
		Status:      types.StatusTypeActive,
		UserName:    userName,
	}

	_, err = ActivateKey(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Access Key activated")
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/UpdateAccessKey/UpdateAccessKeyv2.go).