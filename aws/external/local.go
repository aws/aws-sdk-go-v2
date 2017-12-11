package external

import (
	"fmt"
	"net"
	"net/url"
)

var lookupHostFn = net.LookupHost

func isLoopbackHost(host string) bool {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.IsLoopback()
	}

	// Host is not an ip, perform lookup
	addrs, err := lookupHostFn(host)
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		if !net.ParseIP(addr).IsLoopback() {
			return false
		}
	}

	return true
}

func validateLocalURL(v string) error {
	u, err := url.Parse(v)
	if err != nil {
		return err
	}

	if host := u.Hostname(); len(host) == 0 || !isLoopbackHost(host) {
		return fmt.Errorf("invalid endpoint host, %q, only host resolving to loopback addresses are allowed", host)
	}

	return nil
}
