Release v0.18.0 (2019-12-12)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Bugs
---
* `aws/endpoints`: aws/endpoints: Fix SDK resolving endpoint without region ([#420](https://github.com/aws/aws-sdk-go-v2/pull/420))
  * Fixes the SDK's endpoint resolve incorrectly resolving endpoints for a service when the region is empty. Also fixes the SDK attempting to resolve a service when the service value is empty.
  * Related to [aws/aws-sdk-go#2909](https://github.com/aws/aws-sdk-go/issues/2909)

Release v0.17.0 (2019-11-20)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Enhancements
---
* SDK minimum version requirement has been updated to Go 1.12 ([#432](https://github.com/aws/aws-sdk-go-v2/pull/432))

Release v0.16.0 (2019-11-12)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

Release v0.15.0 (2019-10-18)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

Release v0.14.0 (2019-10-08)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

Release v0.13.0 (2019-10-01)
===

### Services
* Synced the V2 SDK with latest AWS service API definitions.

### SDK Breaking changes
* This update includes breaking changes to how the DynamoDB AttributeValue (un)marshier handles empty collections.

### Deprecations
* `service/s3/s3crypto`: Deprecates the crypto client from the SDK ([#394](https://github.com/aws/aws-sdk-go-v2/pull/394))
  * s3crypto client is now deprecated and may be removed from the future versions of the SDK. 
* `aws`: Removes plugin credential provider ([#391](https://github.com/aws/aws-sdk-go-v2/pull/391))
  * Removing plugin credential provider from the v2 SDK developer preview. This feature may be made available as a separate module.
* Removes support for deprecated Go versions ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
  * Removes support for Go version specific files from the SDK. Also removes irrelevant build tags, and updates the README.md file. 
  * Raises the minimum supported version to Go 1.11 for the SDK. Older versions may work, but are not actively supported
  
### SDK Features
* `service/s3/s3manager`: Add Upload Buffer Provider ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Adds a new `BufferProvider` member for specifying how part data can be buffered in memory.
  * Windows platforms will now default to buffering 1MB per part to reduce contention when uploading files.
  * Non-Windows platforms will continue to employ a non-buffering behavior.
* `service/s3/s3manager`: Add Download Buffer Provider ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Adds a new `BufferProvider` member for specifying how part data can be buffered in memory when copying from the http response body.
  * Windows platforms will now default to buffering 1MB per part to reduce contention when downloading files.
  * Non-Windows platforms will continue to employ a non-buffering behavior.
* `service/dynamodb/dynamodbattribute`: New Encoder and Decoder Behavior for Empty Collections ([#401](https://github.com/aws/aws-sdk-go-v2/pull/401))
  * The `Encoder` and `Decoder` types have been enhanced to support the marshaling of empty structures, maps, and slices to and from their respective DynamoDB AttributeValues.
  * This change incorporates the behavior changes introduced via a marshal option in V1 ([#2834](https://github.com/aws/aws-sdk-go/pull/2834))

### SDK Enhancements
* `internal/awsutil`: Add suppressing logging sensitive API parameters ([#398](https://github.com/aws/aws-sdk-go-v2/pull/398))
  * Adds suppressing logging sensitive API parameters marked with the `sensitive` trait. This prevents the API type's `String` method returning a string representation of the API type with sensitive fields printed such as keys and passwords.
  * Related to [aws/aws-sdk-go#2310](https://github.com/aws/aws-sdk-go/pull/2310)
  * Fixes [#251](https://github.com/aws/aws-sdk-go-v2/issues/251)
* `aws/request` : Retryer is now a named field on Request. ([#393](https://github.com/aws/aws-sdk-go-v2/pull/393))
* `service/s3/s3manager`: Adds `sync.Pool` to allow reuse of part buffers for streaming payloads ([#404](https://github.com/aws/aws-sdk-go-v2/pull/404))
  * Fixes [#402](https://github.com/aws/aws-sdk-go-v2/issues/402) 
  * Uses the new behavior introduced in V1 [#2863](https://github.com/aws/aws-sdk-go/pull/2863) which allows the reuse of the sync.Pool across multiple Upload request that match part sizes.

### SDK Bugs
* `service/s3/s3manager`: Fix index out of range when a streaming reader returns -1 ([#378](https://github.com/aws/aws-sdk-go-v2/pull/378))
  * Fixes the S3 Upload Manager's handling of an unbounded streaming reader that returns negative bytes read.
* `internal/ini`: Fix ini parser to handle empty values [#406](https://github.com/aws/aws-sdk-go-v2/pull/406)
  * Fixes incorrect modifications to the previous token value of the skipper. Adds checks for cases where a skipped statement should be marked as complete and not be ignored.
  * Adds tests for nested and empty field value parsing, along with tests suggested in [aws/aws-sdk-go#2801](https://github.com/aws/aws-sdk-go/pull/2801)

Release v0.12.0 (2019-09-17)
===

Services
---
* Synced the V2 SDK with latest AWS service API definitions.

SDK Enhancements
---
* `aws/endpoints`: Expose DNSSuffix for partitions ([#369](https://github.com/aws/aws-sdk-go-v2/pull/369))
  * Exposes the underlying partition metadata's DNSSuffix value via the `DNSSuffix` method on the endpoint's `Partition` type. This allows access to the partition's DNS suffix, e.g. "amazon.com".
  * Fixes [#347](https://github.com/aws/aws-sdk-go-v2/issues/347)
* `private/protocol`: Add support for parsing fractional timestamp ([#367](https://github.com/aws/aws-sdk-go-v2/pull/367))
  * Fixes the SDK's ability to parse fractional unix timestamp values and adds tests.
  * Fixes [#365](https://github.com/aws/aws-sdk-go-v2/issues/365)
* `aws/ec2metadata`: Add marketplaceProductCodes to EC2 Instance Identity Document ([#374](https://github.com/aws/aws-sdk-go-v2/pull/374))
  * Adds `MarketplaceProductCodes` to the EC2 Instance Metadata's Identity Document. The ec2metadata client will now retrieve these values if they are available.
  * Related to: [aws/aws-sdk-go#2781](https://github.com/aws/aws-sdk-go/issues/2781)
* `aws`: Adds configurations to the default retryer ([#375](https://github.com/aws/aws-sdk-go-v2/pull/375))
  * Provides more customization options for retryer by adding a constructor for default Retryer which accepts functional options. Adds NoOpRetryer to support no retry behavior. Exposes members of default retryer.
  * Updates the underlying logic used by the default retryer to calculate jittered delay for retry. Handles int overflow for retry delay. 
  * Fixes [#370](https://github.com/aws/aws-sdk-go-v2/issues/370)
* `aws` : Refactors request retry behavior path logic ([#384](https://github.com/aws/aws-sdk-go-v2/pull/384))
  * Retry utilities now follow a consistent code path. aws.IsErrorRetryable is the primary entry point to determine if a request is retryable. 
  * Corrects sdk's behavior by not retrying errors with status code 501. Adds support for retrying the Kinesis API error, LimitExceededException.
  * Fixes [#372](https://github.com/aws/aws-sdk-go-v2/issues/372), [#145](https://github.com/aws/aws-sdk-go-v2/issues/145)
  
SDK Bugs
---
* `aws`: Fixes bug in calculating throttled retry delay ([#373](https://github.com/aws/aws-sdk-go-v2/pull/373))
  * The `Retry-After` duration specified in the request is now added to the Retry delay for throttled exception. Adds test for retry delays for throttled exceptions. Fixes bug where the throttled retry's math was off.
  * Fixes [#45](https://github.com/aws/aws-sdk-go-v2/issues/45)
* `aws` : Adds missing sdk error checking when seeking readers ([#379](https://github.com/aws/aws-sdk-go-v2/pull/379))
  * Adds support for nonseekable io.Reader. Adds support for streamed payloads for unsigned body request. 
  * Fixes [#371](https://github.com/aws/aws-sdk-go-v2/issues/371)
* `service/s3` : Fixes unexpected EOF error by s3manager ([#386](https://github.com/aws/aws-sdk-go-v2/pull/386))
  * Fixes bug which threw unexpected EOF error when s3 upload is performed for a file of maximum allowed size
  * Fixes [#316](https://github.com/aws/aws-sdk-go-v2/issues/316)
* `private/model` : Fixes generated API Reference docs links being invalid ([387](https://github.com/aws/aws-sdk-go-v2/pull/387))
  * Fixes [#327](https://github.com/aws/aws-sdk-go-v2/issues/327)

Release v0.11.0 (2019-08-22)
===

### Services
* Synced the V2 SDK with latest AWS service API definitions.
  * Fixes [#359](https://github.com/aws/aws-sdk-go-v2/issues/359)

### SDK Features

### SDK Enhancements
* `private/protocol`: Add support for TimestampFormat in protocols ([#358](https://github.com/aws/aws-sdk-go-v2/pull/358))
  * Adds support for the timestampForamt API model trait to the V2 SDK. The SDK will now generate API client parameters with the correct time format for APIs modeled with custom time stamp formats specified.
  * Fixes [#202](https://github.com/aws/aws-sdk-go-v2/issues/202)
  * Fixes [#286](https://github.com/aws/aws-sdk-go-v2/issues/286)
* `aws`: Add example for custom HTTP client idle connection options ([#350](https://github.com/aws/aws-sdk-go-v2/pull/350))
  * Adds example to the SDK for configuring custom HTTP client idle connection keep alive options.

### SDK Bugs
* `private/model/api`: Fix API doc being generated with wrong value ([#359](https://github.com/aws/aws-sdk-go-v2/pull/359))
  * Fixes the SDK's generated API documentation for structure member being generated with the wrong documentation value when the member was included multiple times in the model doc-2.json file, but under different types.
  * V2 port of to v1 [aws/aws-sdk-go#2748](https://github.com/aws/aws-sdk-go/issues/2748)
* `aws/ec2rolecreds`: Fix security creds path to include trailing slash ([#356](https://github.com/aws/aws-sdk-go-v2/pull/356))
  * Fixes the iamSecurityCredsPath var to include a trailing slash preventing redirects when making requests to the EC2 Instance Metadata service.
  * Fixes [#351](https://github.com/aws/aws-sdk-go-v2/issues/351)
* `service/dynamodb/expression`: Improved reporting of bad key conditions ([#360](https://github.com/aws/aws-sdk-go-v2/pull/360))
  * Improved error reporting when invalid key conditions are constructed using KeyConditionBuilder

Release v0.10.0 (2019-07-24)
===

### Services
* Synced the V2 SDK with latest AWS service API definitions.
* Fixes [#341](https://github.com/aws/aws-sdk-go-v2/issues/341)
* Fixes [#342](https://github.com/aws/aws-sdk-go-v2/issues/342)

### SDK Breaking Changes
* `aws`: Add default HTTP client instead of http.DefaultClient/Transport ([#315](https://github.com/aws/aws-sdk-go-v2/pull/315))
  * Adds a new BuildableHTTPClient type to the SDK's aws package. The type uses the builder pattern with immutable changes. Modifications to the buildable client create copies of the client.  Adds a HTTPClient interface to the aws package that the SDK will use as an abstraction over the specific HTTP client implementation. The SDK will default to the BuildableHTTPClient, but a *http.Client can be also provided for custom configuration.  When the SDK's aws.Config.HTTPClient value is a BuildableHTTPClient the SDK will be able to use API client specific request timeout options.
  * Fixes [#279](https://github.com/aws/aws-sdk-go-v2/issues/279)
  * Fixes [#269](https://github.com/aws/aws-sdk-go-v2/issues/269)

### SDK Enhancements
* `service/s3/s3manager`: Update S3 Upload Multipart location ([#324](https://github.com/aws/aws-sdk-go-v2/pull/324))
  * Updates the Location returned value of S3 Upload's Multipart UploadOutput type to be consistent with single part upload URL. This update also brings the multipart upload Location inline with the S3 object URLs created by the SDK.
  * Fixes [#323](https://github.com/aws/aws-sdk-go-v2/issues/323)
  * V2 Port [aws/aws-sdk-go#2453](https://github.com/aws/aws-sdk-go/issues/2453)

### SDK Bugs
* `private/model`: Handles empty map vs unset map behavior in send request ([#337](https://github.com/aws/aws-sdk-go-v2/pull/337))
  * Updated shape marshal model to handle the empty map vs nil map behavior. Adding a test case to assert behavior when a user sends an empty map vs nil map.
  * Fix [#332](https://github.com/aws/aws-sdk-go-v2/issues/332)
* `service/rds`: Fix presign URL for same region ([#331](https://github.com/aws/aws-sdk-go-v2/pull/331)) 
  * Fixes RDS no-autopresign URL for same region issue for aws-sdk-go-v2. Solves the issue by making sure that the presigned URLs are not created, when the source and destination regions are the same. Added and updated the tests accordingly.
  * Fix [#271](https://github.com/aws/aws-sdk-go-v2/issues/271)
* `private/protocola/json/jsonutil`: Fix Unmarshal map[string]bool ([#320](https://github.com/aws/aws-sdk-go-v2/pull/320))
  * Fixes the JSON unmarshaling of maps of bools. The unmarshal case was missing the condition for bool value, in addition the bool pointer.
  * Fix [#319](https://github.com/aws/aws-sdk-go-v2/issues/319)

Release v0.9.0 (2019-05-28)
===

### Services
* Synced the V2 SDK with latest AWS service API definitions.
* Fixes [#304](https://github.com/aws/aws-sdk-go-v2/issues/304)
* Fixes [#295](https://github.com/aws/aws-sdk-go-v2/issues/295)

### SDK Breaking changes
This update includes multiple breaking changes to the SDK. These updates improve the SDK's usability, consistency.

#### Client type name
The API client type is renamed to `Client` for consistency, and remove stutter between package and client type name. Using Amazon S3 API client as an example, the `s3.S3` type is renamed to `s3.Client`.

#### New API operation response type
API operations' `Request.Send` method now returns a Response type for the specific operation. The Response type wraps the operation's Output parameter, and includes a method for the response's metadata such as RequestID. The Output type is an anonymous embedded field within the Output type. If your application was passing the Output value around you'll need to extract it directly, or pass the Response type instead.

#### New API operation paginator utility
This change removes the `Paginate` method from API operation Request types, (e.g. ListObjectsRequest). A new Paginator constructor is added that can be used to page these operations. To update your application to use the new pattern, where `Paginate` was being called, replace this with the Paginator type's constructor. The usage of the returned Paginator type is unchanged.

```go
req := svc.ListObjectsRequest(params)
p := req.Paginate()
```

Is updated to to use the Paginator constructor instead of Paginate method.

```go
req := svc.ListObjectsRequest(params)
p := s3.NewListObjectsPaginator(req)
```

#### Other changes
  * Standardizes API client package name to be based on the API model's `ServiceID`.
  * Standardizes API client operation input and output type names.
  * Removes `endpoints` package's service identifier constants. These values were unstable. Each API client package contains an `EndpointsID` constant that can be used for service specific endpoint lookup.
  * Fix API endpoint lookups to use the API's modeled `EndpointsID` (aka `enpdointPrefix`). Searching for API endpoints in the `endpoints` package should use the API client package's, `EndpointsID`.

### SDK Enhancements
*  Update CI tests to ensure all codegen changes are accounted for in PR ([#183](https://github.com/aws/aws-sdk-go-v2/issues/183))
  * Updates the CI tests to ensure that any code generation changes are accounted for in the PR, and that there were no mistaken changes made without also running code generation. This change should also help ensure that code generation order is stable, and there are no ordering issues with the SDK's codegen.
  * Related [aws/aws-sdk-go#1966](https://github.com/aws/aws-sdk-go/issues/1966)

### SDK Bugs
* `service/dynamodb/expression`: Fix Builder with KeyCondition example ([#306](https://github.com/aws/aws-sdk-go-v2/issues/306))
  * Fixes the ExampleBuilder_WithKeyCondition example to include the ExpressionAttributeNames member being set.
  * Fixes [#285](https://github.com/aws/aws-sdk-go-v2/issues/285)
* `aws/defaults`: Fix UserAgent execution environment key ([#307](https://github.com/aws/aws-sdk-go-v2/issues/307))
  * Fixes the SDK's UserAgent key for the execution environment.
  * Fixes [#276](https://github.com/aws/aws-sdk-go-v2/issues/276)
* `private/model/api`: Improve SDK API reference doc generation ([#309](https://github.com/aws/aws-sdk-go-v2/issues/309))
  * Improves the SDK's generated documentation for API client, operation, and types. This fixes several bugs in the doc generation causing poor formatting, an difficult to read reference documentation.
  * Fix [#308](https://github.com/aws/aws-sdk-go-v2/issues/308)
  * Related [aws/aws-sdk-go#2617](https://github.com/aws/aws-sdk-go/issues/2617)

Release v0.8.0 (2019-04-25)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Breaking changes
* Update SDK API operation request send to required Context ([#265](https://github.com/aws/aws-sdk-go-v2/pull/265))
  * Updates the SDK's API operation request Send method to require a context.Context value when called. This is done to encourage best applications to use Context for cancellation and request tracing.  Standardizing on this pattern will also help reduce code paths which accidentally do not have the Context causing the cancellation and tracing chain to be lost. Leading to difficult to trace down losses of cancellation and tracing within an application.
  * Fixes [#264](https://github.com/aws/aws-sdk-go-v2/pull/264)

### SDK Enhancements
* Update README.md for getting SDK without Go Modules
  * Updates the README.md with instructions how to get the SDK without Go Modules enabled, or using the SDK within a GOPATH with Go 1.11, and Go 1.12.
* Refactor SDK's integration tests to be code generated ([#283](https://github.com/aws/aws-sdk-go-v2/pull/283))
* `aws`: Add RequestThrottledException to set of throttled exceptions ([#292](https://github.com/aws/aws-sdk-go-v2/pull/292))
* `private/model/api`: Backfill authtype, STS and Cognito Identity ([#293](https://github.com/aws/aws-sdk-go-v2/pull/293))
  * Backfills the authtype=none modeled trait for STS and Cognito Identity services. This removes the in code customization for these two services' APIs that should not be signed.

### SDK Bugs
* Fix HTTP endpoint credential provider test for unresolved hosts ([#262](https://github.com/aws/aws-sdk-go-v2/pull/262))
  * Fixes the HTTP endpoint credential provider's tests to check for a host that resolves to no addresses.
* `example/service/s3/mockPaginator`: Update example to not use internal pkg ([#278](https://github.com/aws/aws-sdk-go-v2/pull/278))
  * Updates the SDK's S3 Mock Paginator example to not use internal SDK packages and instead use the SDK's provided defaults package for default configuration.
  * Fixes [#116](https://github.com/aws/aws-sdk-go-v2/issues/116)
* Cleanup go mod unused dependencies ([#284](https://github.com/aws/aws-sdk-go-v2/pull/284))
* `service/s3/s3manager`: Fix brittle Upload unit test ([#288](https://github.com/aws/aws-sdk-go-v2/pull/288))
* `aws/ec2metadata`: Fix EC2 Metadata client panic with debug logging ([#290](https://github.com/aws/aws-sdk-go-v2/pull/290))
  * Fixes a panic that could occur within the EC2 Metadata client when both AWS_EC2_METADATA_DISABLED env var is set and log level is LogDebugWithHTTPBody. The SDK's client response body debug functionality would panic because the Request.HTTPResponse value was not specified.
* `aws`: Fix RequestUserAgent test to be stable ([#289](https://github.com/aws/aws-sdk-go-v2/pull/289))
* `private/protocol/rest`: Trim space in header key and value ([#291](https://github.com/aws/aws-sdk-go-v2/pull/291))
  * Fixes a bug when using S3 metadata where metadata values with leading spaces would trigger request signature validation errors when the request is received by the service.


Release v0.7.0 (2019-01-03)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Enhancements
* deps: Update SDK to latest go-jmespath ([#254](https://github.com/aws/aws-sdk-go-v2/pull/254))

### SDK Bugs
* `internal/ini`: Fix bug on trimming rhs spaces closes ([#260](https://github.com/aws/aws-sdk-go-v2/pull/260))
  * Fixes a bug trimming RHS spaces not being read correctly from the ini file.
  * Fix [#259](https://github.com/aws/aws-sdk-go-v2/pull/259)

Release v0.6.0 (2018-12-03)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Bugs
* Updates the SDK's release tagging scheme to use `v0` until the v2 SDK reaches
	* General Availability (GA). This allows the SDK to be used with Go 1.11 modules. Post GA, v2 SDK's release tagging version will most likely follow a `v1.<x>.<y>` patter.
	* Fixes [#221](https://github.com/aws/aws-sdk-go-v2/issues/221)

Release v0.5.0 (2018-09-19)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Bugs
* Fix SDK Go 1.11 connection reset handling (#207)
	* Fixes how the SDK checks for connection reset errors when making API calls to be compatiable with Go 1.11.
* `aws/signer/v4`: Fix X-Amz-Content-Sha256 being in to query for presign (#188)
	* Fixes the bug which would allow the X-Amz-Content-Sha256 header to be promoted to the query string when presigning a S3 request.  This bug also was preventing users from setting their own sha256 value for a presigned URL. Presigned requests generated with the custom sha256 would of always failed with invalid signature.
	* Related to aws/aws-sdk-go#1974

### SDK Enhancements
* Cleanup SDK README and example documenation.
* `service/s3/s3manager`: Add doc for sequential download (#201)
	Adds documentation for downloading object sequentially with the S3 download manager.
* `aws/credentials`: Update Credentials cache to have less overhead (#184)
	* Updates the Credentials type's cache of the CredentialsValue to be synchronized with an atomic value in addition to the Mutex. This reduces the overhead applications will encounter when many concurrent API requests are being made.
	* Related to: aws/aws-sdk-go#1973
* `service/dynamodb/dynamodbattribute`: Add support for custom struct tag keys (#203)
	* Adds support for (un)marshaling Go types using custom struct tag keys. The new `MarshalOptions.TagKey` allows the user to specify the tag key to use when (un)marshaling struct fields.  Adds support for struct tags such as `yaml`, `toml`, etc. Support for these keys are in name only, and require the tag value format and values to be supported by the package's Marshalers.
* `internal/ini`: Add custom INI parser for shared config/credentials file (#209)
	* Related to: aws/aws-sdk-go#2024

Release v0.4.0 (2018-05-25)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Bugs
* `private/protocol/xml/xmlutil`: Fix XML unmarshaler not correctly unmarshaling list of timestamp values ([#166](https://github.com/aws/aws-sdk-go-v2/pull/166))
	* Fixes a bug in the XML unmarshaler that would incorrectly try to unmarshal "time.Time" parameters that did not have the struct tag type on them.
	* Related to [aws/aws-sdk-go#1894](https://github.com/aws/aws-sdk-go/pull/1894)
* `service/s3`: Fix typos for migrated S3 specific config options ([#173](https://github.com/aws/aws-sdk-go-v2/pull/173))
	* Updates the S3 specific config error messages to the correct fields.
* `aws/endpoints`: Fix SDK endpoint signing name resolution ([#181](https://github.com/aws/aws-sdk-go-v2/pull/181))
	* Fixes how the SDK derives service signing names. If the signing name is not modeled in the endpoints package the service will fallback to the signing name modeled in the service model.
	* Fix [#163](https://github.com/aws/aws-sdk-go-v2/pull/163)
	* Fix [#153](https://github.com/aws/aws-sdk-go-v2/pull/153)
	* Related to [aws/aws-sdk-go#1854](https://github.com/aws/aws-sdk-go/pull/1854)
* `service/s3`: remove SelectContent until EventStream supported ([#175](https://github.com/aws/aws-sdk-go-v2/pull/175])
	* S3's SelectContent API is not functional in the SDK yet, and was not supposed to be generated until EventStream support is available.
	* Related to [aws/aws-sdk-go#1941](https://github.com/aws/aws-sdk-go/pull/1941)

### SDK Enhancements
* `service/s3/s3manager/s3manageriface`: Add WithIterator to mock interface ([#156](https://github.com/aws/aws-sdk-go-v2/pull/156))
	* Updates the `DownloaderAPI` and `UploaderAPI` mocking interfaces to have parity with the concrete types.
	* Fixes [#155](https://github.com/aws/aws-sdk-go-v2/issues/155)


Release v0.3.0 (2018-03-08)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### Breaking Changes
* `private/mode/api`: Refactor service paginator helpers to iterator pattern ([#119](https://github.com/aws/aws-sdk-go-v2/pull/119))
	* Refactors the generated service client paginators to be an iterator pattern. This pattern improves usability while removing the need for callbacks.
	* See the linked PR for an example.
* `private/model/api`: Removes setter helpers from service API types ([#101](https://github.com/aws/aws-sdk-go-v2/pull/101))
	* Removes the setter helper methods from service API types. Removing clutter and noise from the API type's signature.
	* Based on feedback [#81][https://github.com/aws/aws-sdk-go-v2/issues/81]
* `aws`: Rename CanceledErrorCode to ErrCodeRequestCanceled ([#131](https://github.com/aws/aws-sdk-go-v2/pull/131))
	* Renames CanceledErrorCode to correct naming scheme.

### SDK Bugs
* `internal/awsutil`: Fix DeepEqual to consider string alias type equal to string ([#102](https://github.com/aws/aws-sdk-go-v2/pull/102))
	* Fixes SDK waiters not detecting the correct condition is met. [#92](https://github.com/aws/aws-sdk-go-v2/issues/92)
* `aws/external`: Fix EnvConfig misspelled container endpoint path getter ([#106](https://github.com/aws/aws-sdk-go-v2/pull/106))
	* This caused the type to not satisfy the ContainerCredentialsEndpointPathProvider interface.
	* Fixes [#105](https://github.com/aws/aws-sdk-go-v2/issues/105)
* `service/s3/s3crypto`: Fix S3Crypto's handling of TagLen ([#107](https://github.com/aws/aws-sdk-go-v2/pull/107))
	* Fixes the S3Crypto's handling of TagLen to only be set if present.
	* V2 Fix for [aws/aws-sdk-go#1742](https://github.com/aws/aws-sdk-go/issues/1742)
* `private/model/api`: Update SDK service client initialization documentation ([#141](https://github.com/aws/aws-sdk-go-v2/pull/141))
	* Updates the SDK's service initialization doc template to reflect the v2 SDK's configuration update change from v1.
	* Related to [#136](https://github.com/aws/aws-sdk-go-v2/issues/136)

### SDK Enhancements
* `aws`: Improve pagination unit tests ([#97](https://github.com/aws/aws-sdk-go-v2/pull/97))
	* V2 port of [aws/aws-sdk-go#1733](https://github.com/aws/aws-sdk-go/pull/1733)
* `aws/external`: Add example for shared config and static credential helper ([#109](https://github.com/aws/aws-sdk-go-v2/pull/109))
	* Adds examples for the  config helpers; WithSharedConfigProfile, WithCredentialsValue, WithMFATokenFunc.
* `private/model/api`: Add validation to prevent collision of api defintions ([#112](https://github.com/aws/aws-sdk-go-v2/pull/112)) 
	* V2 port of [aws/aws-sdk-go#1758](https://github.com/aws/aws-sdk-go/pull/1758)
* `aws/ec2metadata`: Add support for AWS_EC2_METADATA_DISABLED env var ([#128](https://github.com/aws/aws-sdk-go-v2/pull/128))
	* When this environment variable is set. The SDK's EC2 Metadata Client will not attempt to make requests. All requests made with the EC2 Metadata Client will fail.
	* V2 port of [aws/aws-sdk-go#1799](https://github.com/aws/aws-sdk-go/pull/1799)
* Add code of conduct ([#138](https://github.com/aws/aws-sdk-go-v2/pull/138))
* Update SDK README dep usage ([#140](https://github.com/aws/aws-sdk-go-v2/pull/140))

Release v0.2.0 (2018-01-15)
===

### Services
* Synced the V2 SDK with latests AWS service API definitions.

### SDK Bugs
* `service/s3/s3manager`: Fix Upload Manger's UploadInput fields ([#89](https://github.com/aws/aws-sdk-go-v2/pull/89))
	* Fixes [#88](https://github.com/aws/aws-sdk-go-v2/issues/88)
* `aws`: Fix Pagination handling of empty string NextToken ([#94](https://github.com/aws/aws-sdk-go-v2/pull/94))
	* Fixes [#84](https://github.com/aws/aws-sdk-go-v2/issues/84)


Release v0.1.0 (2017-12-21)
===

## What has changed?

Our focus for the 2.0 SDK is to improve the SDK’s development experience and performance, make the SDK easy to extend, and add new features. The changes made in the Developer Preview target the major pain points of configuring the SDK and using AWS service API calls. Check out the SDK for details on pending changes that are in development and designs we’re discussing.

The following are some of the larger changes included in the AWS SDK for Go 2.0 Developer Preview.

### SDK configuration

The 2.0 SDK simplifies how you configure the SDK's service clients by using a single `Config` type. This simplifies the `Session` and `Config` type interaction from the 1.x SDK. In addition, we’ve moved the service-specific configuration flags to the individual service client types. This reduces confusion about where service clients will use configuration members.

We added the external package to provide the utilities for you to use external configuration sources to populate the SDK's `Config` type. External sources include environmental variables, shared credentials file (`~/.aws/credentials`), and shared config file (`~/.aws/config`). By default, the 2.0 SDK will now automatically load configuration values from the shared config file. The external package also provides you with the tools to customize how the external sources are loaded and used to populate the `Config` type.

You can even customize external configuration sources to include your own custom sources, for example, JSON files or a custom file format.

Use `LoadDefaultAWSConfig` in the external package to create the default `Config` value, and load configuration values from the external configuration sources.

```go
cfg, err := external.LoadDefaultAWSConfig()
```

To specify the shared configuration profile load used in code, use the `WithSharedConfigProfile` helper passed into `LoadDefaultAWSConfig` with the profile name to use.

```go
cfg, err := external.LoadDefaultAWSConfig(
	external.WithSharedConfigProfile("gopher")
)
```

Once a `Config` value is returned by `LoadDefaultAWSConfig`, you can set or override configuration values by setting the fields on the `Config` struct, such as `Region`.

```go
cfg.Region = endpoints.UsWest2RegionID
```

Use the `cfg` value to provide the loaded configuration to new AWS service clients.

```go
svc := dynamodb.New(cfg)
```

### API request methods

The 2.0 SDK consolidates several ways to make an API call, providing a single request constructor for each API call. This means that your application will create an API request from input parameters, then send it. This enables you to optionally modify and configure how the request will be sent. This includes, but isn’t limited to, modifications such as setting the `Context` per request, adding request handlers, and enabling logging.

Each API request method is suffixed with `Request` and returns a typed value for the specific API request.

As an example, to use the Amazon Simple Storage Service GetObject API, the signature of the method is:

```go
func (c *S3) GetObjectRequest(*s3.GetObjectInput) *s3.GetObjectRequest
```

To use the GetObject API, we pass in the input parameters to the method, just like we would with the 1.x SDK. The 2.0 SDK's method will initialize a `GetObjectRequest` value that we can then use to send our request to Amazon S3.

```go
req := svc.GetObjectRequest(params)

// Optionally, set the context or other configuration for the request to use
req.SetContext(ctx)

// Send the request and get the response
resp, err := req.Send()
```

### API enumerations

The 2.0 SDK uses typed enumeration values for all API enumeration fields. This change provides you with the type safety that you and your IDE can leverage to discover which enumeration values are valid for particular fields. Typed enumeration values also provide a stronger type safety story for your application than using string literals directly. The 2.0 SDK uses string aliases for each enumeration type. The SDK also generates typed constants for each enumeration value. This change removes the need for enumeration API fields to be pointers, as a zero-value enumeration always means the field isn’t set.

For example, the Amazon Simple Storage Service PutObject API has a field, `ACL ObjectCannedACL`. An `ObjectCannedACL` string alias is defined within the s3 package, and each enumeration value is also defined as a typed constant. In this example, we want to use the typed enumeration values to set an uploaded object to have an ACL of `public-read`. The constant that the SDK provides for this enumeration value is `ObjectCannedACLPublicRead`.

```go
svc.PutObjectRequest(&s3.PutObjectInput{
	Bucket: aws.String("myBucket"),
	Key:    aws.String("myKey"),
	ACL:    s3.ObjectCannedACLPublicRead,
	Body:   body,
})
```

### API slice and map elements

The 2.0 SDK removes the need to convert slice and map elements from values to pointers for API calls. This will reduce the overhead of needing to use fields with a type, such as `[]string`, in API calls. The 1.x SDK's pattern of using pointer types for all slice and map elements was a significant pain point for users, requiring them to convert between the types. The 2.0 SDK does away with the pointer types for slices and maps, using value types instead.

The following example shows how value types for the Amazon Simple Queue Service AddPermission API's `AWSAccountIds` and `Actions` member slices are set.

```go
svc := sqs.New(cfg)

svc.AddPermission(&sqs.AddPermissionInput{
	AWSAcountIds: []string{
		"123456789",
	},
	Actions: []string{
		"SendMessage",
	},

	Label:    aws.String("MessageSender"),
	QueueUrl: aws.String(queueURL)
})
```


