package aws_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockClient struct {
	*aws.Client
}
type MockInput struct{}
type MockOutput struct {
	States []MockState
}
type MockState struct {
	StatePtr *string
	State    StateType
}

type StateType string

const (
	StateTypeStopping StateType = "stopping"
	StateTypePending  StateType = "pending"
	StateTypeRunning  StateType = "running"
)

func (c *mockClient) MockRequest(input *MockInput) (*aws.Request, *MockOutput) {
	op := &aws.Operation{
		Name:       "Mock",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &MockInput{}
	}

	output := &MockOutput{}
	req := c.NewRequest(op, input, output)
	req.Data = output
	return req, output
}

func BuildNewMockRequest(c *mockClient, in *MockInput) func([]aws.Option) (*aws.Request, error) {
	return func(opts []aws.Option) (*aws.Request, error) {
		req, _ := c.MockRequest(in)
		req.ApplyOptions(opts...)
		return req, nil
	}
}

func TestWaiterPathAll(t *testing.T) {
	svc := &mockClient{Client: awstesting.NewClient(unit.Config())}

	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	reqNum := 0
	resps := []*MockOutput{
		{ // Request 1
			States: []MockState{
				{State: StateTypePending},
				{State: StateTypePending},
			},
		},
		{ // Request 2
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypePending},
			},
		},
		{ // Request 3
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypeRunning},
			},
		},
	}

	numBuiltReq := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		numBuiltReq++
	})
	svc.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		if reqNum >= len(resps) {
			t.Fatal("too many polling requests made")
			return
		}
		r.Data = resps[reqNum]
		reqNum++
	})

	w := aws.Waiter{
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(0),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.PathAllWaiterMatch,
				Argument: "States[].State",
				Expected: "running",
			},
		},
		NewRequest: BuildNewMockRequest(svc, &MockInput{}),
	}

	err := w.WaitWithContext(aws.BackgroundContext())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 3, numBuiltReq; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 3, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestWaiterPath(t *testing.T) {
	svc := &mockClient{Client: awstesting.NewClient(unit.Config())}

	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	reqNum := 0
	resps := []*MockOutput{
		{ // Request 1
			States: []MockState{
				{State: StateTypePending},
				{State: StateTypePending},
			},
		},
		{ // Request 2
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypePending},
			},
		},
		{ // Request 3
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypeRunning},
			},
		},
	}

	numBuiltReq := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		numBuiltReq++
	})
	svc.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		if reqNum >= len(resps) {
			t.Fatalf("too many polling requests made")
			return
		}
		r.Data = resps[reqNum]
		reqNum++
	})

	w := aws.Waiter{
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(0),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.PathWaiterMatch,
				Argument: "States[].State",
				Expected: "running",
			},
		},
		NewRequest: BuildNewMockRequest(svc, &MockInput{}),
	}

	err := w.WaitWithContext(aws.BackgroundContext())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 3, numBuiltReq; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 3, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestWaiterFailure(t *testing.T) {
	svc := &mockClient{Client: awstesting.NewClient(unit.Config())}

	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	reqNum := 0
	resps := []*MockOutput{
		{ // Request 1
			States: []MockState{
				{State: StateTypePending},
				{State: StateTypePending},
			},
		},
		{ // Request 2
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypePending},
			},
		},
		{ // Request 3
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypeStopping},
			},
		},
	}

	numBuiltReq := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		numBuiltReq++
	})
	svc.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		if reqNum >= len(resps) {
			t.Fatalf("too many polling requests made")
			return
		}
		r.Data = resps[reqNum]
		reqNum++
	})

	w := aws.Waiter{
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(0),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.PathAllWaiterMatch,
				Argument: "States[].State",
				Expected: "running",
			},
			{
				State:    aws.FailureWaiterState,
				Matcher:  aws.PathAnyWaiterMatch,
				Argument: "States[].State",
				Expected: "stopping",
			},
		},
		NewRequest: BuildNewMockRequest(svc, &MockInput{}),
	}

	err := w.WaitWithContext(aws.BackgroundContext()).(awserr.Error)
	if err == nil {
		t.Fatalf("expect error, got none")
	}
	if e, a := aws.WaiterResourceNotReadyErrorCode, err.Code(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "failed waiting for successful resource state", err.Message(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 3, numBuiltReq; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := 3, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestWaiterError(t *testing.T) {
	svc := &mockClient{Client: awstesting.NewClient(unit.Config())}
	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.UnmarshalError.Clear()
	svc.Handlers.ValidateResponse.Clear()

	reqNum := 0
	resps := []*MockOutput{
		{ // Request 1
			States: []MockState{
				{State: StateTypePending},
				{State: StateTypePending},
			},
		},
		{ // Request 1, error case retry
		},
		{ // Request 2, error case failure
		},
		{ // Request 3
			States: []MockState{
				{State: StateTypeRunning},
				{State: StateTypeRunning},
			},
		},
	}
	reqErrs := make([]error, len(resps))
	reqErrs[1] = awserr.New("MockException", "mock exception message", nil)
	reqErrs[2] = awserr.New("FailureException", "mock failure exception message", nil)

	numBuiltReq := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		numBuiltReq++
	})
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		code := 200
		if reqNum == 1 {
			code = 400
		}
		r.HTTPResponse = &http.Response{
			StatusCode: code,
			Status:     http.StatusText(code),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		}
	})
	svc.Handlers.Unmarshal.PushBack(func(r *aws.Request) {
		if reqNum >= len(resps) {
			t.Fatalf("too many polling requests made")
			return
		}
		r.Data = resps[reqNum]
		reqNum++
	})
	svc.Handlers.UnmarshalMeta.PushBack(func(r *aws.Request) {
		// If there was an error unmarshal error will be called instead of unmarshal
		// need to increment count here also
		if err := reqErrs[reqNum]; err != nil {
			r.Error = err
			reqNum++
		}
	})

	w := aws.Waiter{
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(0),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.PathAllWaiterMatch,
				Argument: "States[].State",
				Expected: "running",
			},
			{
				State:    aws.RetryWaiterState,
				Matcher:  aws.ErrorWaiterMatch,
				Argument: "",
				Expected: "MockException",
			},
			{
				State:    aws.FailureWaiterState,
				Matcher:  aws.ErrorWaiterMatch,
				Argument: "",
				Expected: "FailureException",
			},
		},
		NewRequest: BuildNewMockRequest(svc, &MockInput{}),
	}

	err := w.WaitWithContext(aws.BackgroundContext())
	if err == nil {
		t.Fatalf("expected error, but did not get one")
	}
	aerr := err.(awserr.Error)
	if e, a := aws.WaiterResourceNotReadyErrorCode, aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := 3, numBuiltReq; e != a {
		t.Errorf("expect %d built requests got %d", e, a)
	}
	if e, a := 3, reqNum; e != a {
		t.Errorf("expect %d reqNum got %d", e, a)
	}
}

