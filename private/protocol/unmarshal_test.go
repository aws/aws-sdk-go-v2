package protocol_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/private/protocol"
)

type mockCloser struct {
	*strings.Reader
	Closed bool
}

func (m *mockCloser) Close() error {
	m.Closed = true
	return nil
}

func TestUnmarshalDrainBody(t *testing.T) {
	b := &mockCloser{Reader: strings.NewReader("example body")}
	r := &aws.Request{HTTPResponse: &http.Response{
		Body: b,
	}}

	protocol.UnmarshalDiscardBody(r)
	if r.Error != nil {
		t.Fatalf("expect no error, got %v", r.Error)
	}
	if l := b.Len(); l != 0 {
		t.Errorf("expect no body, have %v length", l)
	}
	if !b.Closed {
		t.Errorf("expect closed, was not")
	}
}

func TestUnmarshalDrainBodyNoBody(t *testing.T) {
	r := &aws.Request{HTTPResponse: &http.Response{}}

	protocol.UnmarshalDiscardBody(r)
	if r.Error != nil {
		t.Fatalf("expect no error, got %v", r.Error)
	}
}
