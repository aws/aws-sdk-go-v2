//go:build integration && ec2env
// +build integration,ec2env

package ec2rolecreds

import (
	"context"
	"testing"
)

func TestInteg_RetrieveCredentials(t *testing.T) {
	provider := New()

	creds, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if !creds.HasKeys() {
		t.Errorf("expect credential values, got none")
	}

	t.Logf("AccessKey: %v", creds.AccessKeyID)
}
