//go:build integration && ec2env
// +build integration,ec2env

package imds

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestInteg_GetDynamicData(t *testing.T) {
	client := New(Options{})

	result, err := client.GetDynamicData(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, err := ioutil.ReadAll(result.Content)
	if err != nil {
		t.Fatalf("expect to read content, got %v", err)
	}

	if len(b) == 0 {
		t.Errorf("expect result content, but was empty")
	}
	t.Logf("Result:\n%s", string(b))
}

func TestInteg_GetIAMInfo(t *testing.T) {
	client := New(Options{})

	result, err := client.GetIAMInfo(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	t.Logf("IAMInfo:\n%#v", result.IAMInfo)
}

func TestInteg_GetInstanceIdentityDocument(t *testing.T) {
	client := New(Options{})

	result, err := client.GetInstanceIdentityDocument(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	t.Logf("InstanceIdentityDocument:\n%#v", result.InstanceIdentityDocument)
}

func TestInteg_GetMetadata(t *testing.T) {
	client := New(Options{})

	result, err := client.GetMetadata(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, err := ioutil.ReadAll(result.Content)
	if err != nil {
		t.Fatalf("expect to read content, got %v", err)
	}

	if len(b) == 0 {
		t.Errorf("expect result content, but was empty")
	}
	t.Logf("Result:\n%s", string(b))
}

func TestInteg_GetRegion(t *testing.T) {
	client := New(Options{})

	result, err := client.GetRegion(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if len(result.Region) == 0 {
		t.Errorf("expect region, got none")
	}

	t.Logf("Region: %s", result.Region)
}

func TestInteg_GetUserData(t *testing.T) {
	if !strings.EqualFold(os.Getenv("AWS_TEST_EC2_IMDS_WITH_USER_DATA"), "true") {
		t.Skip("to run test set AWS_TEST_EC2_IMDS_WITH_USER_DATA=true")
	}

	client := New(Options{})

	result, err := client.GetUserData(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	b, err := ioutil.ReadAll(result.Content)
	if err != nil {
		t.Fatalf("expect to read content, got %v", err)
	}

	if len(b) == 0 {
		t.Errorf("expect result content, but was empty")
	}
	t.Logf("Result:\n%s", string(b))
}
