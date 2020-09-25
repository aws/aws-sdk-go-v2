// Package sdk is the official AWS SDK v2 for the Go programming language.
//
// aws-sdk-go-v2 is the Developer Preview for the v2 of the AWS SDK for the Go
// programming language. Look for additional documentation and examples to be
// added.
//
// Getting started
//
// The best way to get started working with the SDK is to use `go get` to add the
// SDK to your Go Workspace manually.
//
// 	go get github.com/aws/aws-sdk-go-v2
//
// You could also use [Dep] to add the SDK to your application's dependencies.
// Using [Dep] will simplify your update story and help your application keep
// pinned to specific version of the SDK
//
// 	dep ensure --add github.com/aws/aws-sdk-go-v2
//
// Hello AWS
//
// This example shows how you can use the v2 SDK to make an API request using the
// SDK's Amazon DynamoDB client.
//
// 	package main
//
// 	import (
// 		"context"
// 		"fmt"
// 		"log"
//
// 		"github.com/aws/aws-sdk-go-v2/aws"
// 		"github.com/aws/aws-sdk-go-v2/config"
// 		"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	)
//
// 	func main() {
// 		// Using the SDK's default configuration, loading additional config
// 		// and credentials values from the environment variables, shared
// 		// credentials, and shared configuration files
// 		cfg, err := config.LoadDefaultConfig()
// 		if err != nil {
// 			log.Fatalf("unable to load SDK config, %v", err)
// 		}
//
// 		// Set the AWS Region that the service clients should use
// 		cfg.Region = "us-west-2"
//
// 		// Using the Config value, create the DynamoDB client
// 		svc := dynamodb.NewFromConfig(cfg)
//
// 		// Send the request, and get the response or error back
// 		resp, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
// 			TableName: aws.String("myTable"),
// 		})
// 		if err != nil {
// 			log.Fatalf("failed to describe table, %v", err)
// 		}
//
// 		fmt.Println("Response", resp)
// 	}
package sdk
