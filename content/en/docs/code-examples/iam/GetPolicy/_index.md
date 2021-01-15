---
title: "GetPolicyv2"
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

// IAMGetPolicyAPI defines the interface for the GetPolicy function.
// We use this interface to test the function using a mocked service.
type IAMGetPolicyAPI interface {
	GetPolicy(ctx context.Context,
		params *iam.GetPolicyInput,
		optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
}

// GetPolicyDescription retrieves the description of the AWS Identity and Access Management (IAM) policy with the specified ARN.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a GetPolicyOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetPolicy.
func GetPolicyDescription(c context.Context, api IAMGetPolicyAPI, input *iam.GetPolicyInput) (*iam.GetPolicyOutput, error) {
	return api.GetPolicy(c, input)
}

func main() {
	policyArn := flag.String("a", "", "The ARN of the policy to retrieve")
	flag.Parse()

	if *policyArn == "" {
		fmt.Println("You must supply the ARN of the policy to retrieve (-a POLICY-ARN)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.GetPolicyInput{
		PolicyArn: policyArn,
	}

	result, err := GetPolicyDescription(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving the description:")
		fmt.Println(err)
		return
	}

	description := ""

	if nil == result.Policy {
		description = "Policy nil"
	} else {
		if nil == result.Policy.Description {
			description = "Description nil"
		} else {
			description = *result.Policy.Description
		}
	}

	fmt.Println("Description:")
	fmt.Println(description)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/iam/GetPolicy/GetPolicyv2.go).