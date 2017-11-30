// +build example

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// This example will list instances with a filter
//
// Usage:
// filter_ec2_by_tag <name_filter>
func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}

	nameFilter := os.Args[1]
	awsRegion := "us-east-1"

	cfg.Region = awsRegion
	svc := ec2.New(cfg)

	fmt.Printf("listing instances with tag %v in: %v\n", nameFilter, awsRegion)
	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					strings.Join([]string{"*", nameFilter, "*"}, ""),
				},
			},
		},
	}

	req := svc.DescribeInstancesRequest(params)
	resp, err := req.Send()
	if err != nil {
		exitErrorf("failed to describe instances, %s, %v", awsRegion, err)
	}

	fmt.Printf("%+v\n", *resp)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
