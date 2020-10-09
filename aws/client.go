package aws

import (
	"net/http"
)

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
