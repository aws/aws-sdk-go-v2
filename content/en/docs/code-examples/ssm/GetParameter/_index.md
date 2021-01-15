---
title: "GetParameterv2"
---
## GetParameterv2.go

This example retrieves a Systems Manager string parameter.

`go run GetParameterv2.go -n NAME`

- _NAME_ is the name of the parameter to retrieve.

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

// SSMGetParameterAPI defines the interface for the GetParameter function.
// We use this interface to test the function using a mocked service.
type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// FindParameter retrieves an AWS Systems Manager string parameter
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetParameterOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to GetParameter
func FindParameter(c context.Context, api SSMGetParameterAPI, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return api.GetParameter(c, input)


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

	input := &ssm.GetParameterInput{
		Name: parameterName,
	}

	results, err := FindParameter(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(*results.Parameter.Value)
}
```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/ssm/GetParameter/GetParameterv2.go).