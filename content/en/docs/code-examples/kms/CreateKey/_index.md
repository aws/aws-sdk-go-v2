---
title: "CreateKeyv2"
---
## CreateKeyv2.go

This example creates an AWS Key Management Service (AWS KMS) customer master key (CMK).

`go run CreateKeyv2.go`

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSCreateKeyAPI defines the interface for the CreateKey function.
// We use this interface to test the function using a mocked service.
type KMSCreateKeyAPI interface {
	CreateKey(ctx context.Context,
		params *kms.CreateKeyInput,
		optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
}

// MakeKey creates an AWS Key Management Service (AWS KMS) customer master key (CMK).
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a CreateKeyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to CreateKey.
func MakeKey(c context.Context, api KMSCreateKeyAPI, input *kms.CreateKeyInput) (*kms.CreateKeyOutput, error) {
	return api.CreateKey(c, input)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := kms.NewFromConfig(cfg)

	input := &kms.CreateKeyInput{}

	result, err := MakeKey(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error creating key:")
		fmt.Println(err)
		return
	}

	fmt.Println(*result.KeyMetadata.KeyId)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/kms/CreateKey/CreateKeyv2.go).