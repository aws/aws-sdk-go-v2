package defaults_test

/*
func TestRequestInvocationIDHeaderHandler(t *testing.T) {
	cfg := unit.Config()
	cfg.Handlers.Send.Clear()

	var invokeID string
	cfg.Handlers.Build.PushBack(func(r *aws.Request) {
		invokeID = r.InvocationID
		if len(invokeID) == 0 {
			t.Fatalf("expect non-empty invocation id")
		}
	})
	cfg.Handlers.Send.PushBack(func(r *aws.Request) {
		if e, a := invokeID, r.InvocationID; e != a {
			t.Errorf("expect %v invoke ID, got %v", e, a)
		}
		r.Error = &aws.RequestSendError{Err: io.ErrUnexpectedEOF}
	})
	retryer := retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = 3
	})
	r := aws.New(cfg, aws.Metadata{}, cfg.Handlers, retryer, &aws.Operation{},
		&struct{}{}, struct{}{})

	if len(r.InvocationID) == 0 {
		t.Fatalf("expect invocation id, got none")
	}

	err := r.Send()
	if err == nil {
		t.Fatalf("expect error got on")
	}
	var maxErr *aws.MaxAttemptsError
	if !errors.As(err, &maxErr) {
		t.Fatalf("expect max errors, got %v", err)
	} else {
		if e, a := 3, maxErr.Attempt; e != a {
			t.Errorf("expect %v attempts, got %v", e, a)
		}
	}
	if len(invokeID) == 0 {
		t.Fatalf("expect non-empty invocation id")
	}

	if e, a := r.InvocationID, r.HTTPRequest.Header.Get("amz-sdk-invocation-id"); e != a {
		t.Errorf("expect %v invocation id, got %v", e, a)
	}
}
*/

/*
func TestRetryMetricHeaderHandler(t *testing.T) {
	nowTime := sdk.NowTime
	defer func() {
		sdk.NowTime = nowTime
	}()
	sdk.NowTime = func() time.Time {
		return time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC)
	}

	cases := map[string]struct {
		Attempt           int
		MaxAttempts       int
		Client            aws.HTTPClient
		ContextDeadline   time.Time
		AttemptClockSkews []time.Duration
		Expect            string
	}{
		"first attempt": {
			Attempt: 1, MaxAttempts: 3,
			Expect: "attempt=1; max=3",
		},
		"last attempt": {
			Attempt: 3, MaxAttempts: 3,
			Expect: "attempt=3; max=3",
		},
		"no max attempt": {
			Attempt: 10,
			Expect:  "attempt=10",
		},
		"with ttl client timeout": {
			Attempt: 2, MaxAttempts: 3,
			AttemptClockSkews: []time.Duration{
				10 * time.Second,
			},
			Client: func() aws.HTTPClient {
				c := &aws.BuildableHTTPClient{}
				return c.WithTimeout(10 * time.Second)
			}(),
			Expect: "attempt=2; max=3; ttl=20200202T000020Z",
		},
		"with ttl context deadline": {
			Attempt: 1, MaxAttempts: 3,
			AttemptClockSkews: []time.Duration{
				10 * time.Second,
			},
			ContextDeadline: sdk.NowTime().Add(10 * time.Second),
			Expect:          "attempt=1; max=3; ttl=20200202T000020Z",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := unit.Config()
			if c.Client != nil {
				cfg.HTTPClient = c.Client
			}
			r := aws.New(cfg, aws.Metadata{}, cfg.Handlers, aws.NoOpRetryer{},
				&aws.Operation{}, &struct{}{}, struct{}{})
			if !c.ContextDeadline.IsZero() {
				ctx, cancel := context.WithDeadline(r.Context(), c.ContextDeadline)
				defer cancel()
				r.SetContext(ctx)
			}

			r.AttemptNum = c.Attempt
			r.AttemptClockSkews = c.AttemptClockSkews
			r.Retryer = retry.AddWithMaxAttempts(r.Retryer, c.MaxAttempts)

			defaults.RetryMetricHeaderHandler.Fn(r)
			if r.Error != nil {
				t.Fatalf("expect no error, got %v", r.Error)
			}

			if e, a := c.Expect, r.HTTPRequest.Header.Get("amz-sdk-request"); e != a {
				t.Errorf("expect %q metric, got %q", e, a)
			}
		})
	}
}
*/
