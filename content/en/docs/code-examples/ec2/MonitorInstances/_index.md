---
title: "MonitorInstancesv2"
---
## MonitorInstancesv2.go

This example enables or disables monitoring for an Amazon EC2 instance.

`go run MonitorInstancesv2.go -m MODE -i INSTANCE-ID`

- _MODE_ is either "OFF" to disable monitoring or "ON" to enable monitoring.
- _INSTANCE-ID_ is the ID of the instance.

The unit test accepts similar values in _config.json_.

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
)

// EC2MonitorInstancesAPI defines the interface for the MonitorInstances and UnmonitorInstances functions.
// We use this interface to test the function using a mocked service.
type EC2MonitorInstancesAPI interface {
	MonitorInstances(ctx context.Context,
		params *ec2.MonitorInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.MonitorInstancesOutput, error)

	UnmonitorInstances(ctx context.Context,
		params *ec2.UnmonitorInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.UnmonitorInstancesOutput, error)
}

// EnableMonitoring enables monitoring for an Amazon EC2 instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a MonitorInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to MonitorInstances.
func EnableMonitoring(c context.Context, api EC2MonitorInstancesAPI, input *ec2.MonitorInstancesInput) (*ec2.MonitorInstancesOutput, error) {
	resp, err := api.MonitorInstances(c, input)

	// Do we have a DryRunOperation error?
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
		fmt.Println("User has permission to enable monitoring.")
		input.DryRun = false
		return api.MonitorInstances(c, input)
	}

	return resp, err
}

// DisableMonitoring disables monitoring for an Amazon EC2 instance.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a UnmonitorInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to UnmonitorInstances.
func DisableMonitoring(c context.Context, api EC2MonitorInstancesAPI, input *ec2.UnmonitorInstancesInput) (*ec2.UnmonitorInstancesOutput, error) {
	resp, err := api.UnmonitorInstances(c, input)

	// Do we have a DryRunOperation error?
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
		fmt.Println("User has permission to disable monitoring.")
		input.DryRun = false
		return api.UnmonitorInstances(c, input)
	}

	return resp, err
}

func main() {
	monitor := flag.String("m", "", "ON to enable monitoring; OFF to disable monitoring")
	instanceID := flag.String("i", "", "The ID of the instance to monitor")
	flag.Parse()

	fmt.Println("Monitor:    " + *monitor)
	fmt.Println("InstanceID: " + *instanceID)

	if *instanceID == "" || (*monitor != "ON" && *monitor != "OFF") {
		fmt.Println("You must supply the ID of the instance to enable/disable monitoring (-i INSTANCE-ID)")
		fmt.Println("and whether to enable monitoring (-m ON) or disable monitoring (-m OFF)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	if *monitor == "ON" {
		input := &ec2.MonitorInstancesInput{
			InstanceIds: []string{
				*instanceID,
			},
			DryRun: true,
		}

		result, err := EnableMonitoring(context.TODO(), client, input)
		if err != nil {
			fmt.Println("Got an error enabling monitoring for instance:")
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully enabled monitoring for instance: " + *result.InstanceMonitorings[0].InstanceId)
	} else if *monitor == "OFF" {
		input := &ec2.UnmonitorInstancesInput{
			InstanceIds: []string{
				*instanceID,
			},
			DryRun: true,
		}

		result, err := DisableMonitoring(context.TODO(), client, input)
		if err != nil {
			fmt.Println("Got an error disabling monitoring for instance:")
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully disabled monitoring for instance: " + *result.InstanceMonitorings[0].InstanceId)
	}
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/ec2/MonitorInstances/MonitorInstancesv2.go).