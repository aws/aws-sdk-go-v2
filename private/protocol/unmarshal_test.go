package protocol_test

import (
	"net/http"
	"strings"
	"testing"

	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/stretchr/testify/assert"
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
	r := &request.Request{HTTPResponse: &http.Response{
		Body: b,
	}}

	protocol.UnmarshalDiscardBody(r)
	assert.NoError(t, r.Error)
	assert.Equal(t, 0, b.Len())
	assert.True(t, b.Closed)
}

func TestUnmarshalDrainBodyNoBody(t *testing.T) {
	r := &request.Request{HTTPResponse: &http.Response{}}

	protocol.UnmarshalDiscardBody(r)
	assert.NoError(t, r.Error)
}
