---
title: "Migrating to the AWS SDK for Go V2"
linkTitle: "Migrating"
description: "How to migrate to the AWS SDK for Go V2 from AWS SDK for Go V1."
weight: 4
---

## Minimum Go Version

The {{% alias sdk-go %}} requires a minimum version of Go 1.15. Migration from AWS SDK for Go to {{% alias sdk-go %}}
might require you to upgrade your application by one or more Go versions. The latest version of Go can be downloaded on
the [Downloads](https://golang.org/dl/) page. See the [Release History](https://golang.org/doc/devel/release.html) for
more information about each Go version release, and relevant information required for upgrading.

## Modularization

The {{% alias sdk-go %}} has been updated to take advantage of the Go modules which became the default development mode
in Go 1.13. A number of packages provided by the SDK have been modularized and are independently versioned and
released respectively. This change enables improved application dependency modeling, and enables the SDK to
provide new features and functionality that follows the Go module versioning strategy.

The following list are some Go modules provided by the SDK:
Module | Description
--- | ---
`github.com/aws/aws-sdk-go-v2` | The SDK core
`github.com/aws/aws-sdk-go-v2/config` | Shared Configuration Loading
`github.com/aws/aws-sdk-go-v2/credentials` | AWS Credential Providers
`github.com/aws/aws-sdk-go-v2/feature/ec2/imds` | {{% alias service=EC2 %}} Instance Metadata Service Client

The SDK's service clients and higher level utilities modules are nested under the following import paths:
Import Root | Description
--- | ---
`github.com/aws/aws-sdk-go-v2/service/` | Service Client Modules
`github.com/aws/aws-sdk-go-v2/feature/` | High-Level utilities for services, for example the {{% alias service=S3 %}} Transfer Manager

## Configuration Loading

The [session]({{< apiref v1="aws/session" >}}) package and associated functionality are replaced with a simplified
configuration system provided by the [config]({{< apiref "config" >}}) package. The `config` package is a separate Go
module, and can be included in your applications dependencies by with `go get`.
```sh
go get github.com/aws/aws-sdk-go-v2/config
```

The [session.New]({{< apiref v1="aws/session#New" >}}), [session.NewSession]({{< apiref v1="aws/session#NewSession" >}}),
[NewSessionWithOptions]({{< apiref v1="aws/session#NewSessionWithOptions" >}}), and
[session.Must]({{< apiref v1="aws/session#Must" >}}) must be migrated to
[config.LoadDefaultConfig]({{< apiref "config#LoadDefaultConfig" >}}).

The `config` package provides several helper functions that aid in overriding the shared configuration loading
programmatically. These function names are prefixed with `With` followed by option that they override. Let's look at
some examples of how to migrate usage of the `session` package.

For more information on loading shared configuration
see [Configuring the {{% alias sdk-go %}}]({{% ref "configuring-sdk" %}}).

#### Examples
#####  Migrating from NewSession to LoadDefaultConfig
The following example shows how usage of `session.NewSession` without additional argument parameters is migrated to
`config.LoadDefaultConfig`.

```go
// V1 using NewSession

import "github.com/aws/aws-sdk-go/aws/session"

// ...

sess, err := session.NewSession()
if err != nil {
	// handle error
}
```

```go
// V2 using LoadDefaultConfig

import "context"
import "github.com/aws/aws-sdk-go-v2/config"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	// handle error
}
```

##### Migrating from NewSession with aws.Config options
The example shows how to migrate overriding of `aws.Config` values during configuration loading. One or more
`config.With*` helper functions can be provided to `config.LoadDefaultConfig` to override the loaded configuration
values. In this example the AWS Region is overridden to `us-west-2` using the
[config.WithRegion]({{< apiref "config#WithRegion" >}}) helper function.

```go
// V1

import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/session"

// ...

sess, err := session.NewSession(aws.Config{
	Region: aws.String("us-west-2")
})
if err != nil {
	// handle error
}
```

```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/config"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO(),
	config.WithRegion("us-west-2"),
)
if err != nil {
	// handle error
}
```

#####  Migrating from NewSessionWithOptions
This example shows how to migrate overriding values during configuration loading. Zero or more
`config.With*` helper functions can be provided to `config.LoadDefaultConfig` to override the loaded configuration
values. In this example we show how to override the target profile that is used when loading the AWS SDK shared
configuration.

```go
// V1

import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/session"

// ...

sess, err := session.NewSessionWithOptions(aws.Config{
	Profile: "my-application-profile"
})
if err != nil {
	// handle error
}
```

```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/config"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO(),
	config.WithSharedConfigProfile("my-application-profile"),
)
if err != nil {
	// handle error
}
```

## Credentials & Credential Providers

The [aws/credentials]({{< apiref v1="aws/credentials" >}}) package and associated credential providers have been
relocated to the [credentials]({{< apiref "credentials" >}}) package location. The `credentials` package is a Go module that
you retrieve by using `go get`.

```sh
go get github.com/aws/aws-sdk-go-v2/credentials
```

The {{% alias sdk-go %}} release updates the AWS Credential Providers to provide a consistent interface for retrieving
AWS Credentials. Each provider implements the [aws.CredentialsProvider]({{< apiref "aws#CredentialsProvider" >}})
interface, which defines a `Retrieve` method that returns a `(aws.Credentials, error)`.
[aws.Credentials]({{< apiref "aws#Credentials" >}}) that is analogous to the AWS SDK for Go
[credentials.Value]({{< apiref v1="aws/credentials#Value" >}}) type.

You must wrap `aws.CredentialsProvider` objects with [aws.CredentialsCache]({{< apiref "aws#CredentialsCache" >}}) to
allow credential caching to occur. You use [NewCredentialsCache]({{< apiref "aws#NewCredentialsCache" >}}) to
construct a `aws.CredentialsCache` object. By default, credentials configured by `config.LoadDefaultConfig` are wrapped
with `aws.CredentialsCache`.

The following table list the location changes of the AWS credential providers from AWS SDK for Go V1 to
{{% alias sdk-go %}}.

Name | V1 Import | V2 Import
--- | --- | ---
{{% alias service=EC2 %}} IAM Role Credentials | `github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds` | `github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds`
Endpoint Credentials | `github.com/aws/aws-sdk-go/aws/credentials/endpointcreds` | `github.com/aws/aws-sdk-go-v2/credentials/endpointcreds`
Process Credentials | `github.com/aws/aws-sdk-go/aws/credentials/processcreds` | `github.com/aws/aws-sdk-go-v2/credentials/processcreds`
{{% alias service=STSlong %}} | `github.com/aws/aws-sdk-go/aws/credentials/stscreds` | `github.com/aws/aws-sdk-go-v2/credentials/stscreds`

### Static Credentials

Applications that use [credentials.NewStaticCredentials]({{< apiref v1="aws/credentials#NewStaticCredentials" >}}) to
construct static credential programmatically must use
[credentials.NewStaticCredentialsProvider]({{< apiref "credentials#NewStaticCredentialsProvider" >}}).

#### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/credentials"

// ...

appCreds := credentials.NewStaticCredentials(accessKey, secretKey, sessionToken)
value, err := appCreds.Get()
if err != nil {
	// handle error
}
```

```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials"

// ...

appCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken))
value, err := appCreds.Retrieve(context.TODO())
if err != nil {
	// handle error
}
```

### {{% alias service=EC2 %}} IAM Role Credentials

You must migrate usage of [NewCredentials]({{< apiref v1="aws/credentials/ec2rolecreds#NewCredentials" >}}) and
[NewCredentialsWithClient]({{< apiref v1="aws/credentials/ec2rolecreds#NewCredentialsWithClient" >}}) to use
[New]({{< apiref "credentials/ec2rolecreds#New" >}}).

The `ec2rolecreds` package's `ec2rolecreds.New` takes functional options of
[ec2rolecreds.Options]({{< apiref "credentials/ec2rolecreds#Options" >}}) as
input, allowing you override the specific {{% alias service=EC2 %}} Instance
Metadata Service client to use, or to override the credential expiry window.

#### Example

```go
// V1

