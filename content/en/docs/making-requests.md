---
title: "Using the AWS SDK for Go V2 with AWS Services"
linkTitle: "Using AWS Services"
date: "2020-11-12"
description: "Construct service clients and make operation calls to send requests to AWS services."
weight: 5
---

To make calls to an AWS service, you must first construct a service client instance. A service client
provides low-level access to every API action for that service. For example, you create an {{% alias service=S3 %}}
service client to make calls to {{% alias service=S3 %}} APIs.

When you call service operations, you pass in input parameters as a struct. A successful call will result in an
output struct containing the service API response. For example, after you successfully call an {{% alias service=S3 %}}
create bucket action, the action returns an output struct with the bucket's location.

For the list of service clients, including their methods and parameters, see the [{{% alias sdk-api %}}]({{< apiref "#" >}})

## Constructing a Service Client

Service clients can be constructed using either the `New` or `NewFromConfig` functions available in service client's
Go package. Each function will return a `Client` struct type containing the methods for invoking the service APIs.
The `New` and `NewFromConfig` each provide the same set of configurable options for constructing a service client, but
provide slightly different construction patterns that we will look at in the following sections.

### NewFromConfig

`NewFromConfig` function provides a consistent interface for constructing service clients using the 
[aws.Config]({{< apiref "aws#Config" >}}). An `aws.Config` can be loaded using the
[config.LoadDefaultConfig]({{< apiref "config#LoadDefaultConfig" >}}). For more information on constructing 
an `aws.Config` see [Configure the SDK]({{% ref "configuring-sdk" %}}). The following example shows how to construct
an {{% alias service=S3 %}} service client using the `aws.Config`and the `NewFromConfig` function:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	panic(err)
}

client := s3.NewFromConfig(cfg)
```

#### Overriding Configuration
`NewFromConfig` can take one or more functional arguments that can mutate a client's configuration `Options` struct.
This allows you to make specific overrides such as changing the Region, or modifying service specific options such as
{{% alias service=S3 %}} `UseAccelerate` option. For example:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	panic(err)
}

client := s3.NewFromConfig(cfg, func(o *s3.Options) {
	o.Region = "us-west-2"
	o.UseAccelerate = true
})
```

Overrides to the client `Options` value is determined by the order that the functional arguments are given to
`NewFromConfig`.

### New

{{% pageinfo color="info" %}}
`New` is considered a more advanced form of client construction. We recommend you use `NewFromConfig` for client
construction, as it allows construction using the `aws.Config` struct. This removes the need to construct an `Options`
struct instance for each service client your application requires.
{{% /pageinfo %}}

`New` function is a client constructor provides an interface for constructing clients using only the client packages
`Options` struct for defining the client's configuration options. For example to construct {{% alias service=S3 %}}
client using `New`:

```go
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

client := s3.New(s3.Options{
	Region:      "us-west-2",
	Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
})
```

#### Overriding Configuration

`New` can take one or more functional arguments that can mutate a client's configuration `Options` struct.
This allows you to make specific overrides such as changing the Region or modifying service specific options
such as {{% alias service=S3 %}} `UseAccelerate` option. For example:

```go
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/credentials"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

options := s3.Options{
    Region:      "us-west-2",
    Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
}

client := s3.New(options, func(o *s3.Options) {
	o.Region = "us-east-1"
	o.UseAccelerate = true
})
```

Overrides to the client `Options` value is determined by the order that the functional arguments are given to `New`.

## Calling Service Operations

After you have a service client instance, you can use it to call a service's operations. For example to call the
{{% alias service=S3 %}} `GetObject` operation:

```go
response, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
	Bucket: aws.String("my-bucket"),
	Key:    aws.String("obj-key"),
})
```

When you call a service operation, the SDK synchronously validates the input, serializes the request, signs it with your
credentials, sends it to AWS, and then deserializes a response or an error. In most cases, you can call service
operations directly. Each service operation client method will return an operation response struct, and an
error interface type. You should always check `error` type to determine if an error occurred before attempting to access
the service operation's response struct.

### Passing Parameters to a Service Operation

