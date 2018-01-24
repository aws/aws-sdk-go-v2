package external_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/aws/stscreds"
)

func ExampleWithSharedConfigProfile() {
	cfg, err := external.LoadDefaultAWSConfig(
		// Specify the shared configuration profile to load.
		external.WithSharedConfigProfile("exampleProfile"),

		// Optionally specify the specific shared configuraiton
		// files to load the profile from.
		external.WithSharedConfigFiles([]string{
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

func ExampleWithCredentialsValue() {
	cfg, err := external.LoadDefaultAWSConfig(
		// Hard coded credentials.
		external.WithCredentialsValue(aws.Credentials{
			AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
			Source: "example hard coded credentials",
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	// Credentials retrieve will be called automatically internally to the SDK
	// service clients created with the cfg value.
	creds, err := cfg.Credentials.Retrieve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get credentials, %v", err)
		os.Exit(1)
	}

	fmt.Println("Credentials Source:", creds.Source)

	// Output:
	// Credentials Source: example hard coded credentials
}

func ExampleWithMFATokenFunc() {
	cfg, err := external.LoadDefaultAWSConfig(
		// Set the provider function for the MFA token.
		external.WithMFATokenFunc(stscreds.StdinTokenProvider),

		// Optionally, specify the shared configuration profile to load.
		external.WithSharedConfigProfile("exampleProfile"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	// If assume role credentials with MFA enabled are specified in the shared
	// 	configuration the MFA token provider function will be called to retrieve
	// the MFA token for the assume role API call.
	fmt.Println(cfg.Credentials.Retrieve())
}
