// Code generated by smithy-go-codegen DO NOT EDIT.

package opsworkscm

import (
	"context"
	"errors"
	"github.com/aws/smithy-go/middleware"
	"slices"
	"strings"
	"testing"
)

var errTestReturnEarly = errors.New("errTestReturnEarly")

func captureMiddlewareStack(stack *middleware.Stack) func(*middleware.Stack) error {
	return func(inner *middleware.Stack) error {
		*stack = *inner
		return errTestReturnEarly
	}
}
func TestOpAssociateNodeSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.AssociateNode(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpCreateBackupSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.CreateBackup(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpCreateServerSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.CreateServer(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDeleteBackupSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DeleteBackup(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDeleteServerSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DeleteServer(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDescribeAccountAttributesSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DescribeAccountAttributes(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDescribeBackupsSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DescribeBackups(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDescribeEventsSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DescribeEvents(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDescribeNodeAssociationStatusSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DescribeNodeAssociationStatus(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDescribeServersSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DescribeServers(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpDisassociateNodeSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.DisassociateNode(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpExportServerEngineAttributeSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.ExportServerEngineAttribute(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpListTagsForResourceSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.ListTagsForResource(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpRestoreServerSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.RestoreServer(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpStartMaintenanceSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.StartMaintenance(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpTagResourceSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.TagResource(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpUntagResourceSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.UntagResource(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpUpdateServerSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.UpdateServer(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
func TestOpUpdateServerEngineAttributesSRAOperationOrder(t *testing.T) {
	expect := []string{
		"OperationSerializer",
		"Retry",
		"ResolveAuthScheme",
		"GetIdentity",
		"ResolveEndpointV2",
		"Signing",
		"OperationDeserializer",
	}

	var captured middleware.Stack
	svc := New(Options{
		APIOptions: []func(*middleware.Stack) error{
			captureMiddlewareStack(&captured),
		},
	})
	_, err := svc.UpdateServerEngineAttributes(context.Background(), nil)
	if err != nil && !errors.Is(err, errTestReturnEarly) {
		t.Fatalf("unexpected error: %v", err)
	}

	var actual, all []string
	for _, step := range strings.Split(captured.String(), "\n") {
		trimmed := strings.TrimSpace(step)
		all = append(all, trimmed)
		if slices.Contains(expect, trimmed) {
			actual = append(actual, trimmed)
		}
	}

	if !slices.Equal(expect, actual) {
		t.Errorf("order mismatch:\nexpect: %v\nactual: %v\nall: %v", expect, actual, all)
	}
}