Each service operation method takes a [context.Context](https://golang.org/pkg/context/#Context) value that can be
used for setting request deadlines that will be honored by the SDK. In addition, each service operation will take a
`<OperationName>Input` struct found in the service's respective Go package. You pass in API input parameters using
the operation input struct.

Operation input structures can have input parameters such as the standard Go numerics, boolean, string, map, and list
types. In more complex API operations a service might have more complex modeling of input parameters. These other types
such as service specific structures and enum values are found in the service's `types` Go package.

In addition, services might distinguish between the default value of a Go type and whether the value was set or not by
the user. In these cases, input parameters might require you to pass a pointer reference to the type in question. For
standard Go types like numerics, boolean, and string there are `<Type>` and `From<Type>` convenience functions available
in the [aws]({{< apiref aws >}}) to ease this conversion. For example [aws.String]({{% apiref "aws#String" %}}) can be
used to convert a `string` to a `*string` type for input parameters that require a pointer to a string. Inversely
[aws.ToString]({{% apiref "aws#ToString" %}}) can be used to transform a `*string` to a `string` while providing
protection from dereferencing a nil pointer. The `To<Type>` functions are helpful when handling service responses.

Let's look at an example of how we can use an {{% alias service=S3 %}} client to call the `GetObject` API, and construct
our input using the `types` package, and `aws.<Type>` helpers.

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/aws-sdk-go-v2/service/s3/types"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	panic(err)
}

client := s3.NewFromConfig(cfg)

resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
    Bucket:       aws.String("my-bucket"),
    Key:          aws.String("keyName"),
    RequestPayer: types.RequestPayerRequester,
})
```

### Overriding Client Options For Operation Call {#OverrideClientOptionsForOperation}

Similar to how client operation options can be modified during construction of a client using functional arguments,
the client options can be modified at the time the operation method is called by providing one or more functional
arguments to the service operation method. This action is concurrency safe and will not affect other concurrent
operations on the client.

For example to override the client region from "us-west-2" to "us-east-1":

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
if err != nil {
	log.Printf("error: %v", err)
	return
}

client := s3.NewFromConfig(cfg)

params := &s3.GetObjectInput{
	// ...
}

resp, err := client.GetObject(context.TODO(), params, func(o *Options) {
	o.Region = "us-east-1"
})
```

### Handling Operation Responses {#HandlingOperationResponses}

Each service operation has an associated output struct that contains the service's operation response members.
The output struct follows the following naming pattern `<OperationName>Output`. Some operations might have no members
defined for their operation output. After calling a service operation, the return `error` argument type should always
be checked to determine if an error occurred while invoking the service operation. Errors returned can range from
client-side input validation errors to service-side error responses returned to the client. The operation's output
struct should not be accessed in the event that a non-nil error is returned by the client.

For example to log an operation error and prematurely return from the calling function:
```go
response, err := client.GetObject(context.TODO())
if err != nil {
	log.Printf("GetObject error: %v", err)
	return
}
```

For more information on error handling, including how to inspect for specific error types, see the
[Handling Errors]({{% ref "handling-errors.md" %}}) documentation.

#### Responses with io.ReadCloser

Some API operations return a response struct that contain an output member that is an `io.ReadCloser`. If you're making
requests with these operations, always be sure to call `io.ReadCloser` member's `Close` method after you've completed
reading the content.

For example {{% alias service=S3 %}} `GetObject` operation returns a response
whose `Body` member is an `io.ReadCloser`:

```go
resp, err := s3svc.GetObject(context.TODO(), &s3.GetObjectInput{...})
if err != nil {
    // handle error
    return
}
// Make sure to always close the response Body when finished
defer resp.Body.Close()

decoder := json.NewDecoder(resp.Body)
if err := decoder.Decode(&myStruct); err != nil {
    // handle error
    return
}
```

#### Response Metadata

All service operation output structs include a `ResultMetadata` member of type 
[middleware.Metadata]({{< apiref smithy="middleware#Metadata" >}}). `middleware.Metadata` is used by the SDK middleware
to provide additional information from a service response that is not modeled by the service. This includes metadata
like the `RequestID`. For example to retrieve the `RequestID` associated with a service response to assit AWS Support in
troubleshooting a request:

