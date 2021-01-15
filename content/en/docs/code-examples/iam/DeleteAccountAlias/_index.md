---
title: "DeleteAccountAliasv2"
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

// IAMDeleteAccountAliasAPI defines the interface for the DeleteAccountAlias function.
// We use this interface to test the function using a mocked service.
type IAMDeleteAccountAliasAPI interface {
	DeleteAccountAlias(ctx context.Context,
		params *iam.DeleteAccountAliasInput,
		optFns ...func(*iam.Options)) (*iam.DeleteAccountAliasOutput, error)
}

// RemoveAccountAlias deletes an alias for your AWS Identity and Access Management (IAM) account.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a DeleteAccountAliasOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteAccountAlias.
func RemoveAccountAlias(c context.Context, api IAMDeleteAccountAliasAPI, input *iam.DeleteAccountAliasInput) (*iam.DeleteAccountAliasOutput, error) {
	return api.DeleteAccountAlias(c, input)
}

func main() {
	alias := flag.String("a", "", "The account alias")
	flag.Parse()

	if *alias == "" {
		fmt.Println("You must supply an account alias (-a ALIAS)")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.DeleteAccountAliasInput{
		AccountAlias: alias,
	}

	_, err = RemoveAccountAlias(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error deleting an account alias")
		fmt.Println(err)
		return
	}

	fmt.Printf("Deleted account alias " + *alias)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/iam/DeleteAccountAlias/DeleteAccountAliasv2.go).