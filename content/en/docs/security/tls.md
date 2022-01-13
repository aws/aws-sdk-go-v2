---
title: "TLS Version in AWS SDK for Go V2"
linkTitle: "TLS Version"
date: "2020-11-09"
---

The {{% alias sdk-go %}} configures the default HTTP client used by the SDK service clients to require a minimum version
of TLS 1.2 or greater. When using the [http.BuildableClient]({{< apiref "aws/transport/http#BuildableClient" >}}) to
customize the SDK HTTP client, the minimum TLS value is configured as TLS 1.2.

{{% pageinfo color="warning" %}}
If your application constructs an HTTP client using a method other than the provided `BuildableClient`, you must
configure your client to set the minimum TLS version to 1.2.
{{% /pageinfo %}}

## Enforcing a Minimum TLS Version

You can construct a custom an [http.Client]() or use the SDK provided
[http.BuildableClient]({{< apiref "aws/transport/http#BuildableClient" >}})
builder. The following example demonstrates how to specify a minimum TLS
version of [1.3](https://pkg.go.dev/crypto/tls#VersionTLS13) using the 
[http.BuildableClient]({{< apiref "aws/transport/http#BuildableClient" >}}).

```go
package main

import (
	"context"
	"crypto/tls"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// Create the custom HTTP client, configured for TLS 1.3 specified as the
	// minimum TLS version.
	httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.MinVersion = tls.VersionTLS13
	})

	// Load the SDK's configuration, and specify the custom HTTP client to be used
	// by all SDK API clients created from this configuration.
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(httpClient))

    // Use the loaded config and custom HTTP client to create SDK API client(s).
    // ...
}
```

