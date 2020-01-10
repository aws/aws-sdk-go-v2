// +build example

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	roleArn         string
	roleSessionName = "assumeTestRole"
)

func main() {

	// input: aws.Config, roleArn, roleSessionName
	flag.StringVar(&roleArn, "roleArn", "", "role ARN to be assumed")
	flag.Parse()

	config, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("loading default config", err)
	}

	// assume role
	svc := sts.New(config)
	input := &sts.AssumeRoleInput{RoleArn: aws.String(roleArn), RoleSessionName: aws.String(roleSessionName)}
	out, err := svc.AssumeRoleRequest(input).Send(context.TODO())
	if err != nil {
		exitErrorf("aws assume role %s: %v", roleArn, err)
	}

	awsConfig := svc.Config.Copy()
	awsConfig.Credentials = CredentialsProvider{Credentials: out.Credentials}
	fmt.Printf("aws config:\n%+v", awsConfig)
}

type CredentialsProvider struct {
	*sts.Credentials
}

func (s CredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	if s.Credentials == nil {
		return aws.Credentials{}, errors.New("sts credentials are nil")
	}

	return aws.Credentials{
		AccessKeyID:     aws.StringValue(s.AccessKeyId),
		SecretAccessKey: aws.StringValue(s.SecretAccessKey),
		SessionToken:    aws.StringValue(s.SessionToken),
		Expires:         aws.TimeValue(s.Expiration),
	}, nil
}

// --- helper functions ---

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func usage() string {
	return `

Missing mandatory flag(s). Please use like below  Example:

To assume role run:
	./assumeRole -roleArn arn:aws:iam::12345:some/Role
`
}
