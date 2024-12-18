//go:build go1.21
// +build go1.21

package checksum

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"testing"

	internalcontext "github.com/aws/aws-sdk-go-v2/internal/context"
	"github.com/aws/smithy-go/middleware"
)

func TestSetupInput(t *testing.T) {
	type Params struct {
		Value string
	}

	cases := map[string]struct {
		inputParams                interface{}
		getAlgorithm               func(interface{}) (string, bool)
		RequireChecksum            bool
		RequestChecksumCalculation aws.RequestChecksumCalculation
		expectValue                string
	}{
		"user config require checksum and algorithm unset": {
			RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			getAlgorithm: func(v interface{}) (string, bool) {
				return "", false
			},
			expectValue: "",
		},
		"require checksum found empty": {
			RequireChecksum: true,
			inputParams:     Params{Value: ""},
			getAlgorithm: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "",
		},
		"user config require checksum found empty": {
			RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			inputParams:                Params{Value: ""},
			getAlgorithm: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "",
		},
		"require checksum and found": {
			RequireChecksum: true,
			inputParams:     Params{Value: "abc123"},
			getAlgorithm: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "abc123",
		},
		"user config support checksum and found": {
			inputParams: Params{Value: "abc123"},
			getAlgorithm: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "abc123",
		},
		"user config require checksum and found": {
			RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			inputParams:                Params{Value: "abc123"},
			getAlgorithm: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "abc123",
		},
		"require checksum unset and use default": {
			RequireChecksum: true,
			getAlgorithm: func(v interface{}) (string, bool) {
				return "", false
			},
			expectValue: "CRC32",
		},
		"user config support checksum and use default": {
			RequestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
			getAlgorithm: func(v interface{}) (string, bool) {
				return "", false
			},
			expectValue: "CRC32",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			m := setupInputContext{
				GetAlgorithm:               c.getAlgorithm,
				RequireChecksum:            c.RequireChecksum,
				RequestChecksumCalculation: c.RequestChecksumCalculation,
			}

			_, _, err := m.HandleInitialize(context.Background(),
				middleware.InitializeInput{Parameters: c.inputParams},
				middleware.InitializeHandlerFunc(
					func(ctx context.Context, input middleware.InitializeInput) (
						out middleware.InitializeOutput, metadata middleware.Metadata, err error,
					) {
						v := internalcontext.GetChecksumInputAlgorithm(ctx)
						if e, a := c.expectValue, v; e != a {
							t.Errorf("expect value %v, got %v", e, a)
						}

						return out, metadata, nil
					},
				))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

		})
	}
}

func TestSetupOutput(t *testing.T) {
	type Params struct {
		Value string
	}

	cases := map[string]struct {
		inputParams                interface{}
		ResponseChecksumValidation aws.ResponseChecksumValidation
		getValidationMode          func(interface{}) (string, bool)
		expectValue                string
	}{
		"user config support checksum found empty": {
			ResponseChecksumValidation: aws.ResponseChecksumValidationWhenSupported,
			inputParams:                Params{Value: ""},
			getValidationMode: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "ENABLED",
		},
		"user config support checksum found invalid value": {
			ResponseChecksumValidation: aws.ResponseChecksumValidationWhenSupported,
			inputParams:                Params{Value: "abc123"},
			getValidationMode: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "ENABLED",
		},
		"user config require checksum found invalid value": {
			ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
			inputParams:                Params{Value: "abc123"},
			getValidationMode: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "",
		},
		"user config require checksum found valid value": {
			ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
			inputParams:                Params{Value: "ENABLED"},
			getValidationMode: func(v interface{}) (string, bool) {
				vv := v.(Params)
				return vv.Value, true
			},
			expectValue: "ENABLED",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			m := setupOutputContext{
				GetValidationMode:          c.getValidationMode,
				ResponseChecksumValidation: c.ResponseChecksumValidation,
			}

			_, _, err := m.HandleInitialize(context.Background(),
				middleware.InitializeInput{Parameters: c.inputParams},
				middleware.InitializeHandlerFunc(
					func(ctx context.Context, input middleware.InitializeInput) (
						out middleware.InitializeOutput, metadata middleware.Metadata, err error,
					) {
						v := getContextOutputValidationMode(ctx)
						if e, a := c.expectValue, v; e != a {
							t.Errorf("expect value %v, got %v", e, a)
						}

						return out, metadata, nil
					},
				))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

		})
	}
}
