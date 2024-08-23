---
title: "Retries and Timeouts"
linkTitle: "Retries and Timeouts"
---

The {{% alias sdk-go %}} enables you to configure the retry behavior of requests to HTTP services. By default,
service clients use [retry.Standard]({{< apiref "aws/retry#Standard" >}}) as their default retryer. If the default
configuration or behavior does not meet your application requirements, you can adjust the retryer configuration or
provide your own retryer implementation.

The {{% alias sdk-go %}} provides a [aws.Retryer]({{< apiref "aws#Retryer" >}}) interface that defines the set of
methods required by a retry implementation to implement. The SDK provides two implementations for retries:
[retry.Standard]({{< apiref "aws/retry#Standard" >}}) and [aws.NoOpRetryer]({{< apiref "aws#NoOpRetryer" >}}).

## Standard Retryer

The [retry.Standard]({{< apiref "aws/retry#Standard" >}})  retryer is the default `aws.Retryer` implementation used
by SDK clients. The standard retryer is a rate limited retryer with a configurable number of max attempts, and the
ability to tune the request back off policy.

The following table defines the default values for this retryer:
Property | Default
 --- | ---
Max Number of Attempts | 3
Max Back Off Delay | 20 seconds

When a retryable error occurs while invoking your request, the standard retryer will use its provided configuration
to delay and subsequently retry the request. Retries add to the overall latency of your request, and you must configure
retryer if the default configuration does not meet your application requirements.

See the [retry]({{< apiref "aws/retry" >}}) package documentation for details on what errors are considered as
retryable by the standard retryer implementation.

## NopRetryer

The [aws.NopRetryer]({{< apiref "aws#NopRetryer" >}}) is a `aws.Retryer` implementation that is provided if you wish
to disable all retry attempts. When invoking a service client operation, this retryer will only allow the request to
be attempted once, and any resulting error will be returned to the calling application.

## Customizing Behavior

The SDK provides a set of helper utilities that wrap an `aws.Retryer` implementation, and returns the provided retryer
wrapped with the desired retry behavior. You can override the default retryer for all clients, per client, or per
operation depending on your applications requirements. To see additional examples showing how to do this see the
[retry]({{< apiref "aws/retry" >}}) package documentation examples.

{{% alert color="warning" %}}
If specifying a global `aws.Retryer` implementation using `config.WithRetryer`, you must ensure that you return a new
instance of the `aws.Retryer` each invocation. This will ensure that you won't create a global retry token bucket across
all service clients.
{{% /alert %}}

### Limiting the max number of attempts

You use [retry.AddWithMaxAttempts]({{< apiref "aws/retry#AddWithMaxAttempts" >}}) to wrap an `aws.Retryer`
implementation to set the max number attempts to your desired value. Setting max attempts to zero will allow the SDK
to retry all retryable errors until the request succeeds, or a non-retryable error is returned.

For example, you can the following code to wrap the standard client retryer with a maximum of five attempts:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/aws/retry"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
	return retry.AddWithMaxAttempts(retry.NewStandard(), 5)
}))
if err != nil {
	return err
}

client := s3.NewFromConfig(cfg)
```

### Limiting the max back off delay

You use [retry.AddWithMaxBackoffDelay]({{< apiref "aws/retry#AddWithMaxBackoffDelay" >}}) to wrap an `aws.Retryer`
implementation and limit the max back off delay that is allowed to occur between retrying a failed request.

For example, you can the following code to wrap the standard client retryer with a desired max delay of five seconds:

```go
import "context"
import "time"
import "github.com/aws/aws-sdk-go-v2/aws/retry"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
	return retry.AddWithMaxBackoffDelay(retry.NewStandard(), time.Second*5)
}))
if err != nil {
	return err
}

