---
title: "DeleteUserv2"
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

// IAMDeleteUserAPI defines the interface for the DeleteUser function.
// We use this interface to test the function using a mocked service.
type IAMDeleteUserAPI interface {
	DeleteUser(ctx context.Context,
		params *iam.DeleteUserInput,
		optFns ...func(*iam.Options)) (*iam.DeleteUserOutput, error)
}

// RemoveUser deletes an AWS Identity and Access Management (IAM) user.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a DeleteUserOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteUser.
func RemoveUser(c context.Context, api IAMDeleteUserAPI, input *iam.DeleteUserInput) (*iam.DeleteUserOutput, error) {
	return api.DeleteUser(c, input)
}

func main() {
	userName := flag.String("u", "", "The name of the user")
	flag.Parse()

	if *userName == "" {
		fmt.Println("You must supply a user name (-u USERNAME)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.DeleteUserInput{
		UserName: userName,
	}

	_, err = RemoveUser(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error deleting user " + *userName)
	}

	fmt.Println("Deleted user " + *userName)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/iam/DeleteUser/DeleteUserv2.go).