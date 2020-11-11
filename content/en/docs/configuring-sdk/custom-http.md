---
title: "Creating a Custom HTTP Client"
linkTitle: "Custom HTTP Client"
description: "Create a custom HTTP client with the AWS SDK for Go V2 to specify custom timeout values."
---

The AWS SDK for Go V2 uses a default HTTP client with default configuration values. Although you can change some of
these configuration values, the default HTTP client and transport are not sufficiently configurable for customers using
the AWS SDK for Go V2 in an environment with high throughput and low latency requirements. This section describes how to create a
custom HTTP client, and use that client to create AWS SDK for Go V2 calls.

To assist you in creating a custom HTTP client, this section describes how to create a structure to encapsulate the
custom settings, create a function to create a custom HTTP client based on those settings, and use that custom HTTP
client to call an AWS SDK for Go V2 service client.

Let's define what we want to customize.

## Dialer.KeepAlive

This setting represents the keep-alive period for an active network connection.

Set to a negative value to disable keep-alives.

Set to **0** to enable keep-alives if supported by the protocol and operating system.

Network protocols or operating systems that do not support keep-alives ignore this field. By default, TCP enables keep
alive.

See <https://golang.org/pkg/net/#Dialer.KeepAlive>

We'll call this ``ConnKeepAlive`` as **time.Duration**.

## Dialer.Timeout

This setting represents the maximum amount of time a dial to wait for a connection to be created.

Default is 30 seconds.

See <https://golang.org/pkg/net/#Dialer.Timeout>

We'll call this ``Connect`` as **time.Duration**.

## Transport.ExpectContinueTimeout

This setting represents the maximum amount of time to wait for a server's first response headers after fully writing the
request headers, if the request has an "Expect: 100-continue" header. This time does not include the time to send the
request header. The HTTP client sends its payload after this timeout is exhausted.

Default 1 second.

Set to **0** for no timeout and send request payload without waiting. One use case is when you run into issues with
proxies or third party services that take a session similar to the use of Amazon S3 in the function shown later.

See <https://golang.org/pkg/net/http/#Transport.ExpectContinueTimeout>

We'll call this ``ExpectContinue`` as **time.Duration**.

## Transport.IdleConnTimeout

This setting represents the maximum amount of time to keep an idle network connection alive between HTTP requests.

Set to **0** for no limit.

See <https://golang.org/pkg/net/http/#Transport.IdleConnTimeout>

We'll call this ``IdleConn`` as **time.Duration**.

## Transport.MaxIdleConns

This setting represents the maximum number of idle (keep-alive) connections across all hosts. One use case for
increasing this value is when you are seeing many connections in a short period from the same clients

**0** means no limit.

See <https://golang.org/pkg/net/http/#Transport.MaxIdleConns>

We'll call this ``MaxAllIdleConns`` as **int**.

## Transport.MaxIdleConnsPerHost

This setting represents the maximum number of idle (keep-alive) connections to keep per-host. One use case for
increasing this value is when you are seeing many connections in a short period from the same clients

Default is two idle connections per host.

Set to **0** to use DefaultMaxIdleConnsPerHost (2).

See <https://golang.org/pkg/net/http/#Transport.MaxIdleConnsPerHost>

We'll call this ``MaxHostIdleConns`` as **int**.

## Transport.ResponseHeaderTimeout

This setting represents the maximum amount of time to wait for a client to read the response header.

If the client isn't able to read the response's header within this duration, the request fails with a timeout error.

Be careful setting this value when using long-running Lambda functions, as the operation does not return any response
headers until the Lambda function has finished or timed out. However, you can still use this option with the **
InvokeAsync** API operation.

Default is no timeout; wait forever.

See <https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout>

We'll call this ``ResponseHeader`` as **time.Duration**.

## Transport.TLSHandshakeTimeout

This setting represents the maximum amount of time waiting for a TLS handshake to be completed.

Default is 10 seconds.

Zero means no timeout.

See <https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout>

We'll call this ``TLSHandshake`` as **time.Duration**.

## Examples

See
a [complete example](https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/s3/CustomClient/CustomHttpClient.go)
on GitHub.
