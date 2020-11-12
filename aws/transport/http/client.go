package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Defaults for the HTTPTransportBuilder.
var (
	// Default connection pool options
	DefaultHTTPTransportMaxIdleConns        = 100
	DefaultHTTPTransportMaxIdleConnsPerHost = 10

	// Default connection timeouts
	DefaultHTTPTransportIdleConnTimeout       = 90 * time.Second
	DefaultHTTPTransportTLSHandleshakeTimeout = 10 * time.Second
	DefaultHTTPTransportExpectContinueTimeout = 1 * time.Second

	// Default to TLS 1.2 for all HTTPS requests.
	DefaultHTTPTransportTLSMinVersion uint16 = tls.VersionTLS12
)

// Timeouts for net.Dialer's network connection.
var (
	DefaultDialConnectTimeout   = 30 * time.Second
	DefaultDialKeepAliveTimeout = 30 * time.Second
)

// BuildableClient provides a HTTPClient implementation with options to
// create copies of the HTTPClient when additional configuration is provided.
//
// The client's methods will not share the http.Transport value between copies
// of the BuildableClient. Only exported member values of the Transport and
// optional Dialer will be copied between copies of BuildableClient.
type BuildableClient struct {
	transport *http.Transport
	dialer    *net.Dialer

	initOnce sync.Once

	clientTimeout time.Duration
	client        *http.Client
}

// NewBuildableClient returns an initialized client for invoking HTTP
// requests.
func NewBuildableClient() *BuildableClient {
	return &BuildableClient{}
}

// Do implements the HTTPClient interface's Do method to invoke a HTTP request,
// and receive the response. Uses the BuildableClient's current
// configuration to invoke the http.Request.
//
// If connection pooling is enabled (aka HTTP KeepAlive) the client will only
// share pooled connections with its own instance. Copies of the
// BuildableClient will have their own connection pools.
//
// Redirect (3xx) responses will not be followed, the HTTP response received
// will returned instead.
func (b *BuildableClient) Do(req *http.Request) (*http.Response, error) {
	b.initOnce.Do(b.build)

	return b.client.Do(req)
}

func (b *BuildableClient) build() {
	b.client = wrapWithoutRedirect(&http.Client{
		Timeout:   b.clientTimeout,
		Transport: b.GetTransport(),
	})
}

func (b *BuildableClient) clone() *BuildableClient {
	cpy := NewBuildableClient()
	cpy.transport = b.GetTransport()
	cpy.dialer = b.GetDialer()
	cpy.clientTimeout = b.clientTimeout

	return cpy
}

// WithTransportOptions copies the BuildableClient and returns it with the
// http.Transport options applied.
//
// If a non (*http.Transport) was set as the round tripper, the round tripper
// will be replaced with a default Transport value before invoking the option
// functions.
func (b *BuildableClient) WithTransportOptions(opts ...func(*http.Transport)) *BuildableClient {
	cpy := b.clone()

	tr := cpy.GetTransport()
	for _, opt := range opts {
		opt(tr)
	}
	cpy.transport = tr

	return cpy
}

// WithDialerOptions copies the BuildableClient and returns it with the
// net.Dialer options applied. Will set the client's http.Transport DialContext
// member.
func (b *BuildableClient) WithDialerOptions(opts ...func(*net.Dialer)) *BuildableClient {
	cpy := b.clone()

	dialer := cpy.GetDialer()
	for _, opt := range opts {
		opt(dialer)
	}
	cpy.dialer = dialer

	tr := cpy.GetTransport()
	tr.DialContext = cpy.dialer.DialContext
	cpy.transport = tr

	return cpy
}

// WithTimeout Sets the timeout used by the client for all requests.
func (b *BuildableClient) WithTimeout(timeout time.Duration) *BuildableClient {
	cpy := b.clone()
	cpy.clientTimeout = timeout
	return cpy
}

