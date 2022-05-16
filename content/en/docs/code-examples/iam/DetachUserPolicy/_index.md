---
title: "DetachUserPolicyv2"
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

// IAMDetachRolePolicyAPI defines the interface for the DetachRolePolicy function.
// We use this interface to test the function using a mocked service.
type IAMDetachRolePolicyAPI interface {
	DetachRolePolicy(ctx context.Context,
		params *iam.DetachRolePolicyInput,
		optFns ...func(*iam.Options)) (*iam.DetachRolePolicyOutput, error)
}

// DetachDynamoFullPolicy detaches an Amazon DynamoDB full-access policy from an AWS Identity and Access Management (IAM) role.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a DetachRolePolicyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DetachRolePolicy.
func DetachDynamoFullPolicy(c context.Context, api IAMDetachRolePolicyAPI, input *iam.DetachRolePolicyInput) (*iam.DetachRolePolicyOutput, error) {
	return api.DetachRolePolicy(c, input)
}

func main() {
	roleName := flag.String("r", "", "The name of the IAM role")
	flag.Parse()

	if *roleName == "" {
		fmt.Println("You must supply a role name (-r ROLE)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	policyArn := "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
	input := &iam.DetachRolePolicyInput{
		PolicyArn: &policyArn,
		RoleName:  roleName,
	}

	_, err = DetachDynamoFullPolicy(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Unable to detach DynamoDB full-access role policy from role")
		return
	}
	fmt.Println("Role detached successfully")
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/DetachUserPolicy/DetachUserPolicyv2.go).