```go
import "fmt"
import "log"
import "github.com/aws/aws-sdk-go-v2/aws/middleware"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ..

resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
	// ...
})
if err != nil {
	log.Printf("error: %v", err)
	return
}

requestID, ok := middleware.GetRequestIDMetadata(resp.ResultMetadata)
if !ok {
	fmt.Println("RequestID not included with request")
}

fmt.Printf("RequestID: %s\n", requestID)
```

## Concurrently Using Service Clients

You can create goroutines that concurrently use the same service client to send multiple requests. You can use a service
client with as many goroutines as you want. 

In the following example, an {{% alias=service=S3 %}} service client is used in multiple goroutines. This example
concurrently uploads two objects to an {{% alias service=S3 %}} bucket.

```go
import "context"
import "log"
import "strings"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	log.Printf("error: %v", err)
	return
}

client := s3.NewFromConfig(cfg)

type result struct {
    Output *s3.PutObjectOutput
    Err    error
}

results := make(chan result, 2)

var wg sync.WaitGroup
wg.Add(2)

go func() {
defer wg.Done()
    output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String("my-bucket"),
        Key:    aws.String("foo"),
        Body:   strings.NewReader("foo body content"),
    })
    results <- result{Output: output, Err: err}
}()

go func() {
    defer wg.Done()
    output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String("my-bucket"),
        Key:    aws.String("bar"),
        Body:   strings.NewReader("bar body content"),
    })
    results <- result{Output: output, Err: err}
}()

wg.Wait()

close(results)

for result := range results {
    if result.Err != nil {
        log.Printf("error: %v", result.Err)
        continue
    }
    fmt.Printf("etag: %v", aws.ToString(result.Output.ETag))
}
```

## Using Operation Paginators {id="using-paginators"}

Typically, when you retrieve a list of items, you might need to check the output struct for a token or marker to confirm
whether the AWS service returned all results from your request. If the token or marker is present, you use it to request
the next page of results. Instead of managing these tokens or markers, you can use the service package's available
paginator types.

Paginator helpers are available for supported service operations, and can be found in the service client's Go package.
To construct a paginator for a supported operation, use the `New<OperationName>Paginator` function. Paginator construct
functions take the service `Client`, the operation's `<OperationName>Input` input parameters, and an optional set of
functional arguments allowing you to configure other optional paginator settings.

The returned operation paginator type provides a convenient way to iterate over a paginated operation until you have
reached the last page, or you have found the item(s) that your application was searching for. A paginator type has
two methods: `HasMorePages` and `NextPage`. `HasMorePages` returns a boolean value of `true` if the first page has not
been retrieved, or if additional pages available to retrieve using the operation. To retrieve the first or subsequent
pages of the operation, the `NextPage` operation must be called. `NextPage` takes `context.Context` and returns
the operation output and any corresponding error. Like the client operation method return parameters, the return error
should always be checked before attempting to use the returned response structure.
See [Handling Operation Responses]({{% ref "#HandlingOperationResponses" %}})

The following example uses the `ListObjectsV2` paginator to list up to three pages of object keys from the
`ListObjectV2`operation. Each page consists of up to 10 keys, which is defined by the `Limit` paginator option.

```go
import "context"
import "log"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	log.Printf("error: %v", err)
	return
}

client := s3.NewFromConfig(cfg)

params := &s3.ListObjectsV2Input{
	Bucket: aws.String("my-bucket"),
}

paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
	o.Limit = 10
})

pageNum := 0
for paginator.HasMorePages() && pageNum < 3 {
    output, err := paginator.NextPage(context.TODO())
    if err != nil {
    	log.Printf("error: %v", err)
    	return
    }
    for _, value := range output.Contents {
        fmt.Println(*value.Key)
    }
    pageNum++
}
```

Similar to client operation method, the client options like the request Region can be modified by providing one or more
functional arguments to `NextPage`. For more information about overriding client options when calling an operation see
[Overriding Clients For Operation]({{% ref "#OverrideClientOptionsForOperation" %}})


## Using Waiters

When interacting with AWS APIs that are asynchronous, you often need to wait 
for a particular resource to become available in order to perform further 
actions on it. 

For example, the {{% alias service=DDBlong %}} `CreateTable` API returns 
immediately with a TableStatus of CREATING, and you can't invoke read or 
write operations until the table status has been transitioned to `ACTIVE`. 

