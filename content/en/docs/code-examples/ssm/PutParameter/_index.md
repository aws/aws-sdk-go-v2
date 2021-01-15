---
title: "PutParameterv2"
---
## PutParameterv2.go

This example creates a Systems Manager string parameter.

`go run PutParameterv2.go -n NAME -v VALUE`

- _NAME_ is the name of the parameter to create.
- _VALUE_ is the value of the parameter to create.

The unit test accepts similar values in _config.json_.

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
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SSMPutParameterAPI defines the interface for the PutParameter function.
// We use this interface to test the function using a mocked service.
type SSMPutParameterAPI interface {
	PutParameter(ctx context.Context,
		params *ssm.PutParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
}

// AddStringParameter creates an AWS Systems Manager string parameter
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutParameterOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutParameter
func AddStringParameter(c context.Context, api SSMPutParameterAPI, input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return api.PutParameter(c, input)
}

func main() {
	parameterName := flag.String("n", "", "The name of the parameter")
	parameterValue := flag.String("v", "", "The value of the parameter")
	flag.Parse()

	if *parameterName == "" {
		fmt.Println("You must supply the name of the parameter")
		fmt.Println("-n NAME")
		return
	}

	if *parameterValue == "" {
		fmt.Println("You must supply the value of the parameter")
		fmt.Println("-v VALUE")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.PutParameterInput{
		Name:  parameterName,
		Value: parameterValue,
		Type:  types.ParameterTypeString,
	}

	results, err := AddStringParameter(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Parameter version:", results.Version)
}
```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/ssm/PutParameter/PutParameterv2.go).