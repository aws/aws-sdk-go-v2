---
title: "Configuring Authentication"
linkTitle: "Authentication"
date: "2023-12-01"
description: Customizing service client authentication.
---

The {{% alias sdk-go %}} provides the ability to configure a custom
auth scheme resolver to be used for a service. In most cases, the default
configuration will suffice, but configuring custom authentication allows for
additional behavior such as working with pre-release service features.

## Auth scheme

An auth scheme is defined as the workflow through which the SDK supplements an
operation request with a caller identity.

An auth scheme uses the following components:
* A unique ID which identifies the scheme
* An identity resolver, which returns a caller identity used in the signing process (e.g. your AWS credentials)
* A signer, which performs the actual injection of caller identity into the operation's transport request (e.g. the `Authorization` HTTP header)

The following auth schemes are used within the SDK:

| Name | Scheme ID | Identity resolver | Signer | Notes |
| ---  | --------- | ----------------- | ------ | ----- |
| [SigV4](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_aws-signing.html) | `aws.auth#sigv4` | [aws.CredentialsProvider](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Credentials) | [v4.HTTPSigner](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/signer/v4#Signer) | The current default for most AWS service operations. |
| SigV4A | `aws.auth#sigv4a` | aws.CredentialsProvider | n/a | SigV4A usage is limited at this time, the signer implementation is internal. |
| SigV4Express | `com.amazonaws.s3#sigv4express` | [s3.ExpressCredentialsProvider](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#ExpressCredentialsProvider) | v4.HTTPSigner | Used for [Express One Zone](https://aws.amazon.com/s3/storage-classes/express-one-zone/) requests. |
| HTTP Bearer | `smithy.api#httpBearerAuth` | [smithybearer.TokenProvider](https://pkg.go.dev/github.com/aws/smithy-go/auth/bearer#TokenProvider) | [smithybearer.Signer](https://pkg.go.dev/github.com/aws/smithy-go/auth/bearer#Signer) | Used by the codecatalyst service. |
| Anonymous | `smithy.api#noAuth` | n/a | n/a | No authentication - no identity is required, and the request is not signed or authenticated. |

In {{% alias sdk-go %}}, the identity and signer components of an auth scheme
are configured on a service's client Options. The SDK will automatically pick
up and use the values for these components for the scheme it selects when an
operation is called.

**If you want to override the identity or signing behavior for a particular
scheme, you can do so by modifying the appropriate concrete field on client
Options. Modifying the default AuthScheme itself is not necessary.**

**Note:** For backwards compatibility reasons, the SDK implicitly allows the
use of the anonymous auth scheme if no identity resolvers are configured.
This can be manually achieved by setting all identity resolvers on a client's
`Options` to `nil` (the sigv4 identity resolver can also be set to
`aws.AnonymousCredentials{}`).

## Auth scheme resolution

### `AuthSchemes`

Each service client options includes an `AuthSchemes` field, which by default
is populated with the list of auth schemes supported by that service.

### `AuthSchemeResolver`

Each service client options includes an `AuthSchemeResolver` field. This
interface, defined per-service, is the API called by the SDK to determine the
possible authentication options for each operation.

{{% pageinfo color="warning" %}}
The following services have unique or customized authentication behavior. We
recommend you delegate to the default implementation and wrap accordingly if
you require custom authentication behavior therein:

| Service | Notes |
| ------- | ----- |
| S3 | Conditional use of SigV4A and SigV4Express depending on operation input. |
| EventBridge | Conditional use of SigV4A depending on operation input. |
| Cognito | Certain operations are anonymous-only. |
| SSO | Certain operations are anonymous-only. |
| STS | Certain operations are anonymous-only. |
{{% /pageinfo %}}

### `AuthSchemeResolver` parameters

Each service takes a specific set of inputs which are passed to its resolution
function, defined in each service package as `AuthResolverParameters`.

The base resolver parameters are as follows:

| name           | type     | description |
| -------------- | -------- | ----------- |
| `Operation` | `string` | The name of the operation being invoked. |
| `Region` | `string` | The client's AWS region. Only present for services that use SigV4[A]. |

If you are implementing your own resolver, you should never need to construct
your own instance of its parameters. The SDK will source these values
per-request and pass them to your implementation.

### Auth scheme selection

Once the auth scheme resolver has been called to provide a list of
authentication options, the SDK checks that list and selects the first scheme
that satisfies the following conditions:

* A scheme with matching ID is present in the client's `AuthSchemes`
* The scheme's identity resolver exists (is non-`nil`) on the client's Options
  (checked via the scheme's `GetIdentityResolver` method, the mapping to the
  concrete identity resolver types described above is handled internally).

If the SDK encounters the anonymous option in the list, it is selected
automatically. **If you are overriding a service's auth scheme resolver and
need to include anonymous auth, the resolver SHOULD return it at the end of the
list of options.**