Writing logic to continuously poll the table status can be cumbersome 
and error-prone. The waiters help take the complexity out of it and 
are simple APIs that handle the polling task for you.

For example, you can use waiters to poll if a {{% alias service=DDB %}} table 
is created and ready for a write operation.

```go
import "context"
import "fmt"
import "log"
import "time"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    log.Printf("error: %v", err)
    return
}

client := dynamodb.NewFromConfig(cfg)

// we create a waiter instance by directly passing in a client
// that satisfies the waiters client Interface. 
waiter :=  dynamodb.NewTableExistsWaiter(client)

// params is the input to api operation used by the waiter
params := &dynamodb.DescribeTableInput {
	TableName: aws.String("test-table")
}

// maxWaitTime is the maximum wait time, the waiter will wait for 
// the resource status.
maxWaitTime := 5 * time.Minutes

// Wait will poll until it gets the resource status, or max wait time 
// expires.
err := waiter.Wait(context.TODO(), params, maxWaitTime)  
if err != nil {
    log.Printf("error: %v", err)
    return 
}
fmt.Println("Dynamodb table is now ready for write operations")

```

#### Overriding waiter configuration

By default, the SDK uses the minimum delay and maximum delay value configured with 
optimal values defined by AWS services for different APIs. You can override waiter 
configuration by providing functional options during waiter construction, or when 
invoking a waiter operation. 

For example, to override waiter configuration during waiter construction

```go
import "context"
import "fmt"
import "log"
import "time"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    log.Printf("error: %v", err)
    return
}

client := dynamodb.NewFromConfig(cfg)

// we create a waiter instance by directly passing in a client
// that satisfies the waiters client Interface. 
waiter :=  dynamodb.NewTableExistsWaiter(client, func (o *dynamodb.TableExistsWaiterOptions) {
	
	// override minimum delay to 10 seconds
	o.MinDelay = 10 * time.Second
	
	// override maximum default delay to 300 seconds
	o.MaxDelay = 300 * time.Second
})
```

The `Wait` function on each waiter also takes in functional options.  
Similar to the above example, you can override waiter configuration per `Wait` request. 

```go
// params is the input to api operation used by the waiter
params := &dynamodb.DescribeTableInput {
	TableName: aws.String("test-table")
}

// maxWaitTime is the maximum wait time, the waiter will wait for 
// the resource status.
maxWaitTime := 5 * time.Minutes

// Wait will poll until it gets the resource status, or max wait time 
// expires.
err := waiter.Wait(context.TODO(), params, maxWaitTime, func (o *dynamodb.TableExistsWaiterOptions) {

    // override minimum delay to 5 seconds
    o.MinDelay = 5 * time.Second

    // override maximum default delay to 120 seconds
    o.MaxDelay = 120 * time.Second
})
if err != nil {
    log.Printf("error: %v", err)
    return 
}
fmt.Println("Dynamodb table is now ready for write operations")

```

#### Advanced waiter configuration overrides

You can additionally customize the waiter default behavior by providing a custom 
retryable function. The waiter-specific options also provides `APIOptions` to
[customize operation middlewares](https://aws.github.io/aws-sdk-go-v2/docs/middleware/#writing-a-custom-middleware).

For example, to configure advanced waiter overrides.

```go
import "context"
import "fmt"
import "log"
import "time"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/dynamodb"
import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    log.Printf("error: %v", err)
    return
}

client := dynamodb.NewFromConfig(cfg)

// custom retryable defines if a waiter state is retryable or a terminal state.
// For example purposes, we will configure the waiter to not wait 
// if table status is returned as `UPDATING`
customRetryable := func(ctx context.Context, params *dynamodb.DescribeTableInput, 
	output *dynamodb.DescribeTableOutput, err error) (bool, error) {
	if output.Table != nil {
		if output.Table.TableStatus == types.TableStatusUpdating {
			// if table status is `UPDATING`, no need to wait
		    return false, nil	
        }
    }
}

// we create a waiter instance by directly passing in a client
// that satisfies the waiters client Interface. 
waiter :=  dynamodb.NewTableExistsWaiter(client, func (o *dynamodb.TableExistsWaiterOptions) {
	
	// override the service defined waiter-behavior
	o.Retryable = customRetryable
})

```
