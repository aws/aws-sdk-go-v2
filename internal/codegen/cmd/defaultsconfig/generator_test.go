package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateConfigPackage_RealModel(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "sdk-default-configuration.json")
	outputFile := filepath.Join(t.TempDir(), "defaults.go")

	err := generateConfigPackage(modelPath, outputFile, "defaults", "GetModeConfiguration")
	if err != nil {
		t.Fatalf("failed to parse model: %v", err)
	}

	output, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("generated output is empty")
	}
}
