package aws

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"runtime"
	"testing"
)

type mockCloser struct {
	closed bool
}

func (closer *mockCloser) Read(b []byte) (int, error) {
	return 0, io.EOF
}

func (closer *mockCloser) Close() error {
	closer.closed = true
	return nil
}

func TestTeeReaderCloser(t *testing.T) {
	expected := "FOO"
	buf := bytes.NewBuffer([]byte(expected))
	lw := bytes.NewBuffer(nil)
	c := &mockCloser{}
	closer := teeReaderCloser{
		io.TeeReader(buf, lw),
		c,
	}

	b := make([]byte, len(expected))
	_, err := closer.Read(b)
	closer.Close()

	if expected != lw.String() {
		t.Errorf("Expected %q, but received %q", expected, lw.String())
	}

	if err != nil {
		t.Errorf("Expected 'nil', but received %v", err)
	}

	if !c.closed {
		t.Error("Expected 'true', but received 'false'")
	}
}

func TestLogWriter(t *testing.T) {
	expected := "FOO"
	lw := &logWriter{nil, bytes.NewBuffer(nil)}
	lw.Write([]byte(expected))

	if expected != lw.buf.String() {
		t.Errorf("Expected %q, but received %q", expected, lw.buf.String())
	}
}

func TestLogRequest(t *testing.T) {
	cases := []struct {
		Body       io.ReadSeeker
		ExpectBody []byte
		LogLevel   LogLevel
	}{
		{
			Body:       ReadSeekCloser(bytes.NewBuffer([]byte("body content"))),
			ExpectBody: []byte("body content"),
		},
		{
			Body:       ReadSeekCloser(bytes.NewBuffer([]byte("body content"))),
			LogLevel:   LogDebugWithHTTPBody,
			ExpectBody: []byte("body content"),
		},
		{
			Body:       bytes.NewReader([]byte("body content")),
			ExpectBody: []byte("body content"),
		},
		{
			Body:       bytes.NewReader([]byte("body content")),
			LogLevel:   LogDebugWithHTTPBody,
			ExpectBody: []byte("body content"),
		},
	}

	for i, c := range cases {
		var logW bytes.Buffer
		req := New(
			Config{
				EndpointResolver: ResolveWithEndpointURL("https://endpoint"),
				Credentials:      AnonymousCredentials,
				Logger:           &bufLogger{w: &logW},
				LogLevel:         c.LogLevel,
				Region:           "mock-region",
			},
			Metadata{
				EndpointsID: "https://mock-service.mock-region.amazonaws.com",
			},
			testHandlers(),
			nil,
			&Operation{
				Name:       "APIName",
				HTTPMethod: "POST",
				HTTPPath:   "/",
			},
			struct{}{}, nil,
		)
		req.SetReaderBody(c.Body)
		req.Build()

		logRequest(req)

		b, err := ioutil.ReadAll(req.HTTPRequest.Body)
		if err != nil {
			t.Fatalf("%d, expect to read SDK request Body", i)
		}

		if e, a := c.ExpectBody, b; !reflect.DeepEqual(e, a) {
			t.Errorf("%d, expect %v body, got %v", i, e, a)
		}
	}
}

type bufLogger struct {
	w *bytes.Buffer
}

func (l *bufLogger) Log(args ...interface{}) {
	fmt.Fprintln(l.w, args...)
}

func testHandlers() Handlers {
	var handlers Handlers
	handler := NamedHandler{
		Name: "core.SDKVersionUserAgentHandler",
		Fn: MakeAddToUserAgentHandler(SDKName, SDKVersion,
			runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
	handlers.Build.PushBackNamed(handler)

	return handlers
}
