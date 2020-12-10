# AWS SDK for Go v2

[![Build Status](https://travis-ci.org/aws/aws-sdk-go-v2.svg?branch=master)](https://travis-ci.org/aws/aws-sdk-go-v2) [![SDK Documentation](https://img.shields.io/badge/SDK-Documentation-blue)](https://aws.github.io/aws-sdk-go-v2/docs/) [![API Reference](https://img.shields.io/badge/api-reference-blue.svg)](https://pkg.go.dev/mod/github.com/aws/aws-sdk-go-v2) [![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/aws/aws-sdk-go/blob/master/LICENSE.txt)


`aws-sdk-go-v2` is the **Developer Preview** (aka **beta**) for the v2 AWS SDK for the Go programming language. This Developer Preview is provided to receive feedback from the language community on SDK changes prior to the final release. As such users should expect the SDK to release minor version releases that break backwards compatability. The release notes for the breaking change will include information about the breaking change, and how you can migrate to the latest version.

Check out the [Issues] and [Projects] for design and updates being made to the SDK. The v2 SDK requires a minimum version of `Go 1.15`.

We'll be expanding out the [Issues] and [Projects] sections with additional changes to the SDK based on your feedback, and SDK's core's improvements. Check the the SDK's [CHANGELOG] for information about the latest updates to the SDK.

Jump To:
* [Project Status](#project-status)
* [Getting Started](#getting-started)
* [Getting Help](#getting-help)
* [Contributing](#feedback-and-contributing)
* [More Resources](#resources)

## Project Status
The service API clients, and core features of the SDK have been significantly rewritten as part of the `v0.25.0` release. API clients are now independently versioned Go modules that are bundled with their respective API types and endpoints. Client API serialization and deserialization are now fully code-generated from the service model, providing performance improvements over previous releases.

A number of API clients may be immediately available or usable without workarounds. Additional support for these clients will be added over the next series of releases as support is extended to support their feature set.

We are actively seeking community feedback for API clients, and other upcoming design changes to the SDK. Please review our [design] page on issues
that are currently pending community feedback.

Users should expect significant changes that could affect the following (non-exhaustive) areas:
* Package Locations & Modularization
* Credential Providers
* Paginators
* Waiters
* Minimum Supported Go version following the [AWS SDKs and Tools Maintenance Policy](https://docs.aws.amazon.com/credref/latest/refdocs/maint-policy.html)

## Getting started
To get started working with the SDK setup your project for Go modules, and retrieve the SDK dependencies with `go get`.
This example shows how you can use the v2 SDK to make an API request using the SDK's [Amazon DynamoDB] client.

###### Initialize Project
```sh
$ mkdir ~/helloaws
$ cd ~/helloaws
$ go mod init helloaws
```
###### Add SDK Dependencies
```sh
$ go get github.com/aws/aws-sdk-go-v2/aws
$ go get github.com/aws/aws-sdk-go-v2/config
$ go get github.com/aws/aws-sdk-go-v2/service/dynamodb
```

###### Write Code
In your preferred editor add the following content to `main.go`

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
    // Using the SDK's default configuration, loading additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Using the Config value, create the DynamoDB client
    svc := dynamodb.NewFromConfig(cfg)

    // Build the request with its input parameters
    resp, err := svc.ListTables(context.Background(), &dynamodb.ListTablesInput{
        Limit: aws.Int32(5),
    })
    if err != nil {
        log.Fatalf("failed to list tables, %v", err)
    }

    fmt.Println("Tables:")
    for _, tableName := range resp.TableNames {
        fmt.Println(aws.ToString(tableName))
    }
}
```

###### Compile and Execute
```sh
$ go run .
Table:
tableOne
tableTwo
```

## Getting Help

Please use these community resources for getting help. We use the GitHub issues
for tracking bugs and feature requests.

* Ask a question on [StackOverflow](http://stackoverflow.com/) and tag it with the [`aws-sdk-go`](http://stackoverflow.com/questions/tagged/aws-sdk-go) tag.
* Come join the AWS SDK for Go community chat on [gitter][Gitter channel].
* Open a support ticket with [AWS Support](http://docs.aws.amazon.com/awssupport/latest/user/getting-started.html).
* If you think you may have found a bug, please open an [issue](https://github.com/aws/aws-sdk-go-v2/issues/new/choose).

This SDK implements AWS service APIs. For general issues regarding the AWS services and their limitations, you may also take a look at the [Amazon Web Services Discussion Forums](https://forums.aws.amazon.com/).

### Opening Issues

If you encounter a bug with the AWS SDK for Go we would like to hear about it.
Search the [existing issues][Issues] and see
if others are also experiencing the issue before opening a new issue. Please
include the version of AWS SDK for Go, Go language, and OS you’re using. Please
also include reproduction case when appropriate.

The GitHub issues are intended for bug reports and feature requests. For help
and questions with using AWS SDK for Go please make use of the resources listed
in the [Getting Help](#getting-help) section.
Keeping the list of open issues lean will help us respond in a timely manner.

## Feedback and contributing

The v2 SDK will use GitHub [Issues] to track feature requests and issues with the SDK. In addition, we'll use GitHub [Projects] to track large tasks spanning multiple pull requests, such as refactoring the SDK's internal request lifecycle. You can provide feedback to us in several ways. 

**GitHub issues**. To provide feedback or report bugs, file GitHub [Issues] on the SDK. This is the preferred mechanism to give feedback so that other users can engage in the conversation, +1 issues, etc. Issues you open will be evaluated, and included in our roadmap for the GA launch.

**Gitter channel**. For more informal discussions or general feedback, check out our [Gitter channel] for the SDK. The [Gitter channel] is also a great place to ask general questions, and find help to get started with the 2.0 SDK Developer Preview.

**Contributing**. You can open pull requests for fixes or additions to the AWS SDK for Go 2.0 Developer Preview release. All pull requests must be submitted under the Apache 2.0 license and will be reviewed by an SDK team member before being merged in. Accompanying unit tests, where possible, are appreciated.

## Resources

[SDK API Reference Documentation](https://pkg.go.dev/mod/github.com/aws/aws-sdk-go-v2) - Use this
document to look up all API operation input and output parameters for AWS
services supported by the SDK. The API reference also includes documentation of
the SDK, and examples how to using the SDK, service client API operations, and
API operation require parameters.

[Service Documentation](https://aws.amazon.com/documentation/) - Use this
documentation to learn how to interface with AWS services. These guides are
great for getting started with a service, or when looking for more 
information about a service. While this document is not required for coding, 
services may supply helpful samples to look out for.

[Forum](https://forums.aws.amazon.com/forum.jspa?forumID=293) - Ask questions, get help, and give feedback

[Issues] - Report issues, submit pull requests, and get involved
  (see [Apache 2.0 License][license])

[Dep]: https://github.com/golang/dep
[Issues]: https://github.com/aws/aws-sdk-go-v2/issues
[Projects]: https://github.com/aws/aws-sdk-go-v2/projects
[CHANGELOG]: https://github.com/aws/aws-sdk-go-v2/blob/master/CHANGELOG.md
[Amazon DynamoDB]: https://aws.amazon.com/dynamodb/
[Gitter channel]: https://gitter.im/aws/aws-sdk-go-v2
[design]: https://github.com/aws/aws-sdk-go-v2/blob/master/DESIGN.md  
[license]: http://aws.amazon.com/apache2.0/