func TestWaiterStatus(t *testing.T) {
	svc := &mockClient{Client: awstesting.NewClient(unit.Config())}
	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	reqNum := 0
	svc.Handlers.Build.PushBack(func(r *aws.Request) {
		reqNum++
	})
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		code := 200
		if reqNum == 3 {
			code = 404
			r.Error = awserr.New("NotFound", "Not Found", nil)
		}
		r.HTTPResponse = &http.Response{
			StatusCode: code,
			Status:     http.StatusText(code),
			Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
		}
	})

	w := aws.Waiter{
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(0),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.StatusWaiterMatch,
				Argument: "",
				Expected: 404,
			},
		},
		NewRequest: BuildNewMockRequest(svc, &MockInput{}),
	}

	err := w.WaitWithContext(aws.BackgroundContext())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := 3, reqNum; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestWaiter_ApplyOptions(t *testing.T) {
	w := aws.Waiter{}

	logger := aws.NewDefaultLogger()

	w.ApplyOptions(
		aws.WithWaiterLogger(logger),
		aws.WithWaiterRequestOptions(aws.WithLogLevel(aws.LogDebug)),
		aws.WithWaiterMaxAttempts(2),
		aws.WithWaiterDelay(aws.ConstantWaiterDelay(5*time.Second)),
	)

	if e, a := logger, w.Logger; e != a {
		t.Errorf("expect logger to be set, and match, was not, %v, %v", e, a)
	}

	if len(w.RequestOptions) != 1 {
		t.Fatalf("expect request options to be set to only a single option, %v", w.RequestOptions)
	}
	r := aws.Request{}
	r.ApplyOptions(w.RequestOptions...)
	if e, a := aws.LogDebug, r.Config.LogLevel; e != a {
		t.Errorf("expect %v loglevel got %v", e, a)
	}

	if e, a := 2, w.MaxAttempts; e != a {
		t.Errorf("expect %d retryer max attempts, got %d", e, a)
	}
	if e, a := 5*time.Second, w.Delay(0); e != a {
		t.Errorf("expect %d retryer delay, got %d", e, a)
	}
}

