package external

import "testing"

func TestDefaultConfig(t *testing.T) {
	t.Errorf("not implemented")
}

func TestConfig_WithLoaders(t *testing.T) {
	cfgs, err := Configs{
		StaticSharedConfigProfile("default"),
		StaticSharedConfigFiles([]string{"testdata/creds"}),
	}.AppendFromLoaders(
		DefaultConfigLoaders...,
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	// TODO validate cfg
	_, err = cfgs.ResolveAWSConfig(
		DefaultAWSConfigResolvers...,
	)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

}
