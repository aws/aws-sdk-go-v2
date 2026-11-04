package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	smithymiddleware "github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// These tests replay the cross-SDK Clock Skew Correction SEP test vectors
// (testdata/clock-skew-test-cases.json) against the SDK's skew computation. For
// each operation it seeds AttemptSkew from ClientSkew, signs each attempt at
// clientTimeAtSend+AttemptSkew (asserting expectedSigningTime), computes the
// candidate skew via computeClockSkew, records it unconditionally when present,
// and asserts the resulting ClientSkew (expectedClientSkew).
type skewVectorFile struct {
	Tests []skewVectorTest `json:"tests"`
}

type skewVectorTest struct {
	Description string         `json:"description"`
	Operations  []skewVectorOp `json:"operations"`
}

type skewVectorOp struct {
	InitialClientSkew  int64               `json:"initialClientSkew"`
	MaxAttempts        int                 `json:"maxAttempts"`
	Attempts           []skewVectorAttempt `json:"attempts"`
	ExpectedClientSkew int64               `json:"expectedClientSkew"`
	ExpectedOutcome    string              `json:"expectedOutcome"`
}

type skewVectorAttempt struct {
	ClientTimeAtSend    string         `json:"clientTimeAtSend"`
	ClientTimeAtReceive string         `json:"clientTimeAtReceive"`
	ExpectedSigningTime string         `json:"expectedSigningTime"`
	Response            skewVectorResp `json:"response"`
}

type skewVectorResp struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	ErrorCode  string            `json:"errorCode"`
}

func TestClockSkewSEPVectors(t *testing.T) {
	raw, err := os.ReadFile("testdata/clock-skew-test-cases.json")
	if err != nil {
		t.Fatalf("read vectors: %v", err)
	}

	var file skewVectorFile
	if err := json.Unmarshal(raw, &file); err != nil {
		t.Fatalf("parse vectors: %v", err)
	}

	if len(file.Tests) == 0 {
		t.Fatal("no test vectors loaded")
	}

	for _, tc := range file.Tests {
		t.Run(tc.Description, func(t *testing.T) {
			var prevExpected int64
			for oi, op := range tc.Operations {
				if oi > 0 && op.InitialClientSkew != prevExpected {
					t.Fatalf("op %d: initialClientSkew %d does not match prior expectedClientSkew %d",
						oi, op.InitialClientSkew, prevExpected)
				}

				clientSkew := time.Duration(op.InitialClientSkew) * time.Second
				attemptSkew := clientSkew

				for ai, at := range op.Attempts {
					send := mustParseISO(t, at.ClientTimeAtSend)
					recv := mustParseISO(t, at.ClientTimeAtReceive)

					signing := send.Add(attemptSkew)
					if want := mustParseISO(t, at.ExpectedSigningTime); !signing.Equal(want) {
						t.Fatalf("op %d attempt %d: signing time = %s, want %s", oi, ai, signing.Format(time.RFC3339), want.Format(time.RFC3339))
					}

					var serverTime time.Time
					if d := at.Response.Headers["Date"]; d != "" {
						st, perr := smithyhttp.ParseTime(d)
						if perr != nil {
							t.Fatalf("op %d attempt %d: parse Date %q: %v", oi, ai, d, perr)
						}
						serverTime = st
					}
					hasAge := at.Response.Headers["Age"] != ""

					if skew, ok := computeClockSkew(serverTime, send, recv, hasAge); ok {
						attemptSkew = skew
						clientSkew = skew
					}
				}

				if want := time.Duration(op.ExpectedClientSkew) * time.Second; clientSkew != want {
					t.Fatalf("op %d: ClientSkew = %s, want %s", oi, clientSkew, want)
				}
				prevExpected = op.ExpectedClientSkew
			}
		})
	}
}

func mustParseISO(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parse time %q: %v", s, err)
	}

	return v
}