func TestWaiter_WithContextCanceled(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	reqCount := 0

	w := aws.Waiter{
		Name:             "TestWaiter",
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(1 * time.Millisecond),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.StatusWaiterMatch,
				Expected: 200,
			},
		},
		Logger: aws.NewDefaultLogger(),
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
			req.HTTPResponse = &http.Response{StatusCode: http.StatusNotFound}
			req.Handlers.Clear()
			req.Data = struct{}{}
			req.Handlers.Send.PushBack(func(r *aws.Request) {
				if reqCount == 1 {
					ctx.Error = fmt.Errorf("context canceled")
					close(ctx.DoneCh)
				}
				reqCount++
			})

			return req, nil
		},
	}

	w.SleepWithContext = func(c aws.Context, delay time.Duration) error {
		context := c.(*awstesting.FakeContext)
		select {
		case <-context.DoneCh:
			return context.Err()
		default:
			return nil
		}
	}

	err := w.WaitWithContext(ctx)

	if err == nil {
		t.Fatalf("expect waiter to be canceled.")
	}
	aerr := err.(awserr.Error)
	if e, a := aws.ErrCodeRequestCanceled, aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := 2, reqCount; e != a {
		t.Errorf("expect %d requests, got %d", e, a)
	}
}

func TestWaiter_WithContext(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	reqCount := 0

	statuses := []int{http.StatusNotFound, http.StatusOK}

	w := aws.Waiter{
		Name:             "TestWaiter",
		MaxAttempts:      10,
		Delay:            aws.ConstantWaiterDelay(1 * time.Millisecond),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.StatusWaiterMatch,
				Expected: 200,
			},
		},
		Logger: aws.NewDefaultLogger(),
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
			req.HTTPResponse = &http.Response{StatusCode: statuses[reqCount]}
			req.Handlers.Clear()
			req.Data = struct{}{}
			req.Handlers.Send.PushBack(func(r *aws.Request) {
				if reqCount == 1 {
					ctx.Error = fmt.Errorf("context canceled")
					close(ctx.DoneCh)
				}
				reqCount++
			})

			return req, nil
		},
	}

	err := w.WaitWithContext(ctx)

	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := 2, reqCount; e != a {
		t.Errorf("expect %d requests, got %d", e, a)
	}
}

func TestWaiter_AttemptsExpires(t *testing.T) {
	c := awstesting.NewClient(unit.Config())

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	reqCount := 0

	w := aws.Waiter{
		Name:             "TestWaiter",
		MaxAttempts:      2,
		Delay:            aws.ConstantWaiterDelay(1 * time.Millisecond),
		SleepWithContext: aws.SleepWithContext,
		Acceptors: []aws.WaiterAcceptor{
			{
				State:    aws.SuccessWaiterState,
				Matcher:  aws.StatusWaiterMatch,
				Expected: 200,
			},
		},
		Logger: aws.NewDefaultLogger(),
		NewRequest: func(opts []aws.Option) (*aws.Request, error) {
			req := c.NewRequest(&aws.Operation{Name: "Operation"}, nil, nil)
			req.HTTPResponse = &http.Response{StatusCode: http.StatusNotFound}
			req.Handlers.Clear()
			req.Data = struct{}{}
			req.Handlers.Send.PushBack(func(r *aws.Request) {
				reqCount++
			})

			return req, nil
		},
	}

	err := w.WaitWithContext(ctx)

	if err == nil {
		t.Fatalf("expect error did not get one")
	}
	aerr := err.(awserr.Error)
	if e, a := aws.WaiterResourceNotReadyErrorCode, aerr.Code(); e != a {
		t.Errorf("expect %q error code, got %q", e, a)
	}
	if e, a := 2, reqCount; e != a {
		t.Errorf("expect %d requests, got %d", e, a)
	}
}

func TestWaiterNilInput(t *testing.T) {
	orig := sdk.SleepWithContext
	defer func() { sdk.SleepWithContext = orig }()
	sdk.SleepWithContext = func(context.Context, time.Duration) error { return nil }

	// Code generation doesn't have a great way to verify the code is correct
	// other than being run via unit tests in the SDK. This should be fixed
	// So code generation can be validated independently.

	svc := s3.New(unit.Config())
	svc.Handlers.Validate.Clear()
	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: http.StatusOK,
		}
	})
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	// Ensure waiters do not panic on nil input. It doesn't make sense to
	// call a waiter without an input, Validation will
	err := svc.WaitUntilBucketExists(nil)
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
}

func TestWaiterWithContextNilInput(t *testing.T) {
	// Code generation doesn't have a great way to verify the code is correct
	// other than being run via unit tests in the SDK. This should be fixed
	// So code generation can be validated independently.

	svc := s3.New(unit.Config())
	svc.Handlers.Validate.Clear()
	svc.Handlers.Send.Clear() // mock sending
	svc.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: http.StatusOK,
		}
	})
	svc.Handlers.Unmarshal.Clear()
	svc.Handlers.UnmarshalMeta.Clear()
	svc.Handlers.ValidateResponse.Clear()

	// Ensure waiters do not panic on nil input
	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	err := svc.WaitUntilBucketExistsWithContext(ctx, nil,
		aws.WithWaiterDelay(aws.ConstantWaiterDelay(0)),
		aws.WithWaiterMaxAttempts(1),
	)
	if err != nil {
		t.Fatalf("expect no error, but got %v", err)
	}
}