import "github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"

// ...

appCreds := ec2rolecreds.NewCredentials(sess)
value, err := appCreds.Get()
if err != nil {
	// handle error
}
```

```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"

// ...

// New returns an object of a type that satisfies the aws.CredentialProvider interface
appCreds := aws.NewCredentialsCache(ec2rolecreds.New())
value, err := appCreds.Retrieve(context.TODO())
if err != nil {
	// handle error
}
```

### Endpoint Credentials

You must migrate usage of [NewCredentialsClient]({{< apiref v1="aws/credentials/endpointcreds#NewCredentialsClient" >}})
and [NewProviderClient]({{< apiref v1="aws/credentials/endpointcreds#NewProviderClient" >}}) to use
[New]({{< apiref "credentials/endpointcreds#New" >}}).

The `endpointcreds` package's `New` function takes a string argument containing the
URL of an HTTP or HTTPS endpoint to retrieve credentials from, and functional
options of [endpointcreds.Options]({{< apiref
"credentials/endpointcreds#Options" >}}) to mutate the credentials provider and
override specific configuration settings.

### Process Credentials

You must migrate usage of [NewCredentials]({{< apiref v1="aws/credentials/processcreds#NewCredentials" >}}),
[NewCredentialsCommand]({{< apiref v1="aws/credentials/processcreds#NewCredentialsCommand" >}}), and
[NewCredentialsTimeout]({{< apiref v1="aws/credentials/processcreds#NewCredentialsTimeout" >}}) to use
[NewProvider]({{< apiref "credentials/processcreds#New" >}}) or
[NewProviderCommand]({{< apiref "credentials/processcreds#NewProviderCommand" >}}).

The `processcreds` package's `NewProvider` function takes a string argument that is the
command to be executed in the host environment's shell, and functional options
of [Options]({{< apiref "credentials/processcreds#Options" >}}) to mutate the
credentials provider and override specific configuration settings.

`NewProviderCommand` takes an implementation of the
[NewCommandBuilder]({{< apiref "credentials/processcreds#NewCommandBuilder" >}}) interface that defines
more complex process commands that might take one or more command-line arguments, or have certain execution environment
requirements. [DefaultNewCommandBuilder]({{< apiref "credentials/processcreds#DefaultNewCommandBuilder" >}}) implements
this interface, and defines a command builder for a process that requires multiple command-line arguments.

#### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/credentials/processcreds"

// ...

appCreds := processcreds.NewCredentials("/path/to/command")
value, err := appCreds.Get()
if err != nil {
	// handle error
}
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials/processcreds"

// ...

appCreds := aws.NewCredentialsCache(processcreds.NewProvider("/path/to/command"))
value, err := appCreds.Retrieve(context.TODO())
if err != nil {
	// handle error
}
```

