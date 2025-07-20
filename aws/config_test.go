package aws

import (
	"testing"
)

// Mock service options struct for testing
type MockServiceOptions struct {
	Field1 bool
	Field2 string
	Field3 int
}

func TestWithServiceOptions(t *testing.T) {
	cfg := NewConfig()

	cfg = cfg.WithServiceOptions("TestService", func(opts any) {
		if mockOpts, ok := opts.(*MockServiceOptions); ok {
			mockOpts.Field1 = true
			mockOpts.Field2 = "test"
		}
	})

	if cfg.ServiceOptions == nil {
		t.Fatal("ServiceOptions should not be nil")
	}

	callbacks, exists := cfg.ServiceOptions["TestService"]
	if !exists {
		t.Fatal("TestService should exist in ServiceOptions")
	}

	if len(callbacks) != 1 {
		t.Fatalf("Expected 1 callback, got %d", len(callbacks))
	}

	mockOpts := &MockServiceOptions{}
	cfg.ApplyServiceOptions("TestService", mockOpts)

	if !mockOpts.Field1 {
		t.Error("Field1 should be true")
	}

	if mockOpts.Field2 != "test" {
		t.Errorf("Field2 should be 'test', got '%s'", mockOpts.Field2)
	}
}

func TestWithServiceOptionsMultiple(t *testing.T) {
	cfg := NewConfig()

	cfg = cfg.WithServiceOptions("TestService",
		func(opts any) {
			if mockOpts, ok := opts.(*MockServiceOptions); ok {
				mockOpts.Field1 = true
			}
		},
		func(opts any) {
			if mockOpts, ok := opts.(*MockServiceOptions); ok {
				mockOpts.Field2 = "test"
			}
		},
	)

	callbacks, exists := cfg.ServiceOptions["TestService"]
	if !exists {
		t.Fatal("TestService should exist in ServiceOptions")
	}

	if len(callbacks) != 2 {
		t.Fatalf("Expected 2 callbacks, got %d", len(callbacks))
	}

	mockOpts := &MockServiceOptions{}
	cfg.ApplyServiceOptions("TestService", mockOpts)

	if !mockOpts.Field1 {
		t.Error("Field1 should be true")
	}

	if mockOpts.Field2 != "test" {
		t.Errorf("Field2 should be 'test', got '%s'", mockOpts.Field2)
	}
}

func TestApplyServiceOptionsNonExistent(t *testing.T) {
	cfg := NewConfig()

	mockOpts := &MockServiceOptions{}

	cfg.ApplyServiceOptions("NonExistentService", mockOpts)

	if mockOpts.Field1 || mockOpts.Field2 != "" || mockOpts.Field3 != 0 {
		t.Error("Options should not be modified for non-existent service")
	}
}

func TestTypeAssertionFailure(t *testing.T) {
	cfg := NewConfig()

	cfg = cfg.WithServiceOptions("TestService", func(opts any) {
		if mockOpts, ok := opts.(*MockServiceOptions); ok {
			mockOpts.Field1 = true
		}
	})

	differentOpts := &struct{ Field string }{Field: "test"}
	cfg.ApplyServiceOptions("TestService", differentOpts)

	if differentOpts.Field != "test" {
		t.Error("Different options should not be modified")
	}
}

func TestChaining(t *testing.T) {
	cfg := NewConfig()

	cfg = cfg.WithRegion("us-west-2").
		WithServiceOptions("Service1", func(opts any) {
			if mockOpts, ok := opts.(*MockServiceOptions); ok {
				mockOpts.Field1 = true
			}
		}).
		WithServiceOptions("Service2", func(opts any) {
			if mockOpts, ok := opts.(*MockServiceOptions); ok {
				mockOpts.Field2 = "chained"
			}
		})

	if cfg.Region != "us-west-2" {
		t.Errorf("Expected region 'us-west-2', got '%s'", cfg.Region)
	}

	if len(cfg.ServiceOptions) != 2 {
		t.Fatalf("Expected 2 services, got %d", len(cfg.ServiceOptions))
	}

	mockOpts1 := &MockServiceOptions{}
	cfg.ApplyServiceOptions("Service1", mockOpts1)

	if !mockOpts1.Field1 {
		t.Error("Service1 Field1 should be true")
	}

	mockOpts2 := &MockServiceOptions{}
	cfg.ApplyServiceOptions("Service2", mockOpts2)

	if mockOpts2.Field2 != "chained" {
		t.Errorf("Service2 Field2 should be 'chained', got '%s'", mockOpts2.Field2)
	}
}
