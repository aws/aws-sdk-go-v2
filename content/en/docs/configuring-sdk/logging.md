---
title: "Logging"
linkTitle: "Logging"
date: "2020-11-12"
description: Using Logging facilities for debugging SDK requests.
---

The {{% alias sdk-go %}} has logging facilities available that allow your application to enable debugging information
for debugging and diagnosing request issues or failures. The [Logger]({{< apiref smithy="logging#Logger" >}}) interface
and [ClientLogMode]({{< apiref "aws#ClientLogMode" >}}) are the main components available to you for determining how and
what should be logged by clients.

## Logger

When constructing an [Config]({{< apiref "aws#Config" >}}) using
[LoadDefaultConfig]({{< apiref "config#LoadDefaultConfig" >}}) a default `Logger` will is configured to send log
messages to process' standard error (stderr). A custom logger that satisfies the 
[Logger]({{< apiref smithy="logging#Logger" >}}) interface can be passed as an argument to `LoadDefaultConfig` 
by wrapping it with [config.WithLogger]({{< apiref "config#WithLogger" >}}).

For example to configure our clients to use our `applicationLogger`:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), 
	config.WithLogger(applicationLogger))
```

Now clients configured using the constructed `aws.Config` will send log messages to `applicationLogger`.

### Context-Aware Loggers

A Logger implementation may implement the optional [ContextLogger]({{< apiref smithy="logging#ContextLogger" >}})
interface. Loggers that implement this interface will have their `WithContext` methods invoked with the current context.
This allows your logging implementations to return a new `Logger` that can write additional logging metadata based
on values present in the context.

## ClientLogMode

By default, service clients do not produce log messages. To configure clients to send log messages for debugging
purposes, use the [ClientLogMode]({{< apiref "aws#ClientLogMode" >}}) member on `Config`. `ClientLogMode`
can be set to enable debugging messaging for:

* Signature Version 4 (SigV4) Signing
* Request Retries
* HTTP Requests
* HTTP Responses

For example to enable logging of HTTP requests and retries:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), 
	config.WithClientLogMode(aws.LogRetries | aws.LogRequest))
```

See [ClientLogMode]({{< apiref "aws#ClientLogMode" >}}) for the different client log modes available.