### {{% alias service=STSlong %}} Credentials

#### AssumeRole

You must migrate usage of [NewCredentials]({{< apiref v1="aws/credentials/stscreds#NewCredentials" >}}), and
[NewCredentialsWithClient]({{< apiref v1="aws/credentials/stscreds#NewCredentialsWithClient" >}}) to
use [NewAssumeRoleProvider]({{< apiref "credentials/stscreds#NewAssumeRoleProvider" >}}).

The `stscreds` package's `NewAssumeRoleProvider` function must be called with
a [sts.Client]({{< apiref "service/sts#Client" >}}), and the {{% alias
service=IAMlong %}} Role ARN to be assumed from the provided `sts.Client`'s
configured credentials. You can also provide a set of functional options of
[AssumeRoleOptions]({{< apiref "credentials/stscreds#AssumeRoleOptions" >}}) to
modify other optional settings of the provider.

##### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/credentials/stscreds"

// ...

appCreds := stscreds.NewCredentials(sess, "arn:aws:iam::123456789012:role/demo")
value, err := appCreds.Get()
if err != nil {
	// handle error
}
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/credentials/stscreds"
import "github.com/aws/aws-sdk-go-v2/service/sts"

// ...

client := sts.NewFromConfig(cfg)

appCreds := stscreds.NewAssumeRoleProvider(client, "arn:aws:iam::123456789012:role/demo")
value, err := appCreds.Retrieve(context.TODO())
if err != nil {
	// handle error
}
```

#### AssumeRoleWithWebIdentity

You must migrate usage of
[NewWebIdentityCredentials]({{< apiref v1="aws/credentials/stscreds#NewWebIdentityCredentials" >}}),
[NewWebIdentityRoleProvider]({{< apiref v1="aws/credentials/stscreds#NewWebIdentityRoleProvider" >}}), and
[NewWebIdentityRoleProviderWithToken]({{< apiref v1="aws/credentials/stscreds#NewWebIdentityRoleProviderWithToken" >}})
to use [NewWebIdentityRoleProvider]({{< apiref "credentials/stscreds#NewWebIdentityRoleProvider" >}}).

The `stscreds` package's `NewWebIdentityRoleProvider` function must be called with a [sts.Client]({{< apiref "service/sts#Client" >}}), and the
{{% alias service=IAMlong %}} Role ARN to be assumed using the provided `sts.Client`'s configured credentials, and an
implementation of a [IdentityTokenRetriever]({{< apiref "credentials/stscreds#IdentityTokenRetriever" >}}) for
providing the OAuth 2.0 or OpenID Connect ID token.
[IdentityTokenFile]({{< apiref "credentials/stscreds#IdentityTokenFile" >}}) is an `IdentityTokenRetriever` that can
be used to provide the web identity token from a file located on the application's host file-system.
You can also provide a set of functional options of
[WebIdentityRoleOptions]({{< apiref "credentials/stscreds#WebIdentityRoleOptions" >}}) to modify other optional
settings for the provider.

##### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/credentials/stscreds"

// ...

appCreds := stscreds.NewWebIdentityRoleProvider(sess, "arn:aws:iam::123456789012:role/demo", "sessionName", "/path/to/token")
value, err := appCreds.Get()
if err != nil {
	// handle error
}
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials/stscreds"
import "github.com/aws/aws-sdk-go-v2/service/sts"

// ...

client := sts.NewFromConfig(cfg)

appCreds := aws.NewCredentialsCache(stscreds.NewWebIdentityRoleProvider(
		client,
		"arn:aws:iam::123456789012:role/demo",
		stscreds.IdentityTokenFile("/path/to/file"),
		func(o *stscreds.WebIdentityRoleOptions) {
			o.RoleSessionName = "sessionName"
		}))
value, err := appCreds.Retrieve(context.TODO())
if err != nil {
	// handle error
}
```

