package aws

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type testReader struct {
	duration time.Duration
	count    int
}

func (r *testReader) Read(b []byte) (int, error) {
	if r.count > 0 {
		r.count--
		return len(b), nil
	}
	time.Sleep(r.duration)
	return 0, io.EOF
}

func (r *testReader) Close() error {
	return nil
}

func TestTimeoutReadCloser(t *testing.T) {
	reader := timeoutReadCloser{
		reader: &testReader{
			duration: time.Second,
			count:    5,
		},
		duration: time.Millisecond,
	}
	b := make([]byte, 100)
	_, err := reader.Read(b)
	if err != nil {
		t.Log(err)
	}
}

func TestTimeoutReadCloserSameDuration(t *testing.T) {
	reader := timeoutReadCloser{
		reader: &testReader{
			duration: time.Millisecond,
			count:    5,
		},
		duration: time.Millisecond,
	}
	b := make([]byte, 100)
	_, err := reader.Read(b)
	if err != nil {
		t.Log(err)
	}
}

func TestWithResponseReadTimeout(t *testing.T) {
	r := Request{
		Retryer: NoOpRetryer{},
		HTTPResponse: &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(nil)),
		},
	}
	r.ApplyOptions(WithResponseReadTimeout(time.Second))
	err := r.Send()
	if err != nil {
		t.Error(err)
	}
	v, ok := r.HTTPResponse.Body.(*timeoutReadCloser)
	if !ok {
		t.Fatalf("Expected the body to be a timeoutReadCloser, got %T", r.HTTPResponse.Body)
	}
	if v.duration != time.Second {
		t.Errorf("Expected %v, but receive %v\n", time.Second, v.duration)
	}
}
