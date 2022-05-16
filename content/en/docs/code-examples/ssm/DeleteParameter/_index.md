---
title: "DeleteParameterv2"
---
## DeleteParameterv2.go

This example deletes a Systems Manager string parameter.

`go run DeleteParameterv2.go -n NAME`

- _NAME_ is the name of the parameter to delete.

The unit test accepts a similar value in _config.json_.

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
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMDeleteParameterAPI defines the interface for the DeleteParameter function.
// We use this interface to test the function using a mocked service.
type SSMDeleteParameterAPI interface {
	DeleteParameter(ctx context.Context,
		params *ssm.DeleteParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
}

// RemoveParameter deletes an AWS Systems Manager string parameter.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteParameterOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteParameter.
func RemoveParameter(c context.Context, api SSMDeleteParameterAPI, input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	return api.DeleteParameter(c, input)
}

func main() {
	parameterName := flag.String("n", "", "The name of the parameter")
	flag.Parse()

	if *parameterName == "" {
		fmt.Println("You must supply the name of the parameter")
		fmt.Println("-n NAME")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.DeleteParameterInput{
		Name: parameterName,
	}

	_, err = RemoveParameter(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Deleted parameter " + *parameterName)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/ssm/DeleteParameter/DeleteParameterv2.go).