## Service Clients

{{% alias sdk-go %}} provides service client modules nested under the `github.com/aws/aws-sdk-go-v2/service` import path.
Each service client is contained in a Go package using each service's unique identifier. The following table provides
some examples of service import paths in the {{% alias sdk-go %}}.

Service Name | V1 Import Path | V2 Import Path
--- | --- | ---
{{% alias service=S3 %}} | `github.com/aws/aws-sdk-go/service/s3` | `github.com/aws/aws-sdk-go-v2/service/s3`
{{% alias service=DDBlong %}} | `github.com/aws/aws-sdk-go/service/dynamodb` | `github.com/aws/aws-sdk-go-v2/service/dynamodb`
{{% alias service=CWLlong %}} | `github.com/aws/aws-sdk-go/service/cloudwatchlogs` | `github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs`

Each service client package is an independently versioned Go module. To add the service client as a dependency of your
application, use the `go get` command with the service's import path. For example, to add the {{% alias service=S3 %}}
client to your dependencies use

```sh
go get github.com/aws/aws-sdk-go-v2/service/s3
```

### Client Construction

You can construct clients in the {{% alias sdk-go %}} using either the `New` or `NewFromConfig` constructor functions in the client's package.
When migrating from the AWS SDK for Go we recommend that you use the `NewFromConfig` variant, which will return a new
service client using values from an `aws.Config`. The `aws.Config` value will have been created while loading the SDK shared configuration using
`config.LoadDefaultConfig`. For details on creating service clients see
[Using AWS Services]({{% ref "making-requests" %}}).

##### Example 1


