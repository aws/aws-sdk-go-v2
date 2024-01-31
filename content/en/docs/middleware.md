---
title: "Customizing the AWS SDK for Go V2 Client Requests"
linkTitle: "Middleware"
description: "How to use stack step middleware to customize AWS SDK for Go V2 client requests."
weight: 8
---

{{% pageinfo color="warning" %}}
Modifying the client request pipeline can result in malformed/invalid requests, or
can result in unexpected application errors. This functionality is meant for advanced uses cases not provided by the SDK
interface by default.
{{% /pageinfo %}}

You can customize {{% alias sdk-go %}} client requests by registering one or more middleware to a service operation's
[stack]({{< apiref smithy="middleware#Stack" >}}). The stack is composed of a series of steps: Initialize, Serialize,
Build, Finalize, and Deserialize. Each step contains zero or more middleware that operate on that step's input and
output types. The following diagram and table provide an overview of how an operation's request and response traverses
the stack.

![Middleware](/aws-sdk-go-v2/diagrams/middleware.png)

Stack Step | Description
--- | ---
Initialize | Prepares the input, and sets any default parameters as needed.
Serialize | Serializes the input to a protocol format suitable for the target transport layer.
Build | Attach additional metadata to the serialized input, such as HTTP Content-Length.
Finalize | Final message preparation, including retries and authentication (SigV4 signing).
Deserialize | Deserialize responses from the protocol format into a structured type or error.

Each middleware within a given step must have a unique identifier, which is determined by the middleware's
`ID` method. Middleware identifiers ensure that only one instance of a given middleware is registered to a step, and
allows other step middleware to be inserted relative to it.

You attach step middleware by using a step's `Insert` or `Add` methods. You use `Add` to attach a middleware to the 
beginning of a step by specifying
[middleware.Before]({{< apiref smithy="middleware#Before" >}}) as the
[RelativePosition]({{< apiref smithy="middleware#RelativePosition" >}}), and
[middleware.After]({{< apiref smithy="middleware#After" >}}) to attach to the end of the step.
You use `Insert` to attach a middleware to a step by inserting the middleware relative to another step middleware.

{{% pageinfo color="warning" %}}
You must use the `Add` method to safely insert custom step middleware. Using `Insert` creates a dependency between your
custom middleware, and the middleware that you are inserting relative to. The middleware within a stack step must be
considered opaque to avoid breaking changes occurring to your application.
{{% /pageinfo %}}

## Writing a Custom Middleware

Each stack step has an interface that you must satisfy in order attach a middleware to a given step. You can use one of
the provided `<Step>MiddlewareFunc` functions to quickly satisfy this interface. The following table outlines the steps,
their interface, and the helper function that can be used to satisfy the interface.

Step | Interface | Helper Function
--- | --- | ---
Initialize | [InitializeMiddleware]({{< apiref smithy="middleware#InitializeMiddleware" >}}) | [InitializeMiddlewareFunc]({{< apiref smithy="middleware#InitializeMiddlewareFunc" >}})
Build | [BuildMiddleware]({{< apiref smithy="middleware#BuildMiddleware" >}}) | [BuildMiddlewareFunc]({{< apiref smithy="middleware#BuildMiddlewareFunc" >}})
Serialize | [SerializeMiddleware]({{< apiref smithy="middleware#SerializeMiddleware" >}}) | [SerializeMiddlewareFunc]({{< apiref smithy="middleware#SerializeMiddlewareFunc" >}})
Finalize | [FinalizeMiddleware]({{< apiref smithy="middleware#FinalizeMiddleware" >}}) | [FinalizeMiddlewareFunc]({{< apiref smithy="middleware#FinalizeMiddlewareFunc" >}})
Deserialize | [DeserializeMiddleware]({{< apiref smithy="middleware#DeserializeMiddleware" >}}) | [DeserializeMiddlewareFunc]({{< apiref smithy="middleware#DeserializeMiddlewareFunc" >}})

The following examples show how you can write a custom middleware to populate the Bucket member of the
{{% alias service=S3 %}} `GetObject` API calls if one is not provided. This middleware will be referenced in proceeding
examples to show how to attach step middleware to the stack.

```go
import "github.com/aws/smithy-go/aws"
import "github.com/aws/smithy-go/middleware"
import "github.com/aws/aws-sdk-go-v2/service/s3"

// ...

var defaultBucket = middleware.InitializeMiddlewareFunc("DefaultBucket", func(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	// Type switch to check if the input is s3.GetObjectInput, if so and the bucket is not set, populate it with
	// our default.
	switch v := in.Parameters.(type) {
	case *s3.GetObjectInput:
		if v.Bucket == nil {
			v.Bucket = aws.String("my-default-bucket")
		}
	}

	// Middleware must call the next middleware to be executed in order to continue execution of the stack.
	// If an error occurs, you can return to prevent further execution.
	return next.HandleInitialize(ctx, in)
})

```

## Attaching Middleware to All Clients

You can attach your custom step middleware to every client by adding the middleware using the `APIOptions` member of the
[aws.Config]({{< apiref "aws#Config" >}}) type. The following examples attaches the `defaultBucket` middleware to every
client constructed using your applications `aws.Config` object:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/smithy-go/middleware"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	// handle error
}

cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
	// Attach the custom middleware to the beginning of the Initialize step
	return stack.Initialize.Add(defaultBucket, middleware.Before)
})

client := s3.NewFromConfig(cfg)
```

## Attaching Middleware to a Specific Operation

You can attach your custom step middleware to a specific client operation by modifying the client's `APIOptions`
member using the variadic argument list for an operation. The following examples attaches the `defaultBucket` middleware
to a specific {{% alias service=S3 %}} `GetObject` operation invocation:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/smithy-go/middleware"

// ...

// registerDefaultBucketMiddleware registers the defaultBucket middleware with the provided stack.
func registerDefaultBucketMiddleware(stack *middleware.Stack) error {
	// Attach the custom middleware to the beginning of the Initialize step
	return stack.Initialize.Add(defaultBucket, middleware.Before)
}

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	// handle error
}

client := s3.NewFromConfig(cfg)

object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
	Key: aws.String("my-key"),
}, func(options *s3.Options) {
	// Register the defaultBucketMiddleware for this operation only
	options.APIOptions = append(options.APIOptions, registerDefaultBucketMiddleware)
})
```

## Passing Metadata Down the Stack

In certain situations, you may find that you require two or more middleware to function in tandem by sharing information
or state. You can use [context.Context](https://golang.org/pkg/context/#Context) to pass this metadata by using 
[middleware.WithStackValue]({{< apiref smithy="middleware#WithStackValue" >}}).
`middleware.WithStackValue` attaches the given key-value pair to the provided context, and safely limits the scope to
the currently executing stack. These stack-scoped values can be retrieved from a context using
[middleware.GetStackValue]({{< apiref smithy="middleware#GetStackValue" >}}) and
providing the key used to stored the corresponding value. Keys must be comparable, and you must define your own
types as context keys to avoid collisions. The following examples shows how two middleware can use `context.Context` to
pass information down the stack.

```go
import "context"
import "github.com/aws/smithy-go/middleware"

// ...

type customKey struct {}

func GetCustomKey(ctx context.Context) (v string) {
	v, _ = middleware.GetStackValue(ctx, customKey{}).(string)
	return v
}

func SetCustomKey(ctx context.Context, value string) context.Context {
	return middleware.WithStackValue(ctx, customKey{}, value)
}

// ...

var customInitalize = middleware.InitializeMiddlewareFunc("customInitialize", func(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	ctx = SetCustomKey(ctx, "my-custom-value")
	
	return next.HandleInitialize(ctx, in)
})

var customBuild = middleware.BuildMiddlewareFunc("customBuild", func(
	ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler,
) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	customValue := GetCustomKey(ctx)
	
	// use customValue

	return next.HandleBuild(ctx, in)
})
```

### Metadata Provided by the SDK

The {{% alias sdk-go %}} provides several metadata values that can be retrieved from the provided context. These values
can be used to enable more dynamic middleware that modifies its behavior based on the executing service, operation, or
target region. A few of the available keys are provided in the table below:

Key | Retriever | Description
--- | --- | ---
ServiceID | [GetServiceID]({{< apiref "aws/middleware#GetServiceID" >}}) | Retrieve the service identifier for the executing stack. This can be compared to the service client package's `ServiceID` constant.
OperationName | [GetOperationName]({{< apiref "aws/middleware#GetOperationName" >}}) | Retrieve the operation name for the executing stack.
Logger | [GetLogger]({{< apiref smithy="middleware#GetLogger" >}}) | Retrieve the logger that can be used for logging message from the middleware.


## Passing Metadata Up the Stack

You can pass metadata up through the stack by adding metadata key and value pairs using the
[middleware.Metadata]({{< apiref smithy="middleware#Metadata" >}}). Each middleware step returns an output structure,
metadata, and an error. Your custom middleware must return the metadata received from calling the next handler in the
step. This ensures that metadata added by downstream middleware propagates to the application invoking the
service operation. The resulting metadata is accessible to the invoking application by either the operation's
output shape via the `ResultMetadata` structure member.

The following examples shows how a custom middleware can add metadata that is returned as part of the operation output.

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/service/s3"
import "github.com/aws/smithy-go/middleware"

// ...

type customKey struct{}

func GetCustomKey(metadata middleware.Metadata) (v string) {
	v, _ = metadata.Get(customKey{}).(string)
	return v
}

func SetCustomKey(metadata *middleware.Metadata, value string) {
	metadata.Set(customKey{}, value)
}

// ...

var customInitalize = middleware.InitializeMiddlewareFunc("customInitialize", func (
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleInitialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	
	SetCustomKey(&metadata, "my-custom-value")
	
	return out, metadata, nil
})

// ...

client := s3.NewFromConfig(cfg, func (options *s3.Options) {
	options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
		return stack.Initialize.Add(customInitalize, middleware.After)
	})
})

out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
	// input parameters
})
if err != nil {
	// handle error
}

customValue := GetCustomKey(out.ResponseMetadata)

```
