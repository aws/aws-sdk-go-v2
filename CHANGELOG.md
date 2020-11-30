# Release 2020-11-30

## Breaking Change
* `codegen`: Add support for slice and maps generated with value members instead of pointer ([#887](https://github.com/aws/aws-sdk-go-v2/pull/887))
    * This update allow the SDK's code generation to be aware of API shapes and members that are not nullable, and can be rendered as value types by the code generation instead of pointer types.
    * Several API client parameter types will change from pointer members to value members for slice, map, number and bool member types.
    * See Migration notes for migrating to v0.30.0 with this change.
* `aws/transport/http`: Move aws.BuildableHTTPClient to HTTP transport package ([#898](https://github.com/aws/aws-sdk-go-v2/pull/898))
    * Moves the `BuildableHTTPClient` from the SDK's `aws` package to the `aws/transport/http` package as `BuildableClient` to with other HTTP specific utilities.
* `feature/cloudfront/sign`: Add CloudFront sign feature as module ([#884](https://github.com/aws/aws-sdk-go-v2/pull/884))
    * Moves `service/cloudfront/sign` package out of the `cloudfront` module, and into its own module as `github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign`.

## New Features
* `config`: Add a WithRetryer provider helper to the config loader ([#897](https://github.com/aws/aws-sdk-go-v2/pull/897))
    * Adds a `WithRetryer` configuration provider to the config loader as a convenience helper to set the `Retryer` on the `aws.Config` when its being loaded.
* `config`: Default to TLS 1.2 for HTTPS requests ([#892](https://github.com/aws/aws-sdk-go-v2/pull/892))
    * Updates the SDK's default HTTP client to use TLS 1.2 as the minimum TLS version for all HTTPS requests by default.

## Bug Fixes
* `config`: Fix AWS_CA_BUNDLE usage while loading default config ([#912](https://github.com/aws/aws-sdk-go-v2/pull/))
    * Fixes the `LoadDefaultConfig`'s configuration provider order to correctly load a custom HTTP client prior to configuring the client for `AWS_CA_BUNDLE` environment variable.
* `service/s3`: Fix signature mismatch error for s3 ([#913](https://github.com/aws/aws-sdk-go-v2/pull/913))
    * Fixes ([#883](https://github.com/aws/aws-sdk-go-v2/issues/883))
* `service/s3control`: 
    * Fix HostPrefix addition behavior for s3control ([#882](https://github.com/aws/aws-sdk-go-v2/pull/882))
        * Fixes ([#863](https://github.com/aws/aws-sdk-go-v2/issues/863))
    * Fix s3control error deserializer ([#875](https://github.com/aws/aws-sdk-go-v2/pull/875))
        * Fixes ([#864](https://github.com/aws/aws-sdk-go-v2/issues/864))

## Service Client Highlights
* Pagination support has been added to supported APIs. See [Using Operation Paginators](https://aws.github.io/aws-sdk-go-v2/docs/making-requests/#using-operation-paginators) in the Developer Guide. ([#885](https://github.com/aws/aws-sdk-go-v2/pull/885))
* Logging support has been added to service clients. See [Logging](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/logging/) in the Developer Guide. ([#872](https://github.com/aws/aws-sdk-go-v2/pull/872))
* `service`: Add support for pre-signed URL clients for S3, RDS, EC2 service ([#888](https://github.com/aws/aws-sdk-go-v2/pull/888))
    * `service/s3`: operations `PutObject` and `GetObject` are now supported with s3 pre-signed url client.
    * `service/ec2`: operation `CopySnapshot` is now supported with ec2 pre-signed url client. 
    * `service/rds`: operations `CopyDBSnapshot`, `CreateDBInstanceReadReplica`, `CopyDBClusterSnapshot`, `CreateDBCluster` are now supported with rds pre-signed url client.
* `service/s3`: Add support for S3 access point and S3 on outposts access point ARNs ([#870](https://github.com/aws/aws-sdk-go-v2/pull/870))
* `service/s3control`: Adds support for S3 on outposts access point and S3 on outposts bucket ARNs ([#870](https://github.com/aws/aws-sdk-go-v2/pull/870))

## Migrating from v2 preview SDK's v0.29.0 to v0.30.0

### aws.BuildableHTTPClient move
The `aws`'s `BuildableHTTPClient` HTTP client implementation was moved to `aws/transport/http` as `BuildableClient`. If your application used the `aws.BuildableHTTPClient` type, update it to use the `BuildableClient` in the `aws/transport/http` package.

### Slice and Map API member types 
This release includes several code generation updates for API client's slice map members. Using API modeling metadata the Slice and map members are now generated as value types instead of pointer types. For your application this means that for these types, the SDK no longer will have pointer member types, and have value member types.

To migrate to this change you'll need to remove the pointer handling for slice and map members, and instead use value type handling of the member values.

### Boolean and Number API member types
Similar to the slice and map API member types being generated as value, the SDK's code generation now has metadata where the SDK can generate boolean and number members as value type instead of pointer types. 

To migrate to this change you'll need to remove the pointer handling for numbers and boolean member types, and instead use value handling.

# Release 2020-10-30

## New Features
* Adds HostnameImmutable flag on aws.Endpoint to direct SDK if the associated endpoint is modifiable.([#848](https://github.com/aws/aws-sdk-go-v2/pull/848))

## Bug Fixes
* Fix SDK handling of xml based services - xml namespaces ([#858](https://github.com/aws/aws-sdk-go-v2/pull/858))
  * Fixes ([#850](https://github.com/aws/aws-sdk-go-v2/issues/850))

## Service Client Highlights
* API Clients have been bumped to version `v0.29.0`
    * Regenerate API Clients from update API models.
* Improve client doc generation.

## Core SDK Highlights
* Dependency Update: Updated SDK dependencies to their latest versions.

## Migrating from v2 preview SDK's v0.28.0 to v0.29.0
* API Clients ResolverOptions type renamed to EndpointResolverOptions 

# Release 2020-10-26

## New Features
* `service/s3`: Add support for Accelerate, and Dualstack ([#836](https://github.com/aws/aws-sdk-go-v2/pull/836))
* `service/s3control`: Add support for Dualstack ([#836](https://github.com/aws/aws-sdk-go-v2/pull/836))

## Service Client Highlights
* API Clients have been bumped to version `v0.28.0`
    * Regenerate API Clients from update API models.
* `service/s3`: Add support for Accelerate, and Dualstack ([#836](https://github.com/aws/aws-sdk-go-v2/pull/836))
* `service/s3control`: Add support for Dualstack ([#836](https://github.com/aws/aws-sdk-go-v2/pull/836))
* `service/route53`: Fix sanitizeURL customization to handle leading slash(`/`) [#846](https://github.com/aws/aws-sdk-go-v2/pull/846)
    * Fixes [#843](https://github.com/aws/aws-sdk-go-v2/issues/843)
* `service/route53`: Fix codegen to correctly look for operations that need sanitize url ([#851](https://github.com/aws/aws-sdk-go-v2/pull/851))

## Core SDK Highlights
* `aws/protocol/restjson`: Fix unexpected JSON error response deserialization ([#837](https://github.com/aws/aws-sdk-go-v2/pull/837))
    * Fixes [#832](https://github.com/aws/aws-sdk-go-v2/issues/832)
* `example/service/s3/listobjects`: Add example for Amazon S3 ListObjectsV2 ([#838](https://github.com/aws/aws-sdk-go-v2/pull/838))

# Release 2020-10-16

## New Features
* `feature/s3/manager`:
  * Initial `v0.1.0` release
  * Add the Amazon S3 Upload and Download transfer manager ([#802](https://github.com/aws/aws-sdk-go-v2/pull/802))

## Service Client Highlights
* Clients have been bumped to version `v0.27.0`
* `service/machinelearning`: Add customization for setting client endpoint with PredictEndpoint value if set ([#782](https://github.com/aws/aws-sdk-go-v2/pull/782))
* `service/s3`: Fix empty response body deserialization in case of error response ([#801](https://github.com/aws/aws-sdk-go-v2/pull/801))
  * Fixes xml deserialization util to correctly handle empty response body in case of an error response.
* `service/s3`: Add customization to auto fill Content-Md5 request header for Amazon S3 operations ([#812](https://github.com/aws/aws-sdk-go-v2/pull/812))
* `service/s3`: Add fallback to using HTTP status code for error code ([#818](https://github.com/aws/aws-sdk-go-v2/pull/818))
  * Adds falling back to using the HTTP status code to create a API Error code when not error code is received from the service, such as HeadObject.
* `service/route53`: Add support for deserialzing `InvalidChangeBatch` API error ([#792](https://github.com/aws/aws-sdk-go-v2/pull/792))
* `codegen`: Remove API client `Options` getter methods ([#788](https://github.com/aws/aws-sdk-go-v2/pull/788))
* `codegen`: Regenerate API Client modeled endpoints ([#791](https://github.com/aws/aws-sdk-go-v2/pull/791))
* `codegen`: Sort API Client struct member paramaters by required and alphabetical ([#787](https://github.com/aws/aws-sdk-go-v2/pull/787))
* `codegen`: Add package docs to API client modules ([#821](https://github.com/aws/aws-sdk-go-v2/pull/821))
* `codegen`: Rename `smithy-go`'s `smithy.OperationError` to `smithy.OperationInvokeError`.

## Core SDK Highlights
* `config`: 
  * Bumped to `v0.2.0`
  * Refactor Config Module, Add Config Package Documentation and Examples, Improve Overall SDK Readme ([#822](https://github.com/aws/aws-sdk-go-v2/pull/822))
* `credentials`:
  * Bumped to `v0.1.2`
  * Strip Monotonic Clock Readings when Comparing Credential Expiry Time ([#789](https://github.com/aws/aws-sdk-go-v2/pull/789))
* `ec2imds`:
  * Bumped to `v0.1.2`
  * Fix refreshing API token if expired ([#789](https://github.com/aws/aws-sdk-go-v2/pull/789))

## Migrating from v0.26.0 to v0.27.0

#### Configuration

The `config` module's exported types were trimmed down to add clarity and reduce confusion. Additional changes to the `config` module' helpers. 

* Refactored `WithCredentialsProvider`, `WithHTTPClient`, and `WithEndpointResolver` to functions instead of structs.
* Removed `MFATokenFuncProvider`, use `AssumeRoleCredentialOptionsProvider` for setting options for `stscreds.AssumeRoleOptions`.
* Renamed `WithWebIdentityCredentialProviderOptions` to `WithWebIdentityRoleCredentialOptions`
* Renamed `AssumeRoleCredentialProviderOptions` to `AssumeRoleCredentialOptionsProvider`
* Renamed `EndpointResolverFuncProvider` to `EndpointResolverProvider`

#### API Client
* API Client `Options` type getter methods have been removed. Use the struct members instead.
* The error returned by API Client operations was renamed from `smithy.OperationError` to `smithy.OperationInvokeError`.

# Release 2020-09-30

## Service Client Highlights
* Service clients have been bumped to `v0.26.0` simplify the documentation experience when using [pkg.go.dev](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2).
* `service/s3`: Disable automatic decompression of getting Amazon S3 objects with the `Content-Encoding: gzip` metadata header. ([#748](https://github.com/aws/aws-sdk-go-v2/pull/748))
  * This changes the SDK's default behavior with regard to making S3 API calls. The client will no longer automatically set the `Accept-Encoding` HTTP request header, nor will it automatically decompress the gzipped response when the `Content-Encoding: gzip` response header was received.
  * If you'd like the client to sent the `Accept-Encoding: gzip` request header, you can add this header to the API operation method call with the [SetHeaderValue](https://pkg.go.dev/github.com/awslabs/smithy-go/transport/http#SetHeaderValue). middleware helper.
* `service/cloudfront/sign`: Fix cloudfront example usage of SignWithPolicy ([#673](https://github.com/aws/aws-sdk-go-v2/pull/673))
  * Fixes [#671](https://github.com/aws/aws-sdk-go-v2/issues/671) documentation typo by correcting the usage of `SignWithPolicy`.

## Core SDK Highlights
* SDK core module released at `v0.26.0`
* `config` module released at `v0.1.1`
* `credentials` module released at `v0.1.1`
* `ec2imds` module released at `v0.1.1`


# Release 2020-09-28
## Announcements
We’re happy to share the updated clients for the v0.25.0 preview version of the AWS SDK for Go V2.

The updated clients leverage new developments and advancements within AWS and the Go software ecosystem at large since 
our original preview announcement. Using the new clients will be a bit different than before. The key differences are: 
simplified API operation invocation, performance improvements, support for error wrapping, and a new middleware architecture.
So below we have a guided walkthrough to help try it out and share your feedback in order to better influence the features 
you’d like to see in the GA version.

See [Announcement Blog Post](https://aws.amazon.com/blogs/developer/client-updates-in-the-preview-version-of-the-aws-sdk-for-go-v2/) for more details.

## Service Client Highlights
* Initial service clients released at version `v0.1.0`
## Core SDK Highlights
* SDK core module released at `v0.25.0`
* `config` module released at `v0.1.0`
* `credentials` module released at `v0.1.0`
* `ec2imds` module released at `v0.1.0`

## Migrating from v2 preview SDK's v0.24.0 to v0.25.0

#### Design changes

The v2 preview SDK `v0.25.0` release represents a significant stepping stone bringing the v2 SDK closer to its target design and usability. This release includes significant breaking changes to the v2 preview SDK. The updates in the `v0.25.0` release focus on refactoring and modularization of the SDK’s API clients to use the new [client design](https://github.com/aws/aws-sdk-go-v2/issues/438), updated request pipeline (aka [middleware](https://pkg.go.dev/github.com/awslabs/smithy-go/middleware)), refactored [credential providers](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/credentials), and [configuration loading](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/config) packages.

We've also bumped the minimum supported Go version with this release. Starting with v0.25.0 the SDK requires a minimum version of Go `v1.15`.

As a part of the refactoring done to v2 preview SDK some components have not been included in this update. The following is a non exhaustive list of features that are not available.

* API Paginators - [#439](https://github.com/aws/aws-sdk-go-v2/issues/439)
* API Waiters - [#442](https://github.com/aws/aws-sdk-go-v2/issues/442)
* Presign URL - [#794](https://github.com/aws/aws-sdk-go-v2/issues/794)
* Amazon S3 Upload and Download manager - [#802](https://github.com/aws/aws-sdk-go-v2/pull/802)
* Amazon DynamoDB's AttributeValue marshaler, and Expression package - [#790](https://github.com/aws/aws-sdk-go-v2/issues/790)
* Debug Logging - [#594](https://github.com/aws/aws-sdk-go-v2/issues/594)

We expect additional breaking changes to the v2 preview SDK in the coming releases. We expect these changes to focus on organizational, naming, and hardening the SDK's design for future feature capabilities after it is released for general availability.


#### Relocated Packages

In this release packages within the SDK were relocated, and in some cases those packages were converted to Go modules. The following is a list of packages have were relocated.

* `github.com/aws/aws-sdk-go-v2/aws/external` => `github.com/aws/aws-sdk-go-v2/config` module
* `github.com/aws/aws-sdk-go-v2/aws/ec2metadata` => `github.com/aws/aws-sdk-go-v2/ec2imds` module

The `github.com/aws/aws-sdk-go-v2/credentials` module contains refactored credentials providers.

* `github.com/aws/aws-sdk-go-v2/ec2rolecreds` => `github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds`
* `github.com/aws/aws-sdk-go-v2/endpointcreds` => `github.com/aws/aws-sdk-go-v2/credentials/endpointcreds`
* `github.com/aws/aws-sdk-go-v2/processcreds` => `github.com/aws/aws-sdk-go-v2/credentials/processcreds`
* `github.com/aws/aws-sdk-go-v2/stscreds` => `github.com/aws/aws-sdk-go-v2/credentials/stscreds`


#### Modularization

New modules were added to the v2 preview SDK to allow the components to be versioned independently from each other. This allows your application to depend on specific versions of an API client module, and take discrete updates from the SDK core and other API client modules as desired.

* [github.com/aws/aws-sdk-go-v2/config](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/config)
* [github.com/aws/aws-sdk-go-v2/credentials](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/credentials)
* Module for each API client, e.g. [github.com/aws/aws-sdk-go-v2/service/s3](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3)


#### API Clients

The following is a list of the major changes to the API client modules

* Removed paginators: we plan to add these back once they are implemented to integrate with the SDK's new API client design.
* Removed waiters: we need to further investigate how the V2 SDK should expose waiters, and how their behavior should be modeled.
* API Clients are now Go modules. When migrating to the v2 preview SDK `v0.25.0`, you'll need to add the API client's module to your application's go.mod file.
* API parameter nested types have been moved to a `types` package within the API client's module, e.g. `github.com/aws/aws-sdk-go-v2/service/s3/types` These types were moved to improve documentation and discovery of the API client, operation, and input/output types. For example Amazon S3's ListObject's operation [ListObjectOutput.Contents](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3/#ListObjectsOutput) input parameter is a slice of [types.Object](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3/types#Object).
* The client operation method has been renamed, removing the `Request` suffix. The method now invokes the operation instead of constructing a request, which needed to be invoked separately. The operation methods were also expanded to include functional options for providing operation specific configuration, such as modifying the request pipeline.

```go
result, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
    TableName: aws.String("exampleTable"),
}, func(o *Options) {
    // Limit operation calls to only 1 attempt.
    o.Retryer = retry.AddWithMaxAttempts(o.Retryer, 1)
})
```


#### Configuration

In addition to the `github.com/aws/aws-sdk-go-v2/aws/external` package being made a module at `github.com/aws/aws-sdk-go-v2/config`, the `LoadDefaultAWSConfig` function was renamed to `LoadDefaultConfig`.

The `github.com/aws/aws-sdk-go-v2/aws/defaults` package has been removed. Its components have been migrated to the `github.com/aws/aws-sdk-go-v2/aws` package, and `github.com/aws/aws-sdk-go-v2/config` module.


#### Error Handling

The `github.com/aws/aws-sdk-go-v2/aws/awserr` package was removed as a part of the SDK error handling refactor. The SDK now uses typed errors built around [Go v1.13](https://golang.org/doc/go1.13#error_wrapping)'s [errors.As](https://pkg.go.dev/errors#As) and [errors.Unwrap](https://pkg.go.dev/errors#Unwrap) features. All SDK error types that wrap other errors implement the `Unwrap` method. Generic v2 preview SDK errors created with `fmt.Errorf` use `%w` to wrap the underlying error.

The SDK API clients now include generated public error types for errors modeled for an API. The SDK will automatically deserialize the error response from the API into the appropriate error type. Your application should use `errors.As` to check if the returned error matches one it is interested in. Your application can also use the generic interface [smithy.APIError](https://pkg.go.dev/github.com/awslabs/smithy-go/#APIError) to test if the API client's operation method returned an API error, but not check against a specific error.

API client errors returned to the caller will use error wrapping to layer the error values. This allows underlying error types to be specific to their use case, and the SDK's more generic error types to wrap the underlying error.

For example, if an [Amazon DynamoDB](https://aws.amazon.com/dynamodb/) [Scan](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Scan) operation call cannot find the `TableName` requested, the error returned will contain [dynamodb.ResourceNotFoundException](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb/types#ResourceNotFoundException). The SDK will return this error value wrapped in a couple layers, with each layer adding additional contextual information such as [ResponseError](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws/transport/http#ResponseError) for AWS HTTP response error metadata , and [smithy.OperationError](https://pkg.go.dev/github.com/awslabs/smithy-go/#OperationError) for API operation call metadata.

```go
result, err := client.Scan(context.TODO(), params)
if err != nil {
    // To get a specific API error
    var notFoundErr *types.ResourceNotFoundException
    if errors.As(err, &notFoundErr) {
        log.Printf("scan failed because the table was not found, %v",
            notFoundErr.ErrorMessage())
    }

    // To get any API error
    var apiErr smithy.APIError
    if errors.As(err, &apiErr) {
        log.Printf("scan failed because of an API error, Code: %v, Message: %v",
            apiErr.ErrorCode(), apiErr.ErrorMessage())
    }

    // To get the AWS response metadata, such as RequestID
    var respErr *awshttp.ResponseError // Using import alias "awshttp" for package github.com/aws/aws-sdk-go-v2/aws/transport/http
    if errors.As(err, &respErr) {
        log.Printf("scan failed with HTTP status code %v, Request ID %v and error %v",
            respErr.HTTPStatusCode(), respErr.ServiceRequestID(), respErr)
    }

    return err
}
```

Logging an error value will include information from each wrapped error. For example, the following is a mock error logged for a Scan operation call that failed because the table was not found.

> 2020/10/15 16:03:37 operation error DynamoDB: Scan, https response error StatusCode: 400, RequestID: ABCREQUESTID123, ResourceNotFoundException: Requested resource not found


#### Endpoints

The `github.com/aws/aws-sdk-go-v2/aws/endpoints` has been removed from the SDK, along with all exported endpoint definitions and iteration behavior. Each generated API client now includes its own endpoint definition internally to the module.

API clients can optionally be configured with a generic [aws.EndpointResolver](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#EndpointResolver) via the [aws.Config.EndpointResolver](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config.EndpointResolver). If the API client is not configured with a custom endpoint resolver it will defer to the endpoint resolver the client module was generated with.