```go
// V1

import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/service/s3"

// ...

sess, err := session.NewSession()
if err != nil {
	// handle error
}

client := s3.New(sess)
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	// handle error
}

client := s3.NewFromConfig(cfg)
```


##### Example 2: Overriding Client Settings


```go
// V1

import "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/aws/session"
import "github.com/aws/aws-sdk-go/service/s3"

// ...

sess, err := session.NewSession()
if err != nil {
	// handle error
}

client := s3.New(sess, &aws.Config{
	Region: aws.String("us-west-2"),
})
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	// handle error
}

client := s3.NewFromConfig(cfg, func(o *s3.Options) {
	o.Region = "us-west-2"
})
```

### Endpoints

The [endpoints]({{< apiref v1="aws/endpoints" >}}) package no longer exists in the {{% alias sdk-go %}}. Each service
client now embeds its required AWS endpoint metadata within the client package. This reduces the overall binary size of
compiled applications by no longer including endpoint metadata for services not used by your application.

By default, service clients use their configured AWS Region to resolve the service endpoint for the target Region. If
your application requires a custom endpoint to be specified for a particular service and region, you can specify
a custom [aws.EndpointResolver]({{< apiref "aws#EndpointResolver" >}}) using the `EndpointResolver` field on the
`aws.Config` structure. If your application implements a custom
[endpoints.Resolver]({{< apiref v1="aws/endpoints#Resolver" >}}) you must migrate it to conform to the
`aws.EndpointResolver` interface. [aws.EndpointResolverFunc]({{< apiref "aws#EndpointResolverFunc" >}}) is provided as
a convenient way to wrap a resolver function to satisfy the `aws.EndpointResolver` interface.

For more information on endpoints and implementing a custom resolver, see
[Configuring Client Endpoints]({{% ref "/docs/configuring-sdk/endpoints.md" %}}).

### Invoking API Operations

The number of service client operation methods have been reduced significantly. The `<OperationName>Request`,
`<OperationName>WithContext`, and `<OperationName>` methods have all been consolidated into single operation method, `<OperationName>`.

#### Example

The following example shows how calls to the {{% alias service=s3 %}} PutObject operation would be migrated from
AWS SDK for Go to {{% alias sdk-go %}}.

```go
// V1

import "context"
import "github.com/aws/aws-sdk-go/service/s3"

// ...

client := s3.New(sess)

// Pattern 1
output, err := client.PutObject(&s3.PutObjectInput{
	// input parameters
})

// Pattern 2
output, err := client.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
	// input parameters
})

// Pattern 3
req, output := client.PutObjectRequest(context.TODO(), &s3.PutObjectInput{
	// input parameters
})
err := req.Send()
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

client : s3.NewFromConfig(cfg)

output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
	// input parameters
})
```

### Service Data Types

The top-level input and output types of an operation are found in the service client package. The input and output type
for a given operation follow the pattern of `<OperationName>Input` and `<OperationName>Output`, where `OperationName`
is the name of the operation you are invoking. For example the input and output shape for the
{{% alias service=S3 %}} PutObject operation are [PutObjectInput]({{< apiref "service/s3#PutObjectInput" >}}) and
[PutObjectOutput]({{< apiref "service/s3#PutObjectOutput" >}}) respectively.

All other service data types, other than input and output types, have been migrated to the `types` package located
under the service client package import path hierarchy. For example, the
[s3.AccessControlPolicy]({{< apiref v1="service/s3#AccessControlPolicy" >}}) type is now located at
[types.AccessControlPolicy]({{< apiref "service/s3/types#AccessControlPolicy" >}}).

#### Enumeration Values

The SDK now provides a typed experience for all API enumeration fields. Rather than using a string literal value copied
from the service API reference documentation, you can now use one of the concrete types found in the service client's
`types` package. For example, you can provide the {{% alias service=S3 %}} PutObjectInput operation
with an ACL to be applied on an object. In the AWS SDK for Go V1, this parameter was a `*string` type. In the
{{& alias sdk-go %}} this parameter is now a
[types.ObjectCannedACL]({{< apiref "service/s3/types#ObjectCannedACL" >}}). The `types` package
provides generated constants for the valid enumeration values that can be assigned to this field. For example
[types.ObjectCannedACLPrivate]({{< apiref "service/s3/types#ObjectCannedACLPrivate" >}}) is the constant for the
"private" canned ACL value. This value can be used in place of managing string constants within your application.

### Pointer Parameters

The AWS SDK for Go v1 required pointer references to be passed for all input parameters to service operations. The
{{% alias sdk-go %}} has simplified the experience with most services by removing the need to pass input values as
pointers where possible. This change means that many service clients operations no longer require your application
to pass pointer references for the following types: `uint8`, `uint16`, `uint32`, `int8`, `int16`, `int32`, `float32`,
`float64`, `bool`. Similarly, slice and map element types have been updated accordingly to reflect whether their
elements must be passed as pointer references.

The [aws]({{< apiref "aws" >}}) package contains helper functions for creating pointers for the Go built-in types, these
helpers should be used to more easily handle creating pointer types for these Go types. Similarly, helper methods are
provided for safely de-referencing pointer values for these types. For example, the
[aws.String]({{< apiref "aws#String" >}}) function converts from `string` &rArr; `*string`. Inversely,
the [aws.ToString]({{< apiref "aws#ToString" >}}) converts from `*string` &rArr; `string`. When upgrading your
application from AWS SDK for Go V1 to {{% alias sdk-go %}}, you must migrate usage of the helpers for
converting from the pointer types to the non-pointer variants. For example,
[aws.StringValue]({{< apiref v1="aws#StringValue" >}}) must be updated to `aws.ToString`.

### Errors Types

The {{% alias sdk-go %}} takes full advantage of the error wrapping functionality
[introduced in Go 1.13](https://blog.golang.org/go1.13-errors). Services that model error responses have generated
types available in their client's `types` package that can be used to test whether a client operation error was caused
by one of these types. For example {{% alias service=S3 %}} `GetObject` operation can return a `NoSuchKey` error if
attempting to retrieve an object key that doesn't exist. You can use [errors.As](https://golang.org/pkg/errors#As) to
test whether the returned operation error is a [types.NoSuchKey]({{< apiref "service/s3/types#NoSuchKey" >}}) error.
In the event a service does not model a specific type for an error, you can utilize the
[smithy.APIError]({{< apiref smithy="#APIError" >}}) interface type for inspecting the returned error code and message
from the service. This functionality replaces [awserr.Error]({{< apiref v1="aws/awserr#Error" >}}) and the other
[awserr]({{< apiref v1="aws/awserr" >}}) functionality from the AWS SDK for Go V1. For more details information on
handling errors see [Handling Errors]({{< ref "handling-errors" >}}).

#### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/awserr"
import "github.com/aws/aws-sdk-go/service/s3"

// ...

client := s3.New(sess)

output, err := s3.GetObject(&s3.GetObjectInput{
	// input parameters
})
if err != nil {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "NoSuchKey" {
			// handle NoSuchKey
		} else {
			// handle other codes
		}
		return
	}
	// handle a error
}
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/aws-sdk-go-v2/service/s3/types"
import "github.com/aws/smithy-go"

// ...

client := s3.NewFromConfig(cfg)

output, err := s3.GetObject(context.TODO(), &s3.GetObjectInput{
	// input parameters
})
if err != nil {
	var nsk *types.NoSuchKey
	if errors.As(err, &nsk) {
		// handle NoSuchKey error
		return
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		code := apiErr.ErrorCode()
		message := apiErr.ErrorMessage()
		// handle error code
		return
	}
	// handle error
	return
}
```

### Paginators

Service operation paginators are no longer invoked as methods on the service client. To use a paginator for an operation
you must construct a paginator for an operation using one of the paginator constructor methods. For example,
to use paginate over the {{% alias service=S3 %}} `ListObjectsV2` operation you must construct its paginator using the
[s3.NewListObjectsV2Paginator]({{< apiref "service/s3#NewListObjectsV2Paginator" >}}). This constructor returns a
[ListObjectsV2Paginator]({{< apiref "service/s3#ListObjectsV2Paginator" >}}) which provides the methods `HasMorePages`,
and `NextPage` for determining whether there are more pages to retrieve and invoking the operation to retrieve the next
page respectively. More details on using the SDK paginators can be found at
[]({{< ref "making-requests.md#using-paginators" >}}).

Let's look at an example of how to migrate from a AWS SDK for Go paginator to the {{% alias sdk-go %}} equivalent.

#### Example


```go
// V1

