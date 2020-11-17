---
title: "Using the AWS SDK for Go V2 with AWS Services"
linkTitle: "Using AWS Services"
date: "2020-11-12"
description: "Construct service clients and make operation calls to send requests to AWS services."
---

To make calls to an AWS service, you must first construct a service
client instance with a session. A service client provides low-level
access to every API action for that service. For example, you create an
{{% alias service=S3 %}} service client to make calls to {{% alias service=S3 %}}.

When you call service operations, you pass in input parameters as a
struct. A successful call usually results in an output struct that you
can use. For example, after you successfully call an {{% alias service=S3 %}} create bucket
action, the action returns an output struct with the bucket's location.

For the list of service clients, including their methods and parameters,
see the [{{% alias sdk-api %}}]({{< apiref >}})

## Constructing a Service

To construct a service client instance, use the `NewSession()`
function. The following example creates an {{% alias service=S3 %}} service client.
```go
sess, err := session.NewSession()
if err != nil {
    fmt.Println("Error creating session ", err)
    return
}
svc := s3.New(sess)
```

After you have a service client instance, you can use it to call service
operations. For more information about configurations, see
[Configuring the SDK]({{< relref "configuring-sdk" >}}).

When you create a service client, you can pass in custom configurations
so that you don't need to create a session for each configuration. The
SDK merges the two configurations, overriding session values with your
custom configuration. For example, in the following snippet, the |S3|
client uses the `mySession` session but overrides the `Region` field
with a custom value (``us-west-2``):

```go
svc := s3.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
```

## Tagging Service Resources
You can tag service resources,
such as {{% alias service=S3 %}} buckets,
so that you can determine the costs of your service resources
at whatever level of granularity you require.

Getting the HTTP Request and Response with Each Service Call
============================================================

You can direct the |sdk-go| to display the HTTP request
and response it sends and receives for each call by
including a configuration option when constructing the service client.

The following example uses the |DDB| **ListTables** operation
to illustrate how to add a custom header to a service call.

```go
package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/request"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"

    "fmt"
    "os"
)

func main() {
    // Initialize a session in us-west-2 that the SDK will use to load credentials
    // from the shared config file. (~/.aws/credentials).
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},
    )
    if err != nil {
        fmt.Println("Error getting session:")
        fmt.Println(err)
        os.Exit(1)
    }

    // Create DynamoDB client
    // and expose HTTP requests/responses
    svc := dynamodb.New(sess, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))

    // Add "CustomHeader" header with value of 10
    svc.Handlers.Send.PushFront(func(r *request.Request) {
        r.HTTPRequest.Header.Set("CustomHeader", fmt.Sprintf("%d", 10))
    })

    // Call ListTables just to see HTTP request/response
    // The request should have the CustomHeader set to 10
    _, _ = svc.ListTables(&dynamodb.ListTablesInput{})
}
```

If you run this program, the output should be similar to the following,
where **ACCESS-KEY** is the access key of the user and
**TABLE-1**, through **TABLE-N** are the names of the tables.

```
2017/10/25 11:10:57 DEBUG: Request dynamodb/ListTables Details:
---[ REQUEST POST-SIGN ]-----------------------------
POST / HTTP/1.1
Host: dynamodb.us-west-2.amazonaws.com
User-Agent: aws-sdk-go/1.10.34 (go1.8; windows; amd64)
Content-Length: 2
Accept-Encoding: identity
Authorization: AWS4-HMAC-SHA256 Credential=ACCESS-KEY/20171025/us-west-2/dynamodb/aws4_request, SignedHeaders=accept-encoding;content-length;content-type;host;x-amz-date;x-amz-target, Signature=9c92efe5d6c597cf29e4f7cc74de6dc2e39f8010a0d4957a397c59ef9cde21f2
Content-Type: application/x-amz-json-1.0
CustomHeader: 10
X-Amz-Date: 20171025T181057Z
X-Amz-Target: DynamoDB_20120810.ListTables

{}
-----------------------------------------------------
2017/10/25 11:10:58 DEBUG: Response dynamodb/ListTables Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK
Content-Length: 177
Connection: keep-alive
Content-Type: application/x-amz-json-1.0
Date: Wed, 25 Oct 2017 18:10:58 GMT
Server: Server
X-Amz-Crc32: 3023160996
X-Amzn-Requestid: M5B4BM4UU569MVBSDG5O2O9ITJVV4KQNSO5AEMVJF66Q9ASUAAJG


-----------------------------------------------------
2017/10/25 11:10:58 {"TableNames":["TABLE-1","...","TABLE-N"]}
```

## Service Operation Calls

You can call a service operation directly or with its request form. When
you call a service operation, the SDK synchronously validates the input,
builds the request, signs it with your credentials, sends it to AWS, and
then gets a response or an error. In most cases, you can call service
operations directly.

### Calling Operations

Calling the operation will sync as the request is built, signed, sent,
and the response is received. If an error occurs during the operation,
it will be returned. The output or resulting structure won't be valid.

For example, to call the {{% alias service=S3 %}} GET Object API, use the {{% alias service=S3 %}}
service client instance and call its [GetObject]({{< apiref "service/s3/#Client.GetObject" >}})
method:

