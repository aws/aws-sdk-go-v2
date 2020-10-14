package config_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	smithymiddleware "github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
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

func ExampleWithAPIOptions() {
	cfg, err := config.LoadDefaultConfig(config.WithAPIOptions([]func(stack *smithymiddleware.Stack) error{
		smithyhttp.AddHeaderValue("X-Custom-Header", "customHeaderValue"),
	}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	fmt.Println("Number of APIOptions:", len(cfg.APIOptions))
	// Output:
	// Number of APIOptions: 1
}

func ExampleWithEndpointResolver() {
	cfg, err := config.LoadDefaultConfig(config.WithEndpointResolver{
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "https://mock.amazonaws.com"}, nil
		}),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}

	resolver, _ := cfg.EndpointResolver.ResolveEndpoint("service", "region")

	fmt.Println(resolver.URL)
	// Output:
	// https://mock.amazonaws.com
}

func ExampleWithHTTPClient() {
	cfg, err := config.LoadDefaultConfig(config.WithHTTPClient{
		HTTPClient: aws.NewBuildableHTTPClient().WithTransportOptions(func(tr *http.Transport) {
			tr.MaxIdleConns = 60
		}),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}
	_ = cfg
}

func ExampleWithAssumeRoleCredentialProviderOptions() {
	cfg, err := config.LoadDefaultConfig(config.WithAssumeRoleCredentialProviderOptions(func(options *stscreds.AssumeRoleOptions) {
		options.RoleSessionName = "customRoleSessionName"
		options.TokenProvider = func() (string, error) {
			return "mfaTokenValue", nil
		}
	}))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}
	_ = cfg
}

func ExampleWithRegion() {
	cfg, err := config.LoadDefaultConfig(config.WithRegion("us-west-2"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}
	fmt.Println(cfg.Region)
	// Output: us-west-2
}

func ExampleWithEC2IMDSRegion() {
	cfg, err := config.LoadDefaultConfig(config.WithEC2IMDSRegion{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, %v", err)
		os.Exit(1)
	}
	_ = cfg
}
