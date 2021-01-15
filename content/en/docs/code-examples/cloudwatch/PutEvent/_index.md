---
title: "PutEventv2"
---
## PutEventv2.go

This example sends an Amazon CloudWatch event to Amazon EventBridge.

`go run PutEventv2.go -l LAMBDA-ARN -f EVENT-FILE`

- _LAMBDA-ARN_ is the ARN of the AWS Lambda function of which the event is concerned.
- _EVENT-FILE_ is the local file specifying details of the event to send to Amazon EventBridge.

The unit test accepts similar values in _config.json_.

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

// CWEPutEventsAPI defines the interface for the PutEvents function.
// We use this interface to test the function using a mocked service.
type CWEPutEventsAPI interface {
	PutEvents(ctx context.Context,
		params *cloudwatchevents.PutEventsInput,
		optFns ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

// Event represents the information for a new event
type Event struct {
	Details []struct {
		Key   string `json:"Key"`
		Value string `json:"Value"`
	} `json:"Details"`
	DetailType string `json:"DetailType"`
	Source     string `json:"Source"`
}

func getEventInfo(eventFile string) (Event, error) {
	var e Event

	content, err := ioutil.ReadFile(eventFile)
	if err != nil {
		return e, err
	}

	// Convert []byte to string
	text := string(content)

	// Marshall JSON string in text into global struct
	err = json.Unmarshal([]byte(text), &e)
	if err != nil {
		return e, err
	}

	// Make sure we got the info we need
	if e.DetailType == "" {
		e.DetailType = "appRequestSubmitted"
	}

	if e.Source == "" {
		e.Source = "com.mycompany.myapp"
	}

	if e.Details == nil {
		d := []byte(`"{ "key1": "value1", "key2": "value2" }`)
		e.DetailType = string(d[:])
	}

	return e, nil
}

// CreateEvent sends an Amazon CloudWatch event to Amazon EventBridge
// Inputs:
//     c is the context of the method call, which includes the Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutEventsOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutEvents
func CreateEvent(c context.Context, api CWEPutEventsAPI, input *cloudwatchevents.PutEventsInput) (*cloudwatchevents.PutEventsOutput, error) {
	return api.PutEvents(c, input)
}

func main() {
	lambdaARN := flag.String("l", "", "The ARN of the AWS Lambda function")
	eventFile := flag.String("f", "", "The JSON file containing the event to send")
	flag.Parse()

	if *lambdaARN == "" || *eventFile == "" {
		fmt.Println("You must supply a Lambda ARN with -l LAMBDA-ARN and event file with -f EVENT-FILE.json")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := cloudwatchevents.NewFromConfig(cfg)

	event, err := getEventInfo(*eventFile)
	if err != nil {
		fmt.Println("Got an error calling getEventInfo:")
		fmt.Println(err)
		return
	}

	myDetails := "{ "
	for _, d := range event.Details {
		myDetails = myDetails + "\"" + d.Key + "\": \"" + d.Value + "\","
	}

	myDetails = myDetails + " }"

	input := &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:     &myDetails,
				DetailType: &event.DetailType,
				Resources: []string{
					*lambdaARN,
				},
				Source: &event.Source,
			},
		},
	}

	_, err = CreateEvent(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Could not create event:")
		fmt.Println(err)
		return
	}

	fmt.Println("Created event")
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/cloudwatch/PutEvent/PutEventv2.go).