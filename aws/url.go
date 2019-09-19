package aws

import "net/url"

// URLHostname will extract the Hostname without port from the URL value.
//
// Wrapper of net/url#URL.Hostname for backwards Go version compatibility.
// Todo: Check if needed for +go1.11
func URLHostname(url *url.URL) string {
	return url.Hostname()
}
