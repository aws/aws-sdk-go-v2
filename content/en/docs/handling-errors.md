---
title: "Handling Errors in the AWS SDK for Go V2"
linkTitle: "Handling Errors"
date: "2020-11-12"
description: "Use the Error interface to handle errors from the AWS SDK for Go V2 or AWS service."
---

The {{% alias sdk-go %}} returns errors that satisfy the Go `error` interface type  You can use the `Error()` method to
get a formatted string of the SDK error message without any special handling. Errors returned by the SDK may implement
an `Unwrap` method. The `Unwrap` method is used by the SDK to provide additional contextual information to errors, while
providing access to the underlying error or chain of errors. The `Unwrap` method should be used with the
[errors.As](https://golang.org/pkg/errors#As) to handle unwrapping error chains.

It is important that your application check whether an error occurred after invoking a function or method that
can return an `error` interface type. The most basic form of error handling looks similar to the following example:

```go
if err != nil {
	// Handle error
	return
}
```

## Logging Errors

The simplest form of error handling is traditionally to log or print the error message before returning or exiting from
the application. For example:

```go
import "log"

// ...

if err != nil {
	log.Printf("error: %s", err.Error())
	return
}
```

## Service Client Errors

The SDK wraps All errors returned by service clients with the
[smithy.OperationError]({{% apiref smithy="#OperationError" %}}) error type. `OperationError` provides contextual
information about the service name and operation that is associated with an underlying error. This information can be
useful for applications that perform batches of operations to one or more services, with a centralized error handling
mechanism. Your application can use `errors.As` to access this `OperationError` metadata.

For example:

```go
import "log"
import "github.com/awslabs/smithy-go"

// ...

if err != nil {
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
    }
    return
}
```

### API Error Responses

Service operations can return modeled error types to indicate specific errors. These modeled types can be used with
`errors.As` to unwrap and determine if the operation failure was due to a specific error. For example
{{% alias service=S3 %}} `CreateBucket` can return a
[BucketAlreadyExists]({{< apiref "service/s3/types#BucketAlreadyExists" >}}) error if a bucket of the same name
already exists.

For example, to check if an error was a `BucketAlreadyExists` error:

```go
import "log"
import "github.com/aws/aws-sdk-go-v2/service/s3/types"

// ...

if err != nil {
	var bne *types.BucketAlreadyExists
	if errors.As(err, &bne) {
		log.Println("error:", bne)
    }
    return
}
```

All service API response errors implement the [smithy.APIError]({{< apiref smithy="#APIError" >}}) interface type.
This interface can be used to handle both modeled or un-modeled service error responses. This type provides
access to the error code and message returned by the service. Additionally, this type provides indication of whether
the fault of the error was due to the client or server if known. For example:

```go
import "log"
import "github.com/awslabs/smithy-go"

// ...

if err != nil {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		log.Printf("code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
    }
    return
}
```

## Retrieving Request Identifiers

When working with AWS Support, you may be asked to provide the request identifier that identifiers the request you
are attempting to troubleshoot. You can use [http.ResponseError]({{< apiref "aws/transport/http#ResponseError" >}})
and use the `ServiceRequestID()` method to retrieve the request identifier associated with error response.

For example:

```go
import "log"
import awshttp "github.com/aws/transport/http"

// ...

if err != nil {
	var re *awshttp.ResponseError
	if errors.As(err, &re) {
		log.Printf("requestID: %s, error: %v", re.ServiceRequestID(), re.Unwrap());
    }
    return
}
```

### {{% alias service=S3 %}} Request Identifiers

{{% alias service=S3 %}} requests contain additional identifiers that can be used to assist AWS Support with
troubleshooting your request. {{% alias service=S3 %}} requests can contain a `RequestId` and `HostId` pair that can be
retrieved from {{% alias service=S3 %}} operation errors if available.

For example:

```go
import "log"

// ...

if err != nil {
	var re interface {
		ServiceHostID()    string
		ServiceRequestID() string
	}
	if errors.As(err, &re) {
		log.Printf("requestID: %s, hostID: %s request failure", re.ServiceRequestID(), re.ServiceHostID());
	}
	return
}
```

