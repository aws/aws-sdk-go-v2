---
title: "AttachUserPolicyv2"
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

// IAMAttachRolePolicyAPI defines the interface for the AttachRolePolicy function.
// We use this interface to test the function using a mocked service.
type IAMAttachRolePolicyAPI interface {
	AttachRolePolicy(ctx context.Context,
		params *iam.AttachRolePolicyInput,
		optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
}

// AttachDynamoFullPolicy attaches an Amazon DynamoDB full-access policy to an AWS Identity and Access Management (IAM) role.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, an AttachRolePolicyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to AttachRolePolicy.
func AttachDynamoFullPolicy(c context.Context, api IAMAttachRolePolicyAPI, input *iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
	return api.AttachRolePolicy(c, input)
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

	input := &iam.AttachRolePolicyInput{
		PolicyArn: &policyArn,
		RoleName:  roleName,
	}

	_, err = AttachDynamoFullPolicy(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Unable to attach DynamoDB full-access role policy to role")
		return
	}

	fmt.Println("Role attached successfully")
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/AttachUserPolicy/AttachUserPolicyv2.go).