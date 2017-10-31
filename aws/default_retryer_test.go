package aws

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryThrottleStatusCodes(t *testing.T) {
	cases := []struct {
		expectThrottle bool
		expectRetry    bool
		r              Request
	}{
		{
			false,
			false,
			Request{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 429},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 502},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 503},
			},
		},
		{
			true,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 504},
			},
		},
		{
			false,
			true,
			Request{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
		},
	}

	d := DefaultRetryer{NumMaxRetries: 10}
	for i, c := range cases {
		throttle := d.shouldThrottle(&c.r)
		retry := d.ShouldRetry(&c.r)

		if e, a := c.expectThrottle, throttle; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}

		if e, a := c.expectRetry, retry; e != a {
			t.Errorf("%d: expected %v, but received %v", i, e, a)
		}
	}
}

func TestCanUseRetryAfter(t *testing.T) {
	cases := []struct {
		r Request
		e bool
	}{
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 200},
			},
			false,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 500},
			},
			false,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 429},
			},
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503},
			},
			true,
		},
	}

	for i, c := range cases {
		a := canUseRetryAfterHeader(&c.r)
		if c.e != a {
			t.Errorf("%d: expected %v, but received %v", i, c.e, a)
		}
	}
}

func TestGetRetryDelay(t *testing.T) {
	cases := []struct {
		r     Request
		e     time.Duration
		equal bool
		ok    bool
	}{
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 429, Header: http.Header{"Retry-After": []string{"3600"}}},
			},
			3600 * time.Second,
			true,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{"120"}}},
			},
			120 * time.Second,
			true,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{"120"}}},
			},
			1 * time.Second,
			false,
			true,
		},
		{
			Request{
				HTTPResponse: &http.Response{StatusCode: 503, Header: http.Header{"Retry-After": []string{""}}},
			},
			0 * time.Second,
			true,
			false,
		},
	}

	for i, c := range cases {
		a, ok := getRetryDelay(&c.r)
		if c.ok != ok {
			t.Errorf("%d: expected %v, but received %v", i, c.ok, ok)
		}

		if (c.e != a) == c.equal {
			t.Errorf("%d: expected %v, but received %v", i, c.e, a)
		}
	}
}
