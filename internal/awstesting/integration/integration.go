// +build integration

// Package integration performs initialization and validation for integration
// tests.
package integration

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func init() {
	var err error
	config, err = external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err)
	}

	// Disable parameter validation for integration tests.
	config.Handlers.Validate.Remove(defaults.ValidateParametersHandler)

	logLevel := config.LogLevel
	if os.Getenv("DEBUG") != "" {
		logLevel = aws.LogLevel(aws.LogDebug)
	}
	if os.Getenv("DEBUG_SIGNING") != "" {
		logLevel = aws.LogLevel(aws.LogDebugWithSigning)
	}
	if os.Getenv("DEBUG_BODY") != "" {
		logLevel = aws.LogLevel(aws.LogDebugWithSigning | aws.LogDebugWithHTTPBody)
	}
	config.LogLevel = logLevel

	if config.Region == "" {
		panic("AWS_REGION must be configured to run integration tests")
	}
}

var config aws.Config

// Config returns an initialized configuration
func Config() aws.Config {
	return config.Copy()
}

// UniqueID returns a unique UUID-like identifier for use in generating
// resources for integration tests.
func UniqueID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	return fmt.Sprintf("%x", uuid)
}