// TestRecordResponseTimingClockSkewDisabled verifies that the disable flag
// suppresses skew recording while a default (enabled) middleware records it.
func TestRecordResponseTimingClockSkewDisabled(t *testing.T) {
	responseAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	restore := sdk.TestingUseReferenceTime(responseAt)
	defer restore()

	dateStr := responseAt.Add(5 * time.Minute).UTC().Format(http.TimeFormat)
	next := smithymiddleware.DeserializeHandlerFunc(func(ctx context.Context, in smithymiddleware.DeserializeInput) (
		out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error,
	) {
		out.RawResponse = &smithyhttp.Response{Response: &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Date": []string{dateStr}},
		}}
		return out, m, nil
	})

	_, mdDisabled, err := RecordResponseTiming{DisableClockSkewCorrection: true}.
		HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, next)
	if err != nil {
		t.Fatalf("disabled: unexpected error: %v", err)
	}
	if skew, ok := GetAttemptSkew(mdDisabled); ok {
		t.Fatalf("disabled: expected no recorded skew, got %v", skew)
	}

	_, mdEnabled, err := RecordResponseTiming{}.
		HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, next)
	if err != nil {
		t.Fatalf("enabled: unexpected error: %v", err)
	}
	if skew, ok := GetAttemptSkew(mdEnabled); !ok || skew != 5*time.Minute {
		t.Fatalf("enabled: expected 5m recorded skew, got %v ok=%v", skew, ok)
	}
}

// TestRecordResponseTimingSkewCases drives RecordResponseTiming with a stubbed
// clock returning send for the request time and recv for the response time.
func TestRecordResponseTimingSkewCases(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	at := func(sec int) time.Time { return base.Add(time.Duration(sec) * time.Second) }

	cases := map[string]struct {
		send, recv time.Time
		headers    http.Header
		wantSkew   time.Duration
		wantOK     bool
	}{
		"clocks agree": {
			send: at(0), recv: at(1),
			headers:  http.Header{"Date": []string{at(1).Format(http.TimeFormat)}},
			wantSkew: 0, wantOK: true,
		},
		"server ahead one hour": {
			send: at(0), recv: at(2),
			headers:  http.Header{"Date": []string{at(3602).Format(http.TimeFormat)}},
			wantSkew: 3600 * time.Second, wantOK: true,
		},
		"server behind one hour": {
			send: at(0), recv: at(2),
			headers:  http.Header{"Date": []string{at(-3598).Format(http.TimeFormat)}},
			wantSkew: -3600 * time.Second, wantOK: true,
		},
		"slow request clocks agree": {
			send: at(0), recv: at(620),
			headers:  http.Header{"Date": []string{at(620).Format(http.TimeFormat)}},
			wantSkew: 0, wantOK: true,
		},
		"slow request server ahead": {
			send: at(0), recv: at(620),
			headers:  http.Header{"Date": []string{at(1220).Format(http.TimeFormat)}},
			wantSkew: 600 * time.Second, wantOK: true,
		},
		"max trusted duration boundary": {
			send: at(0), recv: at(900),
			headers:  http.Header{"Date": []string{at(900).Format(http.TimeFormat)}},
			wantSkew: 0, wantOK: true,
		},
		"over max trusted duration": {
			send: at(0), recv: at(901),
			headers: http.Header{"Date": []string{at(901).Format(http.TimeFormat)}},
		},
		"age header": {
			send: at(0), recv: at(2),
			headers: http.Header{
				"Date": []string{at(3602).Format(http.TimeFormat)},
				"Age":  []string{"30"},
			},
		},
		"no date header": {
			send: at(0), recv: at(2),
			headers: http.Header{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			orig := sdk.NowTime
			defer func() { sdk.NowTime = orig }()
			calls := 0
			sdk.NowTime = func() time.Time {
				calls++
				if calls == 1 {
					return c.send
				}
				return c.recv
			}

			next := smithymiddleware.DeserializeHandlerFunc(func(ctx context.Context, in smithymiddleware.DeserializeInput) (
				out smithymiddleware.DeserializeOutput, m smithymiddleware.Metadata, err error,
			) {
				out.RawResponse = &smithyhttp.Response{Response: &http.Response{
					StatusCode: 200,
					Header:     c.headers,
				}}
				return out, m, nil
			})

			_, md, err := RecordResponseTiming{}.HandleDeserialize(context.Background(), smithymiddleware.DeserializeInput{}, next)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if calls != 2 {
				t.Errorf("expected 2 clock reads, got %d", calls)
			}
			if v, ok := GetResponseAt(md); !ok || !v.Equal(c.recv) {
				t.Errorf("expected response at %v, got %v ok=%v", c.recv, v, ok)
			}
			skew, ok := GetAttemptSkew(md)
			if ok != c.wantOK {
				t.Fatalf("expected skew recorded=%v, got %v (skew %v)", c.wantOK, ok, skew)
			}
			if ok && skew != c.wantSkew {
				t.Errorf("expected skew %v, got %v", c.wantSkew, skew)
			}
		})
	}
}

