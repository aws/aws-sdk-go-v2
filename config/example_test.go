package config_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func ExampleWithSharedConfigProfile() {
	cfg, err := config.LoadDefaultConfig(
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile("exampleProfile"),

		// Optionally specify the specific shared configuraiton
		// files to load the profile from.
		config.WithSharedConfigFiles([]string{
			filepath.Join("testdata", "shared_config"),
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	// Region loaded from credentials file.
	fmt.Println("Region:", cfg.Region)

	// Output:
	// Region: us-west-2
}

func ExampleWithCredentialsProvider() {
	cfg, err := config.LoadDefaultConfig(
		// Hard coded credentials.
		config.WithCredentialsProvider{
			CredentialsProvider: credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
					Source: "example hard coded credentials",
				},
			},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	// Credentials retrieve will be called automatically internally to the SDK
	// service clients created with the cfg value.
	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get credentials, %v", err)
		os.Exit(1)
	}

	fmt.Println("Credentials Source:", creds.Source)

	// Output:
	// Credentials Source: example hard coded credentials
}
