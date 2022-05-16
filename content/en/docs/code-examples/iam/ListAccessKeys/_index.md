---
title: "ListAccessKeysv2"
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMListAccessKeysAPI defines the interface for the ListAccessKeys function.
// We use this interface to test the function using a mocked service.
type IAMListAccessKeysAPI interface {
	ListAccessKeys(ctx context.Context,
		params *iam.ListAccessKeysInput,
		optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error)
}

//  GetAccessKeys retrieves up to the AWS Identity and Access Management (IAM) access keys for a user.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a ListAccessKeysOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListAccessKeys.
func GetAccessKeys(c context.Context, api IAMListAccessKeysAPI, input *iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return api.ListAccessKeys(c, input)
}

func main() {
	maxItems := flag.Int("m", 10, "The maximum number of access keys to show")
	userName := flag.String("u", "", "The name of the user")
	flag.Parse()

	if *userName == "" {
		fmt.Println("You must supply the name of a user (-u USER)")
		return
	}

	if *maxItems < 0 {
		*maxItems = 10
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.ListAccessKeysInput{
		MaxItems: aws.Int32(int32(*maxItems)),
		UserName: userName,
	}

	result, err := GetAccessKeys(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving user access keys:")
		fmt.Println(err)
		return
	}

	for _, key := range result.AccessKeyMetadata {
		fmt.Println("Status for access key " + *key.AccessKeyId + ": " + string(key.Status))
	}
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/ListAccessKeys/ListAccessKeysv2.go).