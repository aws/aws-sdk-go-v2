// +build example

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jviney/aws-sdk-go-v2/aws/external"
	"github.com/jviney/aws-sdk-go-v2/service/sagemaker"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// This code serves an example on how to use it for sagemaker
//
// Usage: go run -tags example listTrainingJobs <int_value>
func main() {
	if len(os.Args) < 2 {
		exitErrorf("you must specify a MaxItems")
	}

	x, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		exitErrorf("failed to parse argument %v", err)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	cfg.Region = "us-west-2"

	sagemakerSvc := sagemaker.New(cfg)

	req := sagemakerSvc.ListTrainingJobsRequest(&sagemaker.ListTrainingJobsInput{MaxResults: &x})
	resp, err := req.Send(context.TODO())
	if err != nil {
		exitErrorf("failed to list training jobs, %v", err)
	}

	fmt.Println(resp)
}
