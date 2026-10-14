package http

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T, h http.HandlerFunc, http2 bool) *httptest.Server {
	t.Helper()

	srv := httptest.NewUnstartedServer(h)
	srv.EnableHTTP2 = http2
	srv.StartTLS()
	t.Cleanup(srv.Close)

	return srv
}

// testClient returns a BuildableClient that trusts srv's certificate with the
// read timeout set explicitly. The package default is a constant and is not
// enforced in public builds, so behavioral tests drive it as a caller would.
func testClient(t *testing.T, srv *httptest.Server, timeout time.Duration) *BuildableClient {
	t.Helper()

	srvTransport, ok := srv.Client().Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expect *http.Transport from test server, got %T", srv.Client().Transport)
	}

	return NewBuildableClient().
		WithTransportOptions(func(tr *http.Transport) {
			tr.TLSClientConfig = srvTransport.TLSClientConfig.Clone()
		}).
		WithReadTimeout(timeout)
}

func TestReadTimeoutBakedIntoDialer(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, false)
	srvTransport := srv.Client().Transport.(*http.Transport)

	for name, tt := range map[string]struct {
		client     func() *BuildableClient
		expectWrap bool
		expect     time.Duration
	}{
		"configured": {
			client: func() *BuildableClient {
				return NewBuildableClient().WithReadTimeout(42 * time.Second)
			},
			expectWrap: true, expect: 42 * time.Second,
		},
		"unset": {
			client:     NewBuildableClient,
			expectWrap: false,
		},
		"explicitly disabled": {
			client: func() *BuildableClient {
				return NewBuildableClient().WithReadTimeout(0)
			},
			expectWrap: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			dialed := make(chan net.Conn, 1)
			client := tt.client().WithTransportOptions(func(tr *http.Transport) {
				tr.TLSClientConfig = srvTransport.TLSClientConfig.Clone()
			})
			client = client.WithTransportOptions(func(tr *http.Transport) {
				client.installReadTimeout(tr)
				inner := tr.DialContext
				tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
					conn, err := inner(ctx, network, addr)
					if err == nil {
						dialed <- conn
					}
					return conn, err
				}
			})

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("do: %v", err)
			}
			resp.Body.Close()

			conn := <-dialed
			wrapped, ok := conn.(*deadlineConn)
			if ok != tt.expectWrap {
				t.Fatalf("expect wrapped=%v, got %T", tt.expectWrap, conn)
			}
			if ok && wrapped.timeout != tt.expect {
				t.Errorf("expect %v, got %v", tt.expect, wrapped.timeout)
			}
		})
	}
}

// TestReadTimeoutPreservesHTTP2 guards the install point: it must be
// DialContext, not DialTLSContext, or HTTP/2 breaks.
func TestReadTimeoutPreservesHTTP2(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}, true)

	client := testClient(t, srv, time.Minute)

	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.Proto != "HTTP/2.0" {
		t.Fatalf("expect HTTP/2.0, got %v: the connection wrapper disabled h2", resp.Proto)
	}
}

func TestReadTimeoutStallMidBody(t *testing.T) {
	for name, http2 := range map[string]bool{"h1": false, "h2": true} {
		t.Run(name, func(t *testing.T) {
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1000")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello"))
				w.(http.Flusher).Flush()
				<-r.Context().Done()
			}, http2)

			client := testClient(t, srv, 150*time.Millisecond)

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("expect response, got error: %v", err)
			}
			defer resp.Body.Close()

			expectProto := "HTTP/1.1"
			if http2 {
				expectProto = "HTTP/2.0"
			}
			if resp.Proto != expectProto {
				t.Fatalf("expect %v, got %v: not exercising the intended protocol",
					expectProto, resp.Proto)
			}

			_, err = io.ReadAll(resp.Body)

			var timeout *ResponseTimeoutError
			if !errors.As(err, &timeout) {
				t.Fatalf("expect *ResponseTimeoutError, got %T: %v", err, err)
			}
			if !timeout.Timeout() {
				t.Error("expect Timeout() true")
			}
			if errors.Is(err, context.Canceled) {
				t.Error("expect a timeout rather than a cancellation, which is not retried")
			}
		})
	}
}

func TestReadTimeoutStallBeforeHeaders(t *testing.T) {
	for name, http2 := range map[string]bool{"h1": false, "h2": true} {
		t.Run(name, func(t *testing.T) {
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				<-r.Context().Done()
			}, http2)

			client := testClient(t, srv, 150*time.Millisecond)

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := client.Do(req)
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}

			var timeout *ResponseTimeoutError
			if !errors.As(err, &timeout) {
				t.Fatalf("expect *ResponseTimeoutError, got %T: %v", err, err)
			}
		})
	}
}