// runSkewServer sends one request through a stack with RecordResponseTiming
// in the deserialize step against a local server.
func runSkewServer(t *testing.T, handler http.HandlerFunc, body []byte, expectContinue bool) (smithymiddleware.Metadata, time.Time) {
	t.Helper()

	srv := httptest.NewServer(handler)
	defer srv.Close()

	tr := srv.Client().Transport.(*http.Transport).Clone()
	tr.ExpectContinueTimeout = 5 * time.Second
	client := &http.Client{Transport: tr}

	stack := smithymiddleware.NewStack("test", smithyhttp.NewStackRequest)
	err := stack.Serialize.Add(smithymiddleware.SerializeMiddlewareFunc("setURL",
		func(ctx context.Context, in smithymiddleware.SerializeInput, next smithymiddleware.SerializeHandler) (
			smithymiddleware.SerializeOutput, smithymiddleware.Metadata, error,
		) {
			req := in.Request.(*smithyhttp.Request)
			req.URL.Scheme = "http"
			req.URL.Host = srv.Listener.Addr().String()
			if body != nil {
				req.Method = "PUT"
				req.ContentLength = int64(len(body))
				if expectContinue {
					req.Header.Set("Expect", "100-continue")
				}
				r, err := req.SetStream(bytes.NewReader(body))
				if err != nil {
					return smithymiddleware.SerializeOutput{}, smithymiddleware.Metadata{}, err
				}
				in.Request = r
			}
			return next.HandleSerialize(ctx, in)
		}), smithymiddleware.After)
	if err != nil {
		t.Fatal(err)
	}
	if err := AddRawResponseToMetadata(stack); err != nil {
		t.Fatal(err)
	}
	if err := stack.Deserialize.Add(&RecordResponseTiming{}, smithymiddleware.After); err != nil {
		t.Fatal(err)
	}

	h := smithymiddleware.DecorateHandler(smithyhttp.NewClientHandler(client), stack)
	_, md, err := h.Handle(context.Background(), struct{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Read the body after the stack returns, as the deserializer does.
	if resp, ok := GetRawResponse(md).(*smithyhttp.Response); ok {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	return md, time.Now()
}

func assertSkewNearHour(t *testing.T, md smithymiddleware.Metadata) {
	t.Helper()
	skew, ok := GetAttemptSkew(md)
	if !ok {
		t.Fatal("expected skew to be recorded")
	}
	if skew < time.Hour-1500*time.Millisecond || skew > time.Hour+500*time.Millisecond {
		t.Errorf("expected skew near 1h, got %v", skew)
	}
}

// TestRecordResponseTimingSkewFromHeaderArrival verifies the response time is
// taken when Do returns with the headers, before a slow body is read.
func TestRecordResponseTimingSkewFromHeaderArrival(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow real server test in short mode")
	}

	md, end := runSkewServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		time.Sleep(4 * time.Second)
		w.Write([]byte("response body"))
	}, nil, false)

	responseAt, ok := GetResponseAt(md)
	if !ok {
		t.Fatal("expected response at to be recorded")
	}
	if d := end.Sub(responseAt); d < 3*time.Second {
		t.Errorf("expected response at to precede the body by at least 3s, got %v", d)
	}
	assertSkewNearHour(t, md)
}

// TestRecordResponseTimingSkewWith100Continue verifies that a "100 Continue"
// interim response sent before the body upload is not used as the response
// time.
func TestRecordResponseTimingSkewWith100Continue(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow real server test in short mode")
	}

	var gotExpect string
	var sleptAt time.Time
	md, _ := runSkewServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotExpect = r.Header.Get("Expect")
		io.Copy(io.Discard, r.Body)
		time.Sleep(3 * time.Second)
		sleptAt = time.Now()
		w.Header().Set("Date", time.Now().Add(time.Hour).UTC().Format(http.TimeFormat))
		w.WriteHeader(http.StatusOK)
	}, make([]byte, 4<<20), true)

	if gotExpect != "100-continue" {
		t.Fatalf("expected Expect 100-continue request, got %q", gotExpect)
	}
	responseAt, ok := GetResponseAt(md)
	if !ok {
		t.Fatal("expected response at to be recorded")
	}
	if responseAt.Before(sleptAt) {
		t.Errorf("expected response at %v not before final response %v", responseAt, sleptAt)
	}
	assertSkewNearHour(t, md)
}
