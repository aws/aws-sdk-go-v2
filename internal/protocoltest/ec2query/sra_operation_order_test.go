// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2query

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
func TestOpDatetimeOffsetsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.DatetimeOffsets(context.Background(), nil)
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
func TestOpEmptyInputAndEmptyOutputSRAOperationOrder(t *testing.T) {
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
	_, err := svc.EmptyInputAndEmptyOutput(context.Background(), nil)
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
func TestOpEndpointOperationSRAOperationOrder(t *testing.T) {
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
	_, err := svc.EndpointOperation(context.Background(), nil)
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
func TestOpEndpointWithHostLabelOperationSRAOperationOrder(t *testing.T) {
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
	_, err := svc.EndpointWithHostLabelOperation(context.Background(), nil)
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
func TestOpFractionalSecondsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.FractionalSeconds(context.Background(), nil)
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
func TestOpGreetingWithErrorsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.GreetingWithErrors(context.Background(), nil)
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
func TestOpHostWithPathOperationSRAOperationOrder(t *testing.T) {
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
	_, err := svc.HostWithPathOperation(context.Background(), nil)
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
func TestOpIgnoresWrappingXmlNameSRAOperationOrder(t *testing.T) {
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
	_, err := svc.IgnoresWrappingXmlName(context.Background(), nil)
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
func TestOpNestedStructuresSRAOperationOrder(t *testing.T) {
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
	_, err := svc.NestedStructures(context.Background(), nil)
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
func TestOpNoInputAndOutputSRAOperationOrder(t *testing.T) {
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
	_, err := svc.NoInputAndOutput(context.Background(), nil)
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
func TestOpPutWithContentEncodingSRAOperationOrder(t *testing.T) {
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
	_, err := svc.PutWithContentEncoding(context.Background(), nil)
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
func TestOpQueryIdempotencyTokenAutoFillSRAOperationOrder(t *testing.T) {
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
	_, err := svc.QueryIdempotencyTokenAutoFill(context.Background(), nil)
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
func TestOpQueryListsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.QueryLists(context.Background(), nil)
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
func TestOpQueryTimestampsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.QueryTimestamps(context.Background(), nil)
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
func TestOpRecursiveXmlShapesSRAOperationOrder(t *testing.T) {
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
	_, err := svc.RecursiveXmlShapes(context.Background(), nil)
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
func TestOpSimpleInputParamsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.SimpleInputParams(context.Background(), nil)
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
func TestOpSimpleScalarXmlPropertiesSRAOperationOrder(t *testing.T) {
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
	_, err := svc.SimpleScalarXmlProperties(context.Background(), nil)
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
func TestOpXmlBlobsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlBlobs(context.Background(), nil)
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
func TestOpXmlEmptyBlobsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlEmptyBlobs(context.Background(), nil)
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
func TestOpXmlEmptyListsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlEmptyLists(context.Background(), nil)
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
func TestOpXmlEnumsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlEnums(context.Background(), nil)
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
func TestOpXmlIntEnumsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlIntEnums(context.Background(), nil)
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
func TestOpXmlListsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlLists(context.Background(), nil)
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
func TestOpXmlNamespacesSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlNamespaces(context.Background(), nil)
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
func TestOpXmlTimestampsSRAOperationOrder(t *testing.T) {
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
	_, err := svc.XmlTimestamps(context.Background(), nil)
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
