package aws

import (
	"io"
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
