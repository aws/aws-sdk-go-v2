---
title: "Configuring Client Endpoints"
linkTitle: "Endpoints"
date: "2020-11-12"
description: Customizing service client endpoints.
---

{{% pageinfo color="warning" %}}
Endpoint resolution is an advanced SDK topic. By changing these settings
you risk breaking your code. The default settings should be applicable to
most users in production environments.
{{% /pageinfo %}}

The {{% alias sdk-go %}} provides the ability to configure a custom
endpoint to be used for a service. In most cases, the default
configuration will suffice. Configuring custom endpoints allows for
additional behavior, such as working with pre-release versions of a
service.

## Customization

There are two "versions" of endpoint resolution config within the SDK.
* v2, released in Q3 of 2023, configured via:
  * `EndpointResolverV2`
  * `BaseEndpoint`
* v1, released alongside the SDK, configured via:
  * `EndpointResolver`

We recommend users of v1 endpoint resolution migrate to v2 to obtain access to
newer endpoint-related service features.

## V2: `EndpointResolverV2` + `BaseEndpoint`

In resolution v2, `EndpointResolverV2` is the definitive mechanism through
which endpoint resolution occurs. The resolver's `ResolveEndpoint` method is
invoked as part of the workflow for every request you make in the SDK. The
hostname of the `Endpoint` returned by the resolver is used **as-is** when
making the request (operation serializers can still append to the HTTP path,
however).

Resolution v2 includes an additional client-level config, `BaseEndpoint`, which
is used to specify a "base" hostname for the instance of your service. The
value set here is not definitive-- it is ultimately passed as a parameter to
the client's `EndpointResolverV2` when final resolution occurs (read on for
more information about `EndpointResolverV2` parameters). The resolver
implementation then has the opportunity to inspect and potentially modify that
value to determine the final endpoint.

