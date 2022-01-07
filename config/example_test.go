package config_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func ExampleWithCredentialsCacheOptions() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsCacheOptions(func(o *aws.CredentialsCacheOptions) {
			o.ExpiryWindow = 10 * time.Minute
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithSharedConfigProfile() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// Specify the shared configuration profile to load.
		config.WithSharedConfigProfile("exampleProfile"),

		// Optionally specify the specific shared configuraiton
		// files to load the profile from.
		config.WithSharedConfigFiles([]string{
			filepath.Join("testdata", "shared_config"),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Region loaded from credentials file.
	fmt.Println("Region:", cfg.Region)

	// Output:
	// Region: us-west-2
}

func ExampleWithCredentialsProvider() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// Hard coded credentials.
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
				Source: "example hard coded credentials",
			},
		}))
	if err != nil {
		log.Fatal(err)
	}

	// Credentials retrieve will be called automatically internally to the SDK
	// service clients created with the cfg value.
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Credentials Source:", creds.Source)
	// Credentials Source: example hard coded credentials
}

func ExampleWithAPIOptions() {
	// import "github.com/aws/smithy-go/middleware"
	// import smithyhttp "github.com/aws/smithy-go/transport/http"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithAPIOptions([]func(*middleware.Stack) error{
			smithyhttp.AddHeaderValue("X-Custom-Header", "customHeaderValue"),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithEndpointResolver() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "https://mock.amazonaws.com"}, nil
			})),
	)

	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithEndpointResolverWithOptions() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "https://mock.amazonaws.com"}, nil
			})),
	)

	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithHTTPClient() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithHTTPClient(awshttp.NewBuildableClient().
			WithTransportOptions(func(tr *http.Transport) {
				tr.MaxIdleConns = 60
			})),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithWebIdentityRoleCredentialOptions() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithWebIdentityRoleCredentialOptions(func(options *stscreds.WebIdentityRoleOptions) {
			options.RoleSessionName = "customSessionName"
		}))
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithRegion() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithEC2IMDSRegion() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEC2IMDSRegion(),
	)
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}

func ExampleWithAssumeRoleCredentialOptions() {
	// WithAssumeRoleCredentialOptions can be used to configure the AssumeRoleOptions for the STS credential provider.
	// For example the TokenProvider can be populated if assuming a role that requires an MFA token.
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = func() (string, error) {
				return "theTokenCode", nil
			}
		}))
	if err != nil {
		log.Fatal(err)
	}
	_ = cfg
}
