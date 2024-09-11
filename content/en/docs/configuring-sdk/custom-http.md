---
title: "Customizing the HTTP Client"
linkTitle: "HTTP Client"
date: "2020-11-12"
description: "Create a custom HTTP client with the AWS SDK for Go V2 to specify custom timeout values."
---

The {{% alias sdk-go %}} uses a default HTTP client with default configuration values. Although you can change some of
these configuration values, the default HTTP client and transport are not sufficiently configured for customers using
the {{% alias sdk-go %}} in an environment with high throughput and low latency requirements. For more information, please refer to our [FAQ]({{% ref "faq.md" %}}) as configuration recommendations vary based on specific workloads.
This section describes how to configure a custom HTTP client, and use that client to create {{% alias sdk-go %}} calls.

To assist you in creating a custom HTTP client, this section describes how to the 
[NewBuildableClient]({{< apiref "aws/transport/http#NewBuildableClient" >}}) to configure custom settings, and use 
that client with an {{% alias sdk-go %}} service client.

Let's define what we want to customize.


## Overriding During Configuration Loading
Custom HTTP clients can be provided when calling [LoadDefaultConfig]({{< apiref "config#LoadDefaultConfig" >}}) by
wrapping the client using [WithHTTPClient]({{< apiref "config#WithHTTP" >}}) and passing the resulting value to 
`LoadDefaultConfig`. For example to pass `customClient` as our client:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(customClient))
```

## Timeout
The `BuildableHTTPClient` can be configured with a request timeout limit. This timeout includes the time to connect, 
process any redirects, and read the complete response body. For example to modify the client timeout:

```go
import "github.com/aws/aws-sdk-go-v2/aws/transport/http"

// ...

httpClient := http.NewBuildableClient().WithTimeout(time.Second*5)
```

## Dialer
The `BuildableHTTPClient` provides a builder mechanics for constructing clients with modified
[Dialer](https://golang.org/pkg/net/#Dialer) options. The following example shows how to configure a clients
`Dialer` settings.

```go
import awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
import "net"

// ...

httpClient := awshttp.NewBuildableClient().WithDialerOptions(func(d *net.Dialer) {
	d.KeepAlive = -1
	d.Timeout = time.Millisecond*500
})
```

### Settings

#### Dialer.KeepAlive

This setting represents the keep-alive period for an active network connection.

Set to a negative value to disable keep-alives.

Set to **0** to enable keep-alives if supported by the protocol and operating system.

Network protocols or operating systems that do not support keep-alives ignore this field. By default, TCP enables keep
alive.

See <https://golang.org/pkg/net/#Dialer.KeepAlive>

Set `KeepAlive` as **time.Duration**.

#### Dialer.Timeout

This setting represents the maximum amount of time a dial waits for a connection to be created.

Default is 30 seconds.

See <https://golang.org/pkg/net/#Dialer.Timeout>

Set `Timeout` as **time.Duration**.

## Transport

The `BuildableHTTPClient` provides a builder mechanics for constructing clients with modified 
[Transport](https://golang.org/pkg/net/http#Transport) options.

### Configuring a Proxy

If you cannot directly connect to the internet, you can use Go-supported
environment variables (`HTTP_PROXY` / `HTTPS_PROXY`) or create a custom HTTP client to
configure your proxy. The following example configures the client to use `PROXY_URL` as the proxy
endpoint:

```go
import awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
import "net/http"

// ...

httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
	proxyURL, err := url.Parse("PROXY_URL")
	if err != nil {
		log.Fatal(err)
	}
	tr.Proxy = http.ProxyURL(proxyURL)
})
```

### Other Settings

Below are a few other `Transport` settings that can be modified to tune the HTTP client. Any additional settings not
described here can be found in the [Transport](https://golang.org/pkg/net/http/#Transport) type documentation.
These settings can be applied as shown in the following example:

```go
import awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
import "net/http"

// ...

httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
	tr.ExpectContinueTimeout = 0
	tr.MaxIdleConns = 10
})
```

#### Transport.ExpectContinueTimeout

If the request has an "Expect: 100-continue" header, this setting represents the maximum amount of time to wait for a
server's first response headers after fully writing the request headers, This time does not include the time to send the
request header. The HTTP client sends its payload after this timeout is exhausted.

Default 1 second.

Set to **0** for no timeout and send request payload without waiting. One use case is when you run into issues with
proxies or third party services that take a session similar to the use of Amazon S3 in the function shown later.

See <https://golang.org/pkg/net/http/#Transport.ExpectContinueTimeout>

Set `ExpectContinue` as **time.Duration**.

#### Transport.IdleConnTimeout

This setting represents the maximum amount of time to keep an idle network connection alive between HTTP requests.

Set to **0** for no limit.

See <https://golang.org/pkg/net/http/#Transport.IdleConnTimeout>

Set `IdleConnTimeout` as **time.Duration**.

#### Transport.MaxIdleConns

This setting represents the maximum number of idle (keep-alive) connections across all hosts. One use case for
increasing this value is when you are seeing many connections in a short period from the same clients

**0** means no limit.

See <https://golang.org/pkg/net/http/#Transport.MaxIdleConns>

Set`MaxIdleConns` as **int**.

#### Transport.MaxIdleConnsPerHost

This setting represents the maximum number of idle (keep-alive) connections to keep per-host. One use case for
increasing this value is when you are seeing many connections in a short period from the same clients

Default is two idle connections per host.

Set to **0** to use DefaultMaxIdleConnsPerHost (2).

See <https://golang.org/pkg/net/http/#Transport.MaxIdleConnsPerHost>

Set `MaxIdleConnsPerHost` as **int**.

#### Transport.ResponseHeaderTimeout

This setting represents the maximum amount of time to wait for a client to read the response header.

If the client isn't able to read the response's header within this duration, the request fails with a timeout error.

Be careful setting this value when using long-running Lambda functions, as the operation does not return any response
headers until the Lambda function has finished or timed out. However, you can still use this option with the **
InvokeAsync** API operation.

Default is no timeout; wait forever.

See <https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout>

Set `ResponseHeaderTimeout` as **time.Duration**.

#### Transport.TLSHandshakeTimeout

This setting represents the maximum amount of time waiting for a TLS handshake to be completed.

Default is 10 seconds.

Zero means no timeout.

See <https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout>

Set `TLSHandshakeTimeout` as **time.Duration**.