// GetTransport returns a copy of the client's HTTP Transport.
func (b *BuildableClient) GetTransport() *http.Transport {
	var tr *http.Transport
	if b.transport != nil {
		tr = b.transport.Clone()
	} else {
		tr = defaultHTTPTransport()
	}

	return tr
}

// GetDialer returns a copy of the client's network dialer.
func (b *BuildableClient) GetDialer() *net.Dialer {
	var dialer *net.Dialer
	if b.dialer != nil {
		dialer = shallowCopyStruct(b.dialer).(*net.Dialer)
	} else {
		dialer = defaultDialer()
	}

	return dialer
}

// GetTimeout returns a copy of the client's timeout to cancel requests with.
func (b *BuildableClient) GetTimeout() time.Duration {
	return b.clientTimeout
}

func defaultDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   DefaultDialConnectTimeout,
		KeepAlive: DefaultDialKeepAliveTimeout,
		DualStack: true,
	}
}

func defaultHTTPTransport() *http.Transport {
	dialer := defaultDialer()

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   DefaultHTTPTransportTLSHandleshakeTimeout,
		MaxIdleConns:          DefaultHTTPTransportMaxIdleConns,
		MaxIdleConnsPerHost:   DefaultHTTPTransportMaxIdleConnsPerHost,
		IdleConnTimeout:       DefaultHTTPTransportIdleConnTimeout,
		ExpectContinueTimeout: DefaultHTTPTransportExpectContinueTimeout,
		ForceAttemptHTTP2:     true,
		TLSClientConfig: &tls.Config{
			MinVersion: DefaultHTTPTransportTLSMinVersion,
		},
	}

	return tr
}

// shallowCopyStruct creates a shallow copy of the passed in source struct, and
// returns that copy of the same struct type.
func shallowCopyStruct(src interface{}) interface{} {
	srcVal := reflect.ValueOf(src)
	srcValType := srcVal.Type()

	var returnAsPtr bool
	if srcValType.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
		srcValType = srcValType.Elem()
		returnAsPtr = true
	}
	dstVal := reflect.New(srcValType).Elem()

	for i := 0; i < srcValType.NumField(); i++ {
		ft := srcValType.Field(i)
		if len(ft.PkgPath) != 0 {
			// unexported fields have a PkgPath
			continue
		}

		dstVal.Field(i).Set(srcVal.Field(i))
	}

	if returnAsPtr {
		dstVal = dstVal.Addr()
	}

	return dstVal.Interface()
}

func wrapWithoutRedirect(c *http.Client) *http.Client {
	tr := c.Transport
	if tr == nil {
		tr = http.DefaultTransport
	}

	cc := *c
	cc.CheckRedirect = limitedRedirect
	cc.Transport = stubBadHTTPRedirectTransport{
		tr: tr,
	}

	return &cc
}

func limitedRedirect(r *http.Request, via []*http.Request) error {
	// Request.Response, in CheckRedirect is the response that is triggering
	// the redirect.
	resp := r.Response
	if r.URL.String() == stubBadHTTPRedirectLocation {
		resp.Header.Del(stubBadHTTPRedirectLocation)
		return http.ErrUseLastResponse
	}

	switch resp.StatusCode {
	case 307, 308:
		// Only allow 307 and 308 redirects as they preserve the method.
		return nil
	}

	return http.ErrUseLastResponse
}

type stubBadHTTPRedirectTransport struct {
	tr http.RoundTripper
}

const stubBadHTTPRedirectLocation = `https://amazonaws.com/badhttpredirectlocation`

func (t stubBadHTTPRedirectTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.tr.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	// TODO S3 is the only known service to return 301 without location header.
	// consider moving this to a S3 customization.
	switch resp.StatusCode {
	case 301, 302:
		if v := resp.Header.Get("Location"); len(v) == 0 {
			resp.Header.Set("Location", stubBadHTTPRedirectLocation)
		}
	}

	return resp, err
}
