package config

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestDisableClockSkewCorrectionEnv(t *testing.T) {
	t.Setenv("AWS_DISABLE_CLOCK_SKEW_CORRECTION", "true")

	ec, err := NewEnvConfig()
	if err != nil {
		t.Fatalf("NewEnvConfig: %v", err)
	}

	v, found, err := ec.getDisableClockSkewCorrection(context.Background())
	if err != nil {
		t.Fatalf("getDisableClockSkewCorrection: %v", err)
	}
	if !found || !v {
		t.Fatalf("env: got value=%v found=%v, want true/true", v, found)
	}
}

func TestDisableClockSkewCorrectionResolve(t *testing.T) {
	enabled := true

	cases := map[string]configs{
		"load option":   {LoadOptions{DisableClockSkewCorrection: &enabled}},
		"shared config": {SharedConfig{DisableClockSkewCorrection: &enabled}},
		"env config":    {EnvConfig{DisableClockSkewCorrection: &enabled}},
	}

	for name, cfgs := range cases {
		t.Run(name, func(t *testing.T) {
			var cfg aws.Config
			if err := resolveDisableClockSkewCorrection(context.Background(), &cfg, cfgs); err != nil {
				t.Fatalf("resolve: %v", err)
			}
			if !cfg.DisableClockSkewCorrection {
				t.Fatalf("expected DisableClockSkewCorrection true")
			}
		})
	}
}

func TestDisableClockSkewCorrectionDefault(t *testing.T) {
	var cfg aws.Config
	if err := resolveDisableClockSkewCorrection(context.Background(), &cfg, configs{}); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if cfg.DisableClockSkewCorrection {
		t.Fatalf("expected correction enabled (false) by default")
	}
}