```go
result, err := s3Svc.GetObject(&s3.GetObjectInput{...})
// result is a *s3.GetObjectOutput struct pointer
// err is a error which can be cast to awserr.Error.
```

### Passing Parameters to a Service Operation

When calling an operation on a service, you pass in input parameters as
option values, similar to passing in a configuration. For example, to
retrieve an object, you must specify a bucket and the object's key by
passing in the following parameters to the `GetObject` method:

```go
svc := s3.New(session.New())
svc.GetObject(&s3.GetObjectInput{
    Bucket: aws.String("bucketName"),
    Key:    aws.String("keyName"),
})
```

Each service operation has an associated input struct and, usually, an
output struct. The structs follow the naming pattern
*OperationName*`Input` and *OperationName*`Output`.

For more information about the parameters of each method, see the
service client documentation in the.

### Handling Operation Response Body

Some API operations return a response struct that contain a `Body`
field that is an `io.ReadCloser`. If you're making requests with
these operations, always be sure to call :code:`Close` on the field.

```go
resp, err := s3svc.GetObject(&s3.GetObjectInput{...})
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
    
## Concurrently Using Service Clients

You can create goroutines that concurrently use the same service client
to send multiple requests. You can use a service client with as many
goroutines as you want. However, you cannot concurrently modify the
service client's configuration and request handlers. If you do, the
service client operations might encounter race conditions. Define
service client settings before you concurrently use it.

In the following example, an {{% alias=service=S3 %}} service client is used in multiple
goroutines. The example concurrently outputs all objects in `bucket1`,
`bucket2`, and `bucket3`, which are all in the same region. To make
sure all objects from the same bucket are printed together, the example
uses a channel.

```go
sess, err := session.NewSession()
if err != nil {
    fmt.Println("Error creating session ", err)
}
var wg sync.WaitGroup
keysCh := make(chan string, 10)

svc := s3.New(sess)
buckets := []string{"bucket1", "bucket2", "bucket3"}
for _, bucket := range buckets {
    params := &s3.ListObjectsInput{
        Bucket:  aws.String(bucket),
        MaxKeys: aws.Int64(100),
    }
    wg.Add(1)
    go func(param *s3.ListObjectsInput) {
        defer wg.Done()

        err = svc.ListObjectsPages(params,
            func(page *s3.ListObjectsOutput, last bool) bool {
                // Add the objects to the channel for each page
                for _, object := range page.Contents {
                    keysCh <- fmt.Sprintf("%s:%s", *params.Bucket, *object.Key)
                }
                return true
            },
        )
        if err != nil {
            fmt.Println("Error listing", *params.Bucket, "objects:", err)
        }
    }(params)
}
go func() {
    wg.Wait()
    close(keysCh)
}()
for key := range keysCh {
    // Print out each object key as its discovered
    fmt.Println(key)
}
```

## Using Pagination Methods

Typically, when you retrieve a list of items, you might need to check
the output for a token or marker to confirm whether AWS returned all
results from your request. If present, you use the token or marker to
request the next set of results. Instead of managing these tokens or
markers, you can use pagination methods provided by the SDK.

Pagination methods iterate over a list operation until the method
retrieves the last page of results or until the callback function
returns `false`. The names of these methods use the following pattern:
*OperationName*\ `Pages`. For example, the pagination method for the
|S3| list objects operation (`ListObjects`) is `ListObjectPages`.

The following example uses the `ListObjectPages` pagination method to
list up to three pages of object keys from the `ListObject`
operation. Each page consists of up to 10 keys, which is defined by
the `MaxKeys` field.

```go
svc, err := s3.NewSession(sess)
if err != nil {
    fmt.Println("Error creating session ", err)
}
inputparams := &s3.ListObjectsInput{
    Bucket:  aws.String("mybucket"),
    MaxKeys: aws.Int64(10),
}
pageNum := 0
svc.ListObjectsPages(inputparams, func(page *s3.ListObjectsOutput, lastPage bool) bool {
    pageNum++
    for _, value := range page.Contents {
        fmt.Println(*value.Key)
    }
    return pageNum < 3
})
```

## Using Waiters

The SDK provides waiters that continuously check for completion of a
job. For example, when you send a request to create an {{% alias="S3" %}} bucket, you
can use a waiter to check when the bucket has been successfully created.
That way, subsequent operations on the bucket are done only after the
bucket has been created.

The following example uses a waiter that waits until specific instances
have stopped.

```go
sess, err := session.NewSession(aws.NewConfig().WithRegion("us-west-2"))
if err != nil {
    fmt.Println("Error creating session ", err)
}
// Create an EC2 client
ec2client := ec2.New(sess)
// Specify two instances to stop
instanceIDsToStop := aws.StringSlice([]string{"i-12345678", "i-23456789"})
// Send request to stop instances
_, err = ec2client.StopInstances(&ec2.StopInstancesInput{
  InstanceIds: instanceIDsToStop,
})
if err != nil {
  panic(err)
}
// Use a waiter function to wait until the instances are stopped
describeInstancesInput := &ec2.DescribeInstancesInput{
  InstanceIds: instanceIDsToStop,
}
if err := ec2client.WaitUntilInstanceStopped(describeInstancesInput); err != nil {
  panic(err)
}
fmt.Println("Instances are stopped.")
```
