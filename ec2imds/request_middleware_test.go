package ec2imds

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
	"github.com/google/go-cmp/cmp"
)

func TestAddRequestMiddleware(t *testing.T) {
	cases := map[string]struct {
		AddMiddleware     func(*middleware.Stack, Options) error
		ExpectInitialize  []string
		ExpectSerialize   []string
		ExpectBuild       []string
		ExpectFinalize    []string
		ExpectDeserialize []string
	}{
		"api request": {
			AddMiddleware: func(stack *middleware.Stack, options Options) error {
				return addAPIRequestMiddleware(stack, options,
					func(interface{}) (string, error) {
						return "/mockPath", nil
					},
					func(*smithyhttp.Response) (interface{}, error) {
						return struct{}{}, nil
					},
				)
			},
			ExpectInitialize: []string{
				(*operationTimeoutMiddleware)(nil).ID(),
			},
			ExpectSerialize: []string{
				"ResolveEndpoint",
				"OperationSerializer",
			},
			ExpectBuild: []string{
				"UserAgent",
			},
			ExpectFinalize: []string{
				"Retry",
				"APITokenProvider",
				"RetryMetricsHeader",
			},
			ExpectDeserialize: []string{
				"APITokenProvider",
				"OperationDeserializer",
			},
		},

		"base request": {
			AddMiddleware: func(stack *middleware.Stack, options Options) error {
				return addRequestMiddleware(stack, options, "POST",
					func(interface{}) (string, error) {
						return "/mockPath", nil
					},
					func(*smithyhttp.Response) (interface{}, error) {
						return struct{}{}, nil
					},
				)
			},
			ExpectInitialize: []string{
				(*operationTimeoutMiddleware)(nil).ID(),
			},
			ExpectSerialize: []string{
				"ResolveEndpoint",
				"OperationSerializer",
			},
			ExpectBuild: []string{
				"UserAgent",
			},
			ExpectFinalize: []string{
				"Retry",
				"RetryMetricsHeader",
			},
			ExpectDeserialize: []string{
				"OperationDeserializer",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			client := New(Options{})

			stack := middleware.NewStack("mockOp", smithyhttp.NewStackRequest)

			if err := c.AddMiddleware(stack, client.options); err != nil {
				t.Fatalf("expect no error adding middleware, got %v", err)
			}

			if diff := cmp.Diff(c.ExpectInitialize, stack.Initialize.List()); len(diff) != 0 {
				t.Errorf("expect initialize middleware\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectSerialize, stack.Serialize.List()); len(diff) != 0 {
				t.Errorf("expect serialize middleware\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectBuild, stack.Build.List()); len(diff) != 0 {
				t.Errorf("expect build middleware\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectFinalize, stack.Finalize.List()); len(diff) != 0 {
				t.Errorf("expect finalize middleware\n%s", diff)
			}

			if diff := cmp.Diff(c.ExpectDeserialize, stack.Deserialize.List()); len(diff) != 0 {
				t.Errorf("expect deserialize middleware\n%s", diff)
			}
		})
	}
}

func TestOperationTimeoutMiddleware(t *testing.T) {
	m := &operationTimeoutMiddleware{
		Timeout: time.Nanosecond,
	}

	_, _, err := m.HandleInitialize(context.Background(), middleware.InitializeInput{},
		middleware.InitializeHandlerFunc(func(
			ctx context.Context, input middleware.InitializeInput,
		) (
			out middleware.InitializeOutput, metadata middleware.Metadata, err error,
		) {
			if err := sdk.SleepWithContext(ctx, time.Second); err != nil {
				return out, metadata, err
			}

			return out, metadata, nil
		}))
	if err == nil {
		t.Fatalf("expect error got none")
	}

	if e, a := "deadline exceeded", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expect %q error in %q", e, a)
	}
}

type successAPIResponseHandler struct {
	t      *testing.T
	path   string
	method string

	// response
	statusCode int
	header     http.Header
	body       []byte
}

func (h *successAPIResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e, a := h.path, r.URL.Path; e != a {
		h.t.Errorf("expect %v path, got %v", e, a)
	}
	if e, a := h.method, r.Method; e != a {
		h.t.Errorf("expect %v method, got %v", e, a)
	}

	for k, vs := range h.header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	if h.statusCode != 0 {
		w.WriteHeader(h.statusCode)
	}
	w.Write(h.body)
}

func TestRequestGetToken(t *testing.T) {
	cases := map[string]struct {
		GetHandler    func(*testing.T) http.Handler
		APICallCount  int
		ExpectTrace   []string
		ExpectContent []byte
		ExpectErr     string
	}{
		"secure": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				"/latest/foo",
			},
			APICallCount: 2,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newSecureAPIHandler(t,
						[]string{"tokenA"},
						5*time.Minute,
						&successAPIResponseHandler{t: t,
							path:   "/latest/foo",
							method: "GET",
							body:   []byte("hello"),
						},
					))
			},
			ExpectContent: []byte("hello"),
		},

		"secure multi token": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				getTokenPath,
				"/latest/foo",
				getTokenPath,
				"/latest/foo",
				getTokenPath,
				"/latest/foo",
			},
			APICallCount: 4,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newSecureAPIHandler(t,
						[]string{"tokenA", "tokenB", "tokenC"},
						1,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							h := &successAPIResponseHandler{t: t,
								path:   "/latest/foo",
								method: "GET",
								body:   []byte("hello"),
							}

							time.Sleep(100 * time.Millisecond)
							h.ServeHTTP(w, r)
						}),
					))
			},
			ExpectContent: []byte("hello"),
		},

		// disables API token, fallback to insecure API calls.
		"insecure 405": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				"/latest/foo",
			},
			APICallCount: 2,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newInsecureAPIHandler(t,
						405,
						&successAPIResponseHandler{t: t,
							path:   "/latest/foo",
							method: "GET",
							body:   []byte("hello"),
						},
					))
			},
			ExpectContent: []byte("hello"),
		},

		"insecure 404": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				"/latest/foo",
			},
			APICallCount: 2,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newInsecureAPIHandler(t,
						404,
						&successAPIResponseHandler{t: t,
							path:   "/latest/foo",
							method: "GET",
							body:   []byte("hello"),
						},
					))
			},
			ExpectContent: []byte("hello"),
		},

		"insecure 403": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				"/latest/foo",
			},
			APICallCount: 2,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newInsecureAPIHandler(t,
						403,
						&successAPIResponseHandler{t: t,
							path:   "/latest/foo",
							method: "GET",
							body:   []byte("hello"),
						},
					))
			},
			ExpectContent: []byte("hello"),
		},

		// Token disabled and becomes re-enabled
		"unauthorized 401 re-enable": {
			ExpectTrace: []string{
				getTokenPath,
				"/latest/foo",
				getTokenPath,
				"/latest/foo",
				"/latest/foo",
			},
			APICallCount: 2,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newUnauthorizedAPIHandler(t,
						newSecureAPIHandler(t,
							[]string{"tokenA"},
							5*time.Minute,
							&successAPIResponseHandler{t: t,
								path:   "/latest/foo",
								method: "GET",
								body:   []byte("hello"),
							},
						)))
			},
			ExpectContent: []byte("hello"),
		},

		// Token and API call both fail
		"bad request 400": {
			ExpectTrace: []string{
				getTokenPath,
			},
			APICallCount: 1,
			GetHandler: func(t *testing.T) http.Handler {
				return newTestServeMux(t,
					newInsecureAPIHandler(t,
						400,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							t.Errorf("expected no call to API handler")
							http.Error(w, "", 400)
						}),
					))
			},
			ExpectContent: []byte("hello"),
			ExpectErr:     "EC2 IMDS failed",
		},
	}

	type mockRequestOutput struct {
		Content io.ReadCloser
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			envs := awstesting.StashEnv()
			defer awstesting.PopEnv(envs)

			trace := newRequestTrace()
			server := httptest.NewServer(trace.WrapHandler(c.GetHandler(t)))
			defer server.Close()

			client := New(Options{
				Endpoint: server.URL,
			})

			ctx := context.Background()
			var result interface{}
			var err error
			for i := 0; i < c.APICallCount; i++ {
				result, _, err = client.invokeOperation(ctx, "TestRequest", struct{}{}, nil,
					func(stack *middleware.Stack, options Options) error {
						return addAPIRequestMiddleware(stack,
							client.options.Copy(),
							func(interface{}) (string, error) {
								return "/latest/foo", nil
							},
							func(resp *smithyhttp.Response) (interface{}, error) {
								return &mockRequestOutput{
									Content: resp.Body,
								}, nil
							},
						)
					},
				)
			}
			if diff := cmp.Diff(c.ExpectTrace, trace.requests); len(diff) != 0 {
				t.Errorf("expect trace to match\n%s", diff)
			}

			if len(c.ExpectErr) != 0 {
				if err == nil {
					t.Fatalf("expect error, got none")
				}
				if e, a := c.ExpectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect error to contain %v, got %v", e, a)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			out, ok := result.(*mockRequestOutput)
			if !ok {
				t.Fatalf("expect output result, got %T", result)
			}

			content, err := ioutil.ReadAll(out.Content)
			if err != nil {
				t.Fatalf("expect to read result, got %v", err)
			}

			if e, a := c.ExpectContent, content; !bytes.Equal(e, a) {
				t.Errorf("expect results to match\nexpect:\n%s\nactual:\n%s",
					hex.Dump(e), hex.Dump(a))
			}
		})
	}
}
