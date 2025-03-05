package kitchensinktest

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type mockHTTP struct {
	resps []*http.Response
	index int

	reqs []*http.Request
}

func (m *mockHTTP) Do(r *http.Request) (*http.Response, error) {
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
