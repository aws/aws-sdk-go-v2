---
title: "AssumeRolev2"
---
## AssumeRolev2.go

This example gets temporary security credentials to access resources.

`go run AssumeRolev2.go -r ROLE-ARN -s SESSION-NAME`

- _ROLE-ARN_ is the ARN of the role to assume.
- _SESSION-NAME_ is the name of the assumed role session.

The unit test accepts similar values in _config.json_.

## Source code

```
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSAssumeRoleAPI defines the interface for the AssumeRole function.
// We use this interface to test the function using a mocked service.
type STSAssumeRoleAPI interface {
	AssumeRole(ctx context.Context,
		params *sts.AssumeRoleInput,
		optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// TakeRole gets temporary security credentials to access resources.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, an AssumeRoleOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to AssumeRole.
func TakeRole(c context.Context, api STSAssumeRoleAPI, input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return api.AssumeRole(c, input)
}

func main() {
	roleARN := flag.String("r", "", "The Amazon Resource Name (ARN) of the role to assume")
	sessionName := flag.String("s", "", "The name of the session")

	if *roleARN == "" || *sessionName == "" {
		fmt.Println("You must supply a role ARN and session name")
		fmt.Println("-r ROLE-ARN -s SESSION-NAME")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sts.NewFromConfig(cfg)

	input := &sts.AssumeRoleInput{
		RoleArn:         roleARN,
		RoleSessionName: sessionName,
	}

	result, err := TakeRole(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error assuming the role:")
		fmt.Println(err)
		return
	}

	fmt.Println(result.AssumedRoleUser)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/sts/AssumeRole/AssumeRolev2.go).