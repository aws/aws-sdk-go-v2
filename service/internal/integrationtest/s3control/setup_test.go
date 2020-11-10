// +build integration

package s3control

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	svc                            *s3control.Client
	s3ControlEndpoint, stsEndpoint string
	accountID                      string
	insecureTLS, useDualstack      bool
)

var region = "us-west-2"

func TestMain(m *testing.M) {
	flag.Parse()
	flag.CommandLine.Visit(func(f *flag.Flag) {
		if !(f.Name == "run" || f.Name == "test.run") {
			return
		}
		value := f.Value.String()
		if value == `NONE` {
			os.Exit(0)
		}
	})

	var result int
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 integration tests panic,", r)
			result = 1
		}
		os.Exit(result)
	}()

	flag.StringVar(&stsEndpoint, "sts-endpoint", "",
		"The optional `URL` endpoint for the STS service.",
	)
	flag.StringVar(&s3ControlEndpoint, "s3-control-endpoint", "",
		"The optional `URL` endpoint for the S3 Control service.",
	)
	flag.BoolVar(&insecureTLS, "insecure-tls", false,
		"Disables TLS validation on request endpoints.",
	)
	flag.BoolVar(&useDualstack, "dualstack", false,
		"Enables usage of dualstack endpoints.",
	)
	flag.StringVar(&accountID, "account", "",
		"The AWS account `ID`.",
	)
	// parse flag
	flag.Parse()

	tlsCfg := &tls.Config{}
	if insecureTLS {
		tlsCfg.InsecureSkipVerify = true
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

	cfg, err := integrationtest.LoadConfigWithDefaultRegion(region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred while loading config with region %v, %v", region, err)
		result = 1
		return
	}
	cfg.HTTPClient = httpClient

	// initialize context
	ctx := context.Background()

	if len(accountID) == 0 {
		var opts = func(options *sts.Options) {}
		if len(stsEndpoint) != 0 {
			opts = func(options *sts.Options) {
				options.EndpointResolver = sts.EndpointResolverFunc(func(region string, options sts.EndpointResolverOptions) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           stsEndpoint,
						PartitionID:   "aws",
						SigningName:   "sts",
						SigningRegion: region,
					}, nil
				})
			}
		}

		// initialize a sts client
		stsClient := sts.NewFromConfig(cfg, opts)

		identity, err := stsClient.GetCallerIdentity(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to get accountID, %v", err))
		}
		accountID = *(identity.Account)
	}

	var s3controlOpts = func(options *s3control.Options) {}
	if len(s3ControlEndpoint) != 0 {
		s3controlOpts = func(options *s3control.Options) {
			options.EndpointResolver = s3control.EndpointResolverFunc(func(region string, options s3control.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           s3ControlEndpoint,
					PartitionID:   "aws",
					SigningName:   "s3-control",
					SigningRegion: region,
				}, nil
			})
		}
	}
	// construct a s3-control client
	svc = s3control.NewFromConfig(cfg, s3controlOpts)

	result = m.Run()
}