client := s3.NewFromConfig(cfg)
```

### Retry additional API error codes

You use [retry.AddWithErrorCodes]({{< apiref "aws/retry#AddWithErrorCodes" >}}) to wrap an `aws.Retryer`
implementation and include additional API error codes that should be considered retryable.

For example, you can the following code to wrap the standard client retryer to include the {{% alias service="S3" %}}
`NoSuchBucketException` exception as retryable.

```go
import "context"
import "time"
import "github.com/aws/aws-sdk-go-v2/aws/retry"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/aws-sdk-go-v2/service/s3/types"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRetryer(func() aws.Retryer {
	return retry.AddWithErrorCodes(retry.NewStandard(), (*types.NoSuchBucketException)(nil).ErrorCode())
}))
if err != nil {
	return err
}

client := s3.NewFromConfig(cfg)
```

### Client-side rate limiting

The {{% alias sdk-go %}} introduces a new client-side rate-limiting mechanism
in the standard retry policy to align with the behavior of modern SDKs. This is
behavior is controlled by the [RateLimiter]({{< apiref "aws/retry#RateLimiter" >}})
field on a retryer's [options]({{< apiref "aws/retry#StandardOptions" >}}).

A RateLimiter operates as a token bucket with a set capacity, where operation
attempt failures consume tokens. A retry that attempts to consume more tokens
than what's available results in operation failure with a
[QuotaExceededError]({{< apiref "aws/ratelimit#QuotaExceededError" >}}).

The default implementation is parameterized as follows (how to modify each setting):
- a capacity of 500 (set the value of RateLimiter on StandardOptions using [NewTokenRateLimit]({{< apiref "aws/ratelimit#NewTokenRateLimit" >}}))
- a retry caused by a timeout costs 10 tokens (set RetryTimeoutCost on StandardOptions)
- a retry caused by other errors costs 5 tokens (set RetryCost on StandardOptions)
- an operation that succeeds on the 1st attempt adds 1 token (set NoRetryIncrement on StandardOptions)
  - operations that succeed on the 2nd or later attempt do not add back any tokens

If you find that the default behavior does not fit your application's needs,
you can disable it with [ratelimit.None]({{< apiref "aws/ratelimit#pkg-variables" >}}).

#### Example: modified rate limiter

```go
import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/aws/ratelimit"
    "github.com/aws/aws-sdk-go-v2/aws/retry"
    "github.com/aws/aws-sdk-go-v2/config"
)

// ...

cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRetryer(func() aws.Retryer {
    return retry.NewStandard(func(o *retry.StandardOptions) {
        // Makes the rate limiter more permissive in general. These values are
        // arbitrary for demonstration and may not suit your specific
        // application's needs.
        o.RateLimiter = ratelimit.NewTokenRateLimit(1000)
        o.RetryCost = 1
        o.RetryTimeoutCost = 3
        o.NoRetryIncrement = 10
    })
}))
```

#### Example: no rate limit using ratelimit.None

```go
import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/aws/ratelimit"
    "github.com/aws/aws-sdk-go-v2/aws/retry"
    "github.com/aws/aws-sdk-go-v2/config"
)

// ...

cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRetryer(func() aws.Retryer {
    return retry.NewStandard(func(o *retry.StandardOptions) {
        o.RateLimiter = ratelimit.None
    })
}))
```

##  Timeouts

You use the [context](https://golang.org/pkg/context/) package to set timeouts or deadlines when invoking a service
client operation. Use the [context.WithDeadline](https://golang.org/pkg/context/#WithDeadline) to wrap your applications
context and set a deadline to a specific time by which the invoked operation must be completed. To set a timeout
after a certain `time.Duration` use [context.WithTimeout](https://golang.org/pkg/context/#WithTimeout). The SDK passes
the provided `context.Context` to the HTTP transport client when invoking a service API. If the context passed to the SDK
is cancelled or becomes cancelled while invoking the operation, the SDK will not retry the request further and will
return to the calling application. You must handle context cancellation appropriately in your application in cases where
the context provided to the SDK has become cancelled.

### Setting a timeout

The following example shows how to set a timeout for a service client operation.

```go
import "context"
import "time"

// ...

ctx := context.TODO() // or appropriate context.Context value for your application

client := s3.NewFromConfig(cfg)

// create a new context from the previous ctx with a timeout, e.g. 5 seconds
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

resp, err := client.GetObject(ctx, &s3.GetObjectInput{
	// input parameters
})
if err != nil {
	// handle error
}
```
