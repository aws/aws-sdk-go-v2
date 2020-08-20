package external

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func resolveHTTPCredProvider(cfg *aws.Config, url, authToken string, configs Configs) error {
	return fmt.Errorf("endpoint credential provider is not currently supported")
}

func resolveEC2RoleCredentials(cfg *aws.Config, configs Configs) error {
	return fmt.Errorf("ec2 role credential provider is not currently supported")
}

func assumeWebIdentity(cfg *aws.Config, filepath string, roleARN, sessionName string, configs Configs) error {
	return fmt.Errorf("assume web identity role is not currently supported")
}

func credsFromAssumeRole(cfg *aws.Config, sharedCfg *SharedConfig, configs Configs) (err error) {
	return fmt.Errorf("assume role credentials is not currently supported")
}
