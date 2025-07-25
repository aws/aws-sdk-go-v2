package kitchensinktest

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type allInterceptors struct {
	called []string
}

func (i *allInterceptors) BeforeExecution(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	if in.Input == nil {
		return errors.New("input should be available")
	}

	i.called = append(i.called, "BeforeExecution")
	return nil
}

func (i *allInterceptors) BeforeSerialization(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeSerialization")
	return nil
}

func (i *allInterceptors) AfterSerialization(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	if in.Request == nil {
		return errors.New("request should be available")
	}

	i.called = append(i.called, "AfterSerialization")
	return nil
}

func (i *allInterceptors) BeforeRetryLoop(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeRetryLoop")
	return nil
}

func (i *allInterceptors) BeforeAttempt(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeAttempt")
	return nil
}

func (i *allInterceptors) BeforeSigning(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeSigning")
	return nil
}

func (i *allInterceptors) AfterSigning(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "AfterSigning")
	return nil
}

func (i *allInterceptors) BeforeTransmit(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeTransmit")
	return nil
}

func (i *allInterceptors) AfterTransmit(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	if in.Response == nil {
		return errors.New("response should be available")
	}

	i.called = append(i.called, "AfterTransmit")
	return nil
}

func (i *allInterceptors) BeforeDeserialization(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "BeforeDeserialization")
	return nil
}

func (i *allInterceptors) AfterDeserialization(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	// presence of output is conditional
	i.called = append(i.called, "AfterDeserialization")
	return nil
}

func (i *allInterceptors) AfterAttempt(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "AfterAttempt")
	return nil
}

func (i *allInterceptors) AfterExecution(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	i.called = append(i.called, "AfterExecution")
	return nil
}

func (i *allInterceptors) Register(r *smithyhttp.InterceptorRegistry) {
	r.AddBeforeExecution(i)
	r.AddBeforeSerialization(i)
	r.AddAfterSerialization(i)
	r.AddBeforeRetryLoop(i)
	r.AddBeforeAttempt(i)
	r.AddBeforeSigning(i)
	r.AddAfterSigning(i)
	r.AddBeforeTransmit(i)
	r.AddAfterTransmit(i)
	r.AddBeforeDeserialization(i)
	r.AddAfterDeserialization(i)
	r.AddAfterAttempt(i)
	r.AddAfterExecution(i)
}

type expectOutputInterceptor struct{}

func (i *expectOutputInterceptor) AfterDeserialization(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	if in.Output == nil {
		return errors.New("output should be available")
	}
	return nil
}

func (i *expectOutputInterceptor) AfterExecution(ctx context.Context, in *smithyhttp.InterceptorContext) error {
	if in.Output == nil {
		return errors.New("output should be available")
	}
	return nil
}

func TestInterceptor_PerOperationConfig(t *testing.T) {
	i := &allInterceptors{}
	svc := New(Options{}, func(o *Options) {
		o.Interceptors.AddBeforeExecution(i)
	})

	svc.GetItem(context.Background(), nil)
	svc.GetItem(context.Background(), nil, func(o *Options) {
		o.Interceptors.BeforeExecution = nil
		o.Interceptors.AddAfterExecution(i)
	})
	svc.GetItem(context.Background(), nil, func(o *Options) {
		o.Interceptors.AddAfterExecution(i)
	})

	expect := []string{
		"BeforeExecution", // 1
		"AfterExecution",  // 2
		"BeforeExecution", // 3
		"AfterExecution",
	}
	if !reflect.DeepEqual(expect, i.called) {
		t.Errorf("expect interceptor calls: %#v != %#v", expect, i.called)
	}
}

func TestInterceptor_CorrectOrder(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody("{}"),
	}

	i := &allInterceptors{}
	svc := New(Options{
		HTTPClient: &mockHTTP{
			resps: []*http.Response{mockResp},
		},
	}, func(o *Options) {
		i.Register(&o.Interceptors)
	})

	_, err := svc.GetItem(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"BeforeExecution",
		"BeforeSerialization",
		"AfterSerialization",
		"BeforeRetryLoop",
		"BeforeAttempt",
		"BeforeSigning",
		"AfterSigning",
		"BeforeTransmit",
		"AfterTransmit",
		"BeforeDeserialization",
		"AfterDeserialization",
		"AfterAttempt",
		"AfterExecution",
	}
	if !reflect.DeepEqual(expect, i.called) {
		t.Errorf("expect interceptor calls: %#v != %#v", expect, i.called)
	}
}

func TestInterceptor_TransmitFailure(t *testing.T) {
	i := &allInterceptors{}
	svc := New(Options{
		HTTPClient: &mockHTTP{
			err: errors.New("the ethernet cable exploded"),
		},
	}, func(o *Options) {
		o.RetryMaxAttempts = 2
		i.Register(&o.Interceptors)
	})

	_, err := svc.GetItem(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got none")
	}

	expect := []string{
		// init
		"BeforeExecution",
		"BeforeSerialization",
		"AfterSerialization",
		"BeforeRetryLoop",
		// attempt 1
		"BeforeAttempt",
		"BeforeSigning",
		"AfterSigning",
		"BeforeTransmit", // AfterTransmit & Deserialization shouldn't fire
		"AfterAttempt",
		// attempt 2
		"BeforeAttempt",
		"BeforeSigning",
		"AfterSigning",
		"BeforeTransmit",
		"AfterAttempt",
		// end
		"AfterExecution",
	}
	if !reflect.DeepEqual(expect, i.called) {
		t.Errorf("expect interceptor calls: %#v != %#v", expect, i.called)
	}
}

func TestInterceptor_Transmit4XX(t *testing.T) {
	mockResp1 := &http.Response{
		StatusCode: 500,
		Body:       mockResponseBody("{}"),
	}
	mockResp2 := &http.Response{
		StatusCode: 500,
		Body:       mockResponseBody("{}"),
	}

	i := &allInterceptors{}
	svc := New(Options{
		HTTPClient: &mockHTTP{
			resps: []*http.Response{mockResp1, mockResp2},
		},
	}, func(o *Options) {
		o.RetryMaxAttempts = 2
		i.Register(&o.Interceptors)
	})

	_, err := svc.GetItem(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got none")
	}

	expect := []string{
		// init
		"BeforeExecution",
		"BeforeSerialization",
		"AfterSerialization",
		"BeforeRetryLoop",
		// attempt 1
		"BeforeAttempt",
		"BeforeSigning",
		"AfterSigning",
		"BeforeTransmit",
		"AfterTransmit",
		"BeforeDeserialization",
		"AfterDeserialization",
		"AfterAttempt",
		// attempt 2
		"BeforeAttempt",
		"BeforeSigning",
		"AfterSigning",
		"BeforeTransmit",
		"AfterTransmit",
		"BeforeDeserialization",
		"AfterDeserialization",
		"AfterAttempt",
		// end
		"AfterExecution",
	}
	if !reflect.DeepEqual(expect, i.called) {
		t.Errorf("expect interceptor calls:\n%#v !=\n%#v", expect, i.called)
	}
}

func TestInterceptor_OutputAvailable(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 200,
		Body:       mockResponseBody("{}"),
	}

	// will assert that outputs are available where we expect them
	i := &expectOutputInterceptor{}
	svc := New(Options{
		HTTPClient: &mockHTTP{
			resps: []*http.Response{mockResp},
		},
	}, func(o *Options) {
		o.Interceptors.AddAfterDeserialization(i)
		o.Interceptors.AddAfterExecution(i)
	})

	_, err := svc.GetItem(context.Background(), nil)
	if err != nil {
		t.Fatalf("expect no err, got %v", err)
	}
}