For example, if you perform an S3 `GetObject` request against a given bucket
with a client where you've specified a `BaseEndpoint`, the default resolver
will inject the bucket into the hostname if it is virtual-host compatible
(assuming you haven't disabled virtual-hosting in client config).

In practice, `BaseEndpoint` will most likely be used to point your client at a
development or preview instance of a service.

### `EndpointResolverV2` parameters

Each service takes a specific set of inputs which are passed to its resolution
function, defined in each service package as `EndpointParameters`.

Every service includes the following base parameters, which are used to
facilitate general endpoint resolution within AWS:

| name           | type     | description                                                |
|----------------|----------|------------------------------------------------------------|
| `Region`       | `string` | The client's AWS region                                    |
| `Endpoint`     | `string` | The value set for `BaseEndpoint` in client config          |
| `UseFips`      | `bool`   | Whether FIPS endpoints are enabled in client config        |
| `UseDualStack` | `bool`   | Whether dual-stack endpoints are enabled in client config  |

Services can specify additional parameters required for resolution. For
example, S3's `EndpointParameters` include the bucket name, as well as several
S3-specific feature settings such as whether virtual host addressing is
enabled.

If you are implementing your own `EndpointResolverV2`, you should never need to
construct your own instance of `EndpointParameters`.  The SDK will source the
values per-request and pass them to your implementation.

### A note about Amazon S3
Amazon S3 is a complex service with many of its features modeled through
complex endpoint customizations, such as bucket virtual hosting, S3 MRAP, and
more.

Because of this, we recommend that you don't replace the `EndpointResolverV2`
implementation in your S3 client. If you need to extend its resolution
behavior, perhaps by sending requests to a local development stack with
additional endpoint considerations, we recommend wrapping the default
implementation such that it delegates back to the default as a fallback (shown
in examples below).

### Examples

#### With `BaseEndpoint`

The following code snippet shows how to point your S3 client at a local
instance of a service, which in this example is hosted on the loopback device
at port 8080.

```go
client := s3.NewFromConfig(cfg, func (o *svc.Options) {
    o.BaseEndpoint = aws.String("https://localhost:8080/")
})
```

#### With `EndpointResolverV2`

The following code snippet shows how to inject custom behavior into S3's
endpoint resolution using `EndpointResolverV2`.

```go
import (
    "context"
    "net/url"

    "github.com/aws/aws-sdk-go-v2/service/s3"
    smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct {
    // you could inject additional application context here as well
}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
        smithyendpoints.Endpoint, error,
    ) {
    if /* input params or caller context indicate we must route somewhere */ {
        u, err := url.Parse("https://custom.service.endpoint/")
        if err != nil {
            return smithyendpoints.Endpoint{}, err
        }
        return smithyEndpoints.Endpoint{
            URI: *u,
        }, nil
    }

    // delegate back to the default v2 resolver otherwise
    return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params), nil
}

func main() {
    // load config...

    client := s3.NewFromConfig(cfg, func (o *s3.Options) {
        o.EndpointResolverV2 = &resolverV2{
            // ...
        }
    })
}
```

#### With both

The following sample program demonstrates the interaction between
`BaseEndpoint` and `EndpointResolverV2`. **This is an advanced use case:**

```kotlin
import (
    "context"
    "fmt"
    "log"
    "net/url"

    "github.com/aws/aws-sdk-go-v2"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct {}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
        smithyendpoints.Endpoint, error,
    ) {
    // s3.Options.BaseEndpoint is accessible here:
    fmt.Printf("The endpoint provided in config is %s\n", *params.Endpoint)

    // fallback to default
    return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func main() {
    cfg, err := config.LoadDefaultConfig(context.Background()
    if (err != nil) {
        log.Fatal(err)
    }

    client := s3.NewFromConfig(cfg, func (o *s3.Options) {
        o.BaseEndpoint = aws.String("https://endpoint.dev/")
        o.EndpointResolverV2 = &resolverV2{}
    })

    // ignore the output, this is just for demonstration
    client.ListBuckets(context.Background(), nil)
}
```

When run, the above program outputs the following:

```
The endpoint provided in config is https://endpoint.dev/
```


## V1: `EndpointResolver`

{{% pageinfo color="warning" %}}
Endpoint resolution v1 is retained for backwards compatibility and is isolated
from the modern behavior in endpoint resolution v2. It will only be used if the
`EndpointResolver` field is set by the caller.

Use of v1 will most likely prevent you from accessing endpoint-related service
features introduced with or after the release of v2 resolution. See "Migration"
for instructions on how to upgrade.
{{% /pageinfo %}}

A [EndpointResolver]({{< apiref "aws#EndpointResolver" >}}) can be configured
to provide custom endpoint resolution logic for service clients. You can use a
custom endpoint resolver to override a service's endpoint resolution logic for
all endpoints, or a just specific regional endpoint. Custom endpoint resolver
can trigger the service's endpoint resolution logic to fallback if a custom
resolver does not wish to resolve a requested endpoint.
[EndpointResolverWithOptionsFunc]({{< apiref
"aws#EndpointResolverWithOptionsFunc" >}}) can be used to easily wrap functions
to satisfy the `EndpointResolverWithOptions` interface.

A `EndpointResolver` can be easily configured by passing the resolver wrapped
with [WithEndpointResolverWithOptions]({{< apiref
"config#WithEndpointResolverWithOptions" >}}) to [LoadDefaultConfig]({{< apiref
"config#LoadDefaultConfig" >}}), allowing for the ability to override endpoints
when loading credentials, as well as configuring the resulting `aws.Config`
with your custom endpoint resolver.

The endpoint resolver is given the service and region as a string, allowing for
the resolver to dynamically drive its behavior. Each service client package has
an exported `ServiceID` constant which can be used to determine which service
client is invoking your endpoint resolver.

An endpoint resolver can use the [EndpointNotFoundError]({{< apiref
"aws#EndpointNotFoundError" >}}) sentinel error value to trigger fallback
resolution to the service clients default resolution logic. This allows you to
selectively override one or more endpoints seamlessly without having to handle
fallback logic.

If your endpoint resolver implementation returns an error other than
`EndpointNotFoundError`, endpoint resolution will stop and the service
operation returns an error to your application.

### Examples

#### With fallback

The following code snippet shows how a single service endpoint can be
overridden for {{% alias service=DDB %}} with fallback behavior for other
endpoints:

```go
customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
    if service == dynamodb.ServiceID && region == "us-west-2" {
        return aws.Endpoint{
            PartitionID:   "aws",
            URL:           "https://test.us-west-2.amazonaws.com",
            SigningRegion: "us-west-2",
        }, nil
    }
    // returning EndpointNotFoundError will allow the service to fallback to it's default resolution
    return aws.Endpoint{}, &aws.EndpointNotFoundError{}
})

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
```

#### Without fallback

The following code snippet shows how a single service endpoint can be
overridden for {{% alias service=DDB %}} without fallback behavior for other
endpoints:

```go
customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
    if service == dynamodb.ServiceID && region == "us-west-2" {
        return aws.Endpoint{
            PartitionID:   "aws",
            URL:           "https://test.us-west-2.amazonaws.com",
            SigningRegion: "us-west-2",
        }, nil
    }
    return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
})

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
```

### Immutable endpoints

{{% pageinfo color="warning" %}}
Setting an endpoint as immutable may prevent some service client features from
functioning correctly, and could result in undefined behavior. Caution should
be taken when defining an endpoint as immutable.
{{% /pageinfo %}}

Some service clients, such as {{% alias service=S3 %}}, may modify the endpoint
returned by the resolver for certain service operations. For example, {{% alias
service=S3 %}} will automatically handle [Virtual Bucket
Addressing](https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html)
by mutating the resolved endpoint. You can prevent the SDK from mutating your
custom endpoints by setting [HostnameImmutable]({{< apiref
"aws#Endpoint.HostnameImmutable" >}}) to `true`. For example:

```go
customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
    if service == dynamodb.ServiceID && region == "us-west-2" {
        return aws.Endpoint{
            PartitionID:   "aws",
            URL:           "https://test.us-west-2.amazonaws.com",
            SigningRegion: "us-west-2",
            HostnameImmutable: true,
        }, nil
    }
    return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
})

cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
```

## Migration

When migrating from v1 to v2 of endpoint resolution, the following general principles apply:
* Returning an [Endpoint]({{< apiref "aws#Endpoint" >}}) with [HostnameImmutable]({{<
  apiref "aws#Endpoint.HostnameImmutable" >}}) set to `false` is roughly
  equivalent to setting `BaseEndpoint` to the originally returned URL from v1 and
  leaving `EndpointResolverV2` as the default.
* Returning an Endpoint with HostnameImmutable set to `true` is roughly
  equivalent to implementing an `EndpointResolverV2` which returns the
  originally returned URL from v1.

Examples for these cases are provided below.

{{% pageinfo color="warning" %}}
V1 immutable endpoints and V2 resolution are not equivalent in behavior. For
example, signing overrides for custom features like S3 Object Lambda would
still be set for immutable endpoints returned via v1 code, but the same will
not be done for v2.
{{% /pageinfo %}}

### Examples

#### Mutable endpoint

The following code sample demonstrates how to migrate a basic v1 endpoint
resolver that returns a modifiable endpoint:

```go
// v1
client := svc.NewFromConfig(cfg, func (o *svc.Options) {
    o.EndpointResolver = svc.EndpointResolverFromURL("https://custom.endpoint.api/")
})

// v2
client := svc.NewFromConfig(cfg, func (o *svc.Options) {
    // the value of BaseEndpoint is passed to the default EndpointResolverV2
    // implementation, which will handle routing for features such as S3 accelerate,
    // MRAP, etc.
    o.BaseEndpoint = aws.String("https://custom.endpoint.api/")
})
```

#### Immutable endpoint
```go
// v1
client := svc.NewFromConfig(cfg, func (o *svc.Options) {
    o.EndpointResolver = svc.EndpointResolverFromURL("https://custom.endpoint.api/", func (e *aws.Endpoint) {
        e.HostnameImmutable = true
    })
})

// v2
import (
    smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type staticResolver struct {}

func (*staticResolver) ResolveEndpoint(ctx context.Context, params svc.EndpointParameters) (
        smithyendpoints.Endpoint, error,
    ) {
    // This value will be used as-is when making the request.
    u, err := url.Parse("https://custom.endpoint.api/")
    if err != nil {
        return smithyendpoints.Endpoint{}, err
    }
    return smithyendpoints.Endpoint{
        URI: *u,
    }, nil
}

client := svc.NewFromConfig(cfg, func (o *svc.Options) {
    o.EndpointResolverV2 = &staticResolver{}
})
```