func TestReadTimeoutSlowButProgressingSucceeds(t *testing.T) {
	const chunks = 10

	for name, http2 := range map[string]bool{"h1": false, "h2": true} {
		t.Run(name, func(t *testing.T) {
			srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				for i := 0; i < chunks; i++ {
					w.Write([]byte("x"))
					w.(http.Flusher).Flush()
					time.Sleep(30 * time.Millisecond)
				}
			}, http2)

			// The whole transfer exceeds the window; no single gap does.
			client := testClient(t, srv, 200*time.Millisecond)

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("expect response, got error: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("expect the window to reset on progress, got error: %v", err)
			}
			if len(body) != chunks {
				t.Errorf("expect %v bytes, got %v", chunks, len(body))
			}
		})
	}
}

// TestReadTimeoutFreezeStillApplies covers a caller-configured read timeout
// surviving Freeze. Freeze's contract is that resolveHTTPClient can no longer
// reach the client to apply a default (the type assertion to *BuildableClient
// fails), not that a value the caller set themselves gets discarded.
func TestReadTimeoutFreezeStillApplies(t *testing.T) {
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
		w.(http.Flusher).Flush()
		time.Sleep(400 * time.Millisecond)
		w.Write(make([]byte, 995))
	}, false)

	frozen := testClient(t, srv, 100*time.Millisecond).Freeze()

	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := frozen.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	defer resp.Body.Close()

	// The stall is four times the configured window; the caller's own read
	// timeout still fires even though the client is frozen.
	_, err = io.ReadAll(resp.Body)
	var timeout *ResponseTimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("expect *ResponseTimeoutError, got %T: %v", err, err)
	}
}

// TestBuilderOptionsPreserveReadTimeout guards against clone() dropping the
// configured value regardless of builder method order.
func TestBuilderOptionsPreserveReadTimeout(t *testing.T) {
	for name, client := range map[string]*BuildableClient{
		"timeout then transport": NewBuildableClient().
			WithReadTimeout(time.Second).WithTransportOptions(func(*http.Transport) {}),
		"transport then timeout": NewBuildableClient().
			WithTransportOptions(func(*http.Transport) {}).WithReadTimeout(time.Second),
		"timeout then dialer": NewBuildableClient().
			WithReadTimeout(time.Second).WithDialerOptions(func(*net.Dialer) {}),
		"timeout then client timeout": NewBuildableClient().
			WithReadTimeout(time.Second).WithTimeout(time.Minute),
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := client.GetReadTimeout()
			if !ok || got != time.Second {
				t.Errorf("expect 1s/true, got %v/%v", got, ok)
			}
		})
	}
}

func TestGetReadTimeoutUnset(t *testing.T) {
	if _, ok := NewBuildableClient().GetReadTimeout(); ok {
		t.Error("expect read timeout unset")
	}
}

// TestWithReadTimeoutDoesNotShareConnections guards WithReadTimeout always
// cloning the receiver rather than mutating it in place. Connections are
// pooled per transport and keyed by address, not by service, so a shared
// client that got mutated could let one service inherit a connection carrying
// another service's deadline.
func TestWithReadTimeoutDoesNotShareConnections(t *testing.T) {
	// Stalls long enough to trip the timed client, and only after headers so
	// both clients get a response first.
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
		w.(http.Flusher).Flush()
		time.Sleep(400 * time.Millisecond)
		w.Write(make([]byte, 995))
	}, false)
	srvTransport := srv.Client().Transport.(*http.Transport)

	shared := NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		tr.TLSClientConfig = srvTransport.TLSClientConfig.Clone()
	})

	// Stands in for two generated resolveHTTPClient calls resolving against
	// the same starting client: one leaves it as-is, one applies a timeout.
	exempt := shared
	timed := shared.WithReadTimeout(100 * time.Millisecond)

	if timed == shared {
		t.Fatal("expect WithReadTimeout to return a distinct client")
	}

	read := func(client *BuildableClient) error {
		req, err := http.NewRequest("GET", srv.URL, nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		_, err = io.ReadAll(resp.Body)
		return err
	}

	// Interleaved, so a shared pool would hand the exempt client a connection
	// the timed client dialed.
	for i := 0; i < 3; i++ {
		var timeout *ResponseTimeoutError
		if err := read(timed); !errors.As(err, &timeout) {
			t.Fatalf("iteration %d: expect the timed client to time out, got %v", i, err)
		}
		if err := read(exempt); err != nil {
			t.Fatalf("iteration %d: expect the untimed client to succeed, got %v", i, err)
		}
	}
}
