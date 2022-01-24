package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/google/go-cmp/cmp"
)

func TestConfigs_SharedConfigOptions(t *testing.T) {
	var options LoadOptions
	optFns := []func(*LoadOptions) error{
		WithSharedConfigProfile("profile-name"),
		WithSharedConfigFiles([]string{"creds-file"}),
	}

	for _, optFn := range optFns {
		optFn(&options)
	}

	_, err := configs{options}.AppendFromLoaders(context.TODO(), []loader{
		func(ctx context.Context, configs configs) (Config, error) {
			var profile string
			var found bool
			var files []string
			var err error

			for _, cfg := range configs {
				if p, ok := cfg.(sharedConfigProfileProvider); ok {
					profile, found, err = p.getSharedConfigProfile(ctx)
					if err != nil || !found {
						return nil, err
					}
				}
				if p, ok := cfg.(sharedConfigFilesProvider); ok {
					files, found, err = p.getSharedConfigFiles(ctx)
					if err != nil || !found {
						return nil, err
					}
				}
			}

			if e, a := "profile-name", profile; e != a {
				t.Errorf("expect %v profile, got %v", e, a)
			}
			if diff := cmp.Diff([]string{"creds-file"}, files); len(diff) != 0 {
				t.Errorf("expect resolved shared config match, got diff: \n %s", diff)
			}

			return nil, nil
		},
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}

func TestConfigs_AppendFromLoaders(t *testing.T) {
	var options LoadOptions
	err := WithRegion("mock-region")(&options)
	if err != nil {
		t.Fatalf("expect not error, got %v", err)
	}

	cfgs, err := configs{}.AppendFromLoaders(
		context.TODO(), []loader{
			func(ctx context.Context, configs configs) (Config, error) {
				if e, a := 0, len(configs); e != a {
					t.Errorf("expect %v configs, got %v", e, a)
				}
				return options, nil
			},
		})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 1, len(cfgs); e != a {
		t.Errorf("expect %v configs, got %v", e, a)
	}

	if diff := cmp.Diff(options, cfgs[0]); len(diff) != 0 {
		t.Errorf("expect config match, got diff: \n %s", diff)
	}
}

func TestConfigs_ResolveAWSConfig(t *testing.T) {
	var options LoadOptions
	optFns := []func(*LoadOptions) error{
		WithRegion("mock-region"),
		WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "AKID", SecretAccessKey: "SECRET",
				Source: "provider",
			},
		}),
	}

	for _, optFn := range optFns {
		optFn(&options)
	}

	config := configs{options}

	cfg, err := config.ResolveAWSConfig(context.TODO(), []awsConfigResolver{
		resolveRegion,
		resolveCredentials,
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expect %v region, got %v", e, a)
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "provider", creds.Source; e != a {
		t.Errorf("expect %v provider, got %v", e, a)
	}

	var expectedSources []interface{}
	for _, s := range cfg.ConfigSources {
		expectedSources = append(expectedSources, s)
	}

	if diff := cmp.Diff(expectedSources, cfg.ConfigSources); len(diff) != 0 {
		t.Errorf("expect config sources match, got diff: \n %s", diff)
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	optWithErr := func(_ *LoadOptions) error {
		return fmt.Errorf("some error")
	}
	_, err := LoadDefaultConfig(context.TODO(), optWithErr)
	if err == nil {
		t.Fatal("expect error when optFn returns error, got nil")
	}
}
