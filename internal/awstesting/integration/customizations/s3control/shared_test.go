// +build integration

package s3control

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	svc                            *s3control.S3Control
	s3ControlEndpoint, stsEndpoint string
	accountID                      string
	insecureTLS, useDualstack      bool
)

func init() {
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
	flag.Parse()
}

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())
}

// Create a bucket for testing
func setup() {
	tlsCfg := &tls.Config{}
	if insecureTLS {
		tlsCfg.InsecureSkipVerify = true
	}

	cfg := integration.Config()
	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}
	if useDualstack {
		resolver := endpoints.NewDefaultResolver()
		resolver.UseDualStack = useDualstack
		cfg.EndpointResolver = resolver
	}

	if len(accountID) == 0 {
		stsCfg := cfg.Copy()
		if len(stsEndpoint) != 0 {
			stsCfg.EndpointResolver = aws.ResolveWithEndpointURL(stsEndpoint)
		}
		stsSvc := sts.New(stsCfg)
		identity, err := stsSvc.GetCallerIdentityRequest(
			&sts.GetCallerIdentityInput{},
		).Send()
		if err != nil {
			panic(fmt.Sprintf("failed to get accountID, %v", err))
		}
		accountID = aws.StringValue(identity.Account)
	}

	s3CtrlCfg := cfg.Copy()
	if len(s3ControlEndpoint) != 0 {
		s3CtrlCfg.EndpointResolver = aws.ResolveWithEndpointURL(s3ControlEndpoint)
	}

	svc = s3control.New(s3CtrlCfg)
}
