package external

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestConfigs_SharedConfigOptions(t *testing.T) {
	_, err := Configs{
		WithSharedConfigProfile("profile-name"),
		WithSharedConfigFiles([]string{"creds-file"}),
	}.AppendFromLoaders([]ConfigLoader{
		func(configs Configs) (Config, error) {
			var profile string
			var files []string
			var err error

			for _, cfg := range configs {
				if p, ok := cfg.(SharedConfigProfileProvider); ok {
					profile, err = p.GetSharedConfigProfile()
					if err != nil {
						return nil, err
					}
				}
				if p, ok := cfg.(SharedConfigFilesProvider); ok {
					files, err = p.GetSharedConfigFiles()
					if err != nil {
						return nil, err
					}
				}

			}

			if e, a := "profile-name", profile; e != a {
				t.Errorf("expect %v profile, got %v", e, a)
			}
			if e, a := []string{"creds-file"}, files; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %v files, got %v", e, a)
			}

			return nil, nil
		},
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
}

func TestConfigs_AppendFromLoaders(t *testing.T) {
	expectCfg := WithRegion("mock-region")

	cfgs, err := Configs{}.AppendFromLoaders([]ConfigLoader{
		func(configs Configs) (Config, error) {
			if e, a := 0, len(configs); e != a {
				t.Errorf("expect %v configs, got %v", e, a)
			}
			return expectCfg, nil
		},
	})

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := 1, len(cfgs); e != a {
		t.Errorf("expect %v configs, got %v", e, a)
	}

	if e, a := expectCfg, cfgs[0]; e != a {
		t.Errorf("expect %v config, got %v", e, a)
	}
}

func TestConfigs_ResolveAWSConfig(t *testing.T) {
	cfg, err := Configs{
		WithRegion("mock-region"),
		WithCredentialsValue(aws.Credentials{
			AccessKeyID: "AKID", SecretAccessKey: "SECRET",
			Source: "provider",
		}),
	}.ResolveAWSConfig([]AWSConfigResolver{
		ResolveRegion,
		ResolveCredentialsValue,
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expect %v region, got %v", e, a)
	}

	creds, err := cfg.Credentials.Retrieve()
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := "provider", creds.Source; e != a {
		t.Errorf("expect %v provider, got %v", e, a)
	}
}
