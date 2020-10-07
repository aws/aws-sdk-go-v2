// +build integration,perftest

package downloader

import (
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func NewHTTPClient(cfg ClientConfig) aws.HTTPClient {
	return aws.NewBuildableHTTPClient().WithTransportOptions(func(tr *http.Transport) {
		*tr = http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   cfg.Timeouts.Connect,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:     90 * time.Second,

			DisableKeepAlives:     !cfg.KeepAlive,
			TLSHandshakeTimeout:   cfg.Timeouts.TLSHandshake,
			ExpectContinueTimeout: cfg.Timeouts.ExpectContinue,
			ResponseHeaderTimeout: cfg.Timeouts.ResponseHeader,

			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
		}
	})
}
