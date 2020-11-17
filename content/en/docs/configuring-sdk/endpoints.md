---
title: "Configuring Client Endpoints"
linkTitle: "Endpoints"
date: "2020-11-12"
description: Customizing service client endpoints.
---

In most cases you use the endpoint that is pre-configured for a service. However, you can specify a custom endpoint,
such as for pre-release versions of the service.

A [EndpointResolver]({{< apiref "aws#EndpointResolver" >}}) can be configured to provide custom endpoint resolution
logic for service clients. A custom endpoint resolver has the ability to override the resolution of a service's endpoint
resolution for all endpoints, or a just specific regional endpoint. Custom endpoint resolver have the ability trigger
the service's endpoint resolution logic to fallback if a custom resolver does not wish to resolve a requested
endpoint. [EndpointResolverFunc]({{< apiref "aws#EndpointResolverFunc" >}}) can be used to easily wrap functions to
satisfy the `EndpointResolver` interface.

A `EndpointResolver` can be easily configured by passing the resolver wrapped with
[WithEndpointResolver]({{< apiref "config#WithEndpointResolver" >}}) to
[LoadDefaultConfig]({{< apiref "config#LoadDefaultConfig" >}}), allowing for the ability to override endpoints when
loading credentials, as well as configuring the resulting `aws.Config` with your custom endpoint resolver.

The endpoint resolver is given the service and region as a string, allowing for the resolver to dynamically drive its
behavior. Each service client package has an exported `ServiceID` constant which can be used to determine which service
client is invoking your endpoint resolver.

An endpoint resolver can use the [EndpointNotFoundError]({{< apiref "aws#EndpointNotFoundError" >}}) sentinel error
value to trigger fallback resolution to the service clients default resolution logic. This allows you to selectively
override one or more endpoints seamlessly without having to handle fallback logic.

If your endpoint resolver implementation returns an error other than `EndpointNotFoundError`, endpoint resolution will
stop and the service operation will return an error to your application.

## Examples

### Overriding Endpoint with Fallback

The following code snippet shows how a single service endpoint can be overridden for {{% alias service=DDB %}} with
fallback behavior for other endpoints:

```go
customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
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

cfg, err := config.LoadDefaultConfig(config.WithEndpointResolver(customResolver))
```

### Overriding Endpoint Without Fallback

The following code snippet shows how a single service endpoint can be overridden for {{% alias service=DDB %}} without
fallback behavior for other endpoints:

```go
customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
    if service == dynamodb.ServiceID && region == "us-west-2" {
        return aws.Endpoint{
            PartitionID:   "aws",
            URL:           "https://test.us-west-2.amazonaws.com",
            SigningRegion: "us-west-2",
        }, nil
    }
    return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
})

cfg, err := config.LoadDefaultConfig(config.WithEndpointResolver(customResolver))
```

### Immutable Endpoints

{{% pageinfo color="warning" %}}
Setting an endpoint as immutable may prevent some service client features from functioning correctly, and could result
in undefined behavior. Caution should be taken when defining an endpoint as immutable.
{{% /pageinfo %}}

Some service clients, such as {{% alias service=S3 %}}, can modify the endpoint returned by the resolver for certain
service operations. For example, the {{% alias service=S3 %}} will automatically handle
[Virtual Bucket Addressing](https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html) by mutating the
resolved endpoint. You can prevent the SDK from mutating your custom endpoints by setting 
[HostnameImmutable]({{< apiref "aws#Endpoint.HostnameImmutable" >}}) to`true`. For example:

```go
customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
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

cfg, err := config.LoadDefaultConfig(config.WithEndpointResolver(customResolver))
```

