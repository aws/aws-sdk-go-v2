package kitchensinktest

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type mockHTTP struct {
	resps []*http.Response
	index int

	err error

	reqs []*http.Request
}

func (m *mockHTTP) Do(r *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	m.reqs = append(m.reqs, r)
	resp := m.resps[m.index]
	m.index++
	return resp, nil
}

type mockCredentials struct {
	credentials []aws.Credentials
	index       int
}

func (m *mockCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	creds := m.credentials[m.index]
	m.index++
	return creds, nil
}

func mockResponseBody(v string) io.ReadCloser {
	return io.NopCloser(bytes.NewBuffer([]byte(v)))
}