import "fmt"
import "github.com/aws/aws-sdk-go/service/s3"

// ...

client := s3.New(sess)

params := &s3.ListObjectsV2Input{
	// input parameters
}

totalObjects := 0
err := client.ListObjectsV2Pages(params, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
	totalObjects += len(output.Contents)
	return !lastPage
})
if err != nil {
	// handle error
}
fmt.Println("total objects:", totalObjects)
```


```go
// V2

import "context"
import "fmt"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

client := s3.NewFromConfig(cfg)

params := &s3.ListObjectsV2Input{
	// input parameters
}

totalObjects := 0
paginator := s3.NewListObjectsV2Paginator(client, params)
for paginator.HasMorePages() {
	output, err := paginator.NextPage(context.TODO())
	if err != nil {
		// handle error
	}
	totalObjects += len(output.Contents)
}
fmt.Println("total objects:", totalObjects)
```

### Waiters

Service operation waiters are no longer invoked as methods on the service client. To use a waiter you first construct
the desired waiter type, and then invoke the wait method. For example,
to wait for a {{% alias service=S3 %}} Bucket to exist, you must construct a `BucketExists` waiter. Use the
[s3.NewBucketExistsWaiter]({{< apiref "service/s3#NewBucketExistsWaiter" >}}) constructor to create a
[s3.BucketExistsWaiter]({{< apiref "service/s3#BucketExistsWaiter" >}}). The `s3.BucketExistsWaiter` provides a
`Wait` method which can be used to wait for a bucket to become available.

## Features

### {{% alias service=EC2 %}} Instance Metadata Service

The {{% alias sdk-go %}} provides an {{% alias service=EC2 %}} Instance Metadata Service (IMDS) client that you can use
to query the local IMDS when executing your application on an {{% alias service=EC2 %}} instance. The IMDS client is
a separate Go module that can be added to your application by using

```sh
go get github.com/aws/aws-sdk-go-v2/feature/ec2/imds
```

The client constructor and method operations have been updated to match the design of the other SDK service clients.

#### Example


```go
// V1

import "github.com/aws/aws-sdk-go/aws/ec2metadata"

// ...

client := ec2metadata.New(sess)

region, err := client.Region()
if err != nil {
	// handle error
}
```


```go
// V2

import "context"
import "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

// ...

client := imds.NewFromConfig(cfg)

region, err := client.GetRegion(context.TODO())
if err != nil {
	// handle error
}
```

### {{% alias service=S3 %}} Transfer Manager

The {{% alias service=S3 %}} transfer manager is available for managing uploads and downloads of objects
concurrently. This package is located in a Go module outside the service client import path. This module
can be retrieved by using `go get github.com/aws/aws-sdk-go-v2/feature/s3/manager`.

[s3.NewUploader]({{< apiref v1="service/s3/s3manager#NewUploader" >}}) and
[s3.NewUploaderWithClient]({{< apiref v1="service/s3/s3manager#NewUploaderWithClient" >}}) have been replaced with the
constructor method [manager.NewUploader]({{< apiref "feature/s3/manager#" >}}) for creating an Upload manager client.

[s3.NewDownloader]({{< apiref v1="service/s3/s3manager#NewDownloader" >}}) and
[s3.NewDownloaderWithClient]({{< apiref v1="service/s3/s3manager#NewDownloaderWithClient" >}}) have been replaced with a
single constructor method [manager.NewDownloader]({{< apiref "feature/s3/manager#NewDownloader" >}}) for creating a
Download manager client.

### {{% alias service=CFlong %}} Signing Utilities

The {{% alias sdk-go %}} provides {{% alias service=CFlong %}} signing utilities in a Go module outside the service
client import path. This module can be retrieved by using `go get`.

```sh
go get github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign
```
