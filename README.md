# AWS SDK for Go v2

![Build Status](https://codebuild.us-west-2.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoib1lGQ3N6RFJsalI5a3BPcXB3Rytaak9kYVh1ZW1lZExPNjgzaU9Udng3VE5OL1I3czIwcVhkMUlUeG91ajBVaWRYcVVJSEVQcmZwTWVyT1p5MGszbnA4PSIsIml2UGFyYW1ldGVyU3BlYyI6IkhrZ1VMN20zRmtYY1BrR0wiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master) [![API Reference](https://img.shields.io/badge/api-reference-blue.svg)](https://docs.aws.amazon.com/sdk-for-go/v2/api) [![Join the chat at https://gitter.im/aws/aws-sdk-go](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/aws/aws-sdk-go-v2?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/aws/aws-sdk-go/blob/master/LICENSE.txt)


`aws-sdk-go-v2` is the **Developer Preview** (aka **beta**) for the v2 AWS SDK for the Go programming language. This Developer Preview is provided to receive feedback from the language community on SDK changes prior to the final release. As such users should expect the SDK to release minor version releases that break backwards compatability. The release notes for the breaking change will include information about the breaking change, and how you can migrate to the latest version.

Check out the [Issues] and [Projects] for design and updates being made to the SDK. The v2 SDK requires a minimum version of `Go 1.13`.

We'll be expanding out the [Issues] and [Projects] sections with additional changes to the SDK based on your feedback, and SDK's core's improvements. Check the the SDK's [CHANGE_LOG] for information about the latest updates to the SDK.

Jump To:
* [Project Status](_#Project-Status_)
* [Getting Started](_#Getting-Started_)
* [Quick Examples](_#Quick-Examples_)
* [Getting Help](_#Getting-Help_)
* [Contributing](_#Feedback-and-contributing_)
* [More Resources](_#Resources_)

## Project Status
The SDK is in preview state as we work to design and implement potentially breaking changes to the SDK as we update the SDK's layout and usage patterns based on your feedback. You can also expect periodic service API model updates as well.

We are actively seeking community feedback for several changes to the SDK. Please review our [design] page on issues
that are currently pending community feedback.

Users should expect significant changes that could affect the following (non-exhaustive) areas:
* Package Locations
  * Includes Location of Service API Types
* Modularization
* Credential Providers
* Paginators
* Waiters
* Service Endpoint Resolution
* Minimum Supported Go Release following the [Language Release Policy](https://golang.org/doc/devel/release.html#policy)

## Getting started

The best way to get started working with the SDK is to use `go get` to add the SDK to your Go Workspace or application using Go modules.

```sh
go get github.com/aws/aws-sdk-go-v2
```

Without Go Modules, or in a GOPATH use the `/...` suffix on the `go get` to retrieve all of the SDK's dependencies.

```sh
go get github.com/aws/aws-sdk-go-v2/...
```

## Quick Examples

### Hello AWS

This example shows how you can use the v2 SDK to make an API request using the SDK's [Amazon DynamoDB] client.

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = "us-west-2"

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.New(cfg)

	// Build the request with its input parameters
	req := svc.DescribeTableRequest(&dynamodb.DescribeTableInput{
		TableName: aws.String("myTable"),
	})

	// Send the request, and get the response or error back
	resp, err := req.Send(context.Background())
	if err != nil {
		panic("failed to describe table, "+err.Error())
	}

	fmt.Println("Response", resp)
}
```

## Getting Help

Please use these community resources for getting help. We use the GitHub issues
for tracking bugs and feature requests.

* Ask a question on [StackOverflow](http://stackoverflow.com/) and tag it with the [`aws-sdk-go`](http://stackoverflow.com/questions/tagged/aws-sdk-go) tag.
* Come join the AWS SDK for Go community chat on [gitter](https://gitter.im/aws/aws-sdk-go).
* Open a support ticket with [AWS Support](http://docs.aws.amazon.com/awssupport/latest/user/getting-started.html).
* If you think you may have found a bug, please open an [issue](https://github.com/aws/aws-sdk-go-v2/issues/new/choose).

This SDK implements AWS service APIs. For general issues regarding the AWS services and their limitations, you may also take a look at the [Amazon Web Services Discussion Forums](https://forums.aws.amazon.com/).

### Opening Issues

If you encounter a bug with the AWS SDK for Go we would like to hear about it.
Search the [existing issues](https://github.com/aws/aws-sdk-go-v2/issues) and see
if others are also experiencing the issue before opening a new issue. Please
include the version of AWS SDK for Go, Go language, and OS you’re using. Please
also include reproduction case when appropriate.

The GitHub issues are intended for bug reports and feature requests. For help
and questions with using AWS SDK for Go please make use of the resources listed
in the [Getting Help](https://github.com/aws/aws-sdk-go-v2#getting-help) section.
Keeping the list of open issues lean will help us respond in a timely manner.

## Feedback and contributing

The v2 SDK will use GitHub [Issues] to track feature requests and issues with the SDK. In addition, we'll use GitHub [Projects] to track large tasks spanning multiple pull requests, such as refactoring the SDK's internal request lifecycle. You can provide feedback to us in several ways. 

**GitHub issues**. To provide feedback or report bugs, file GitHub [Issues] on the SDK. This is the preferred mechanism to give feedback so that other users can engage in the conversation, +1 issues, etc. Issues you open will be evaluated, and included in our roadmap for the GA launch.

**Gitter channel**. For more informal discussions or general feedback, check out our [Gitter channel] for the SDK. The [Gitter channel] is also a great place to ask general questions, and find help to get started with the 2.0 SDK Developer Preview.

**Contributing**. You can open pull requests for fixes or additions to the AWS SDK for Go 2.0 Developer Preview release. All pull requests must be submitted under the Apache 2.0 license and will be reviewed by an SDK team member before being merged in. Accompanying unit tests, where possible, are appreciated.

## Resources

[Developer guide](https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/welcome.html/) - This document
is a general introduction on how to configure and make requests with the SDK.
If this is your first time using the SDK, this documentation and the API
documentation will help you get started. This document focuses on the syntax
and behavior of the SDK. The [Service Developer Guide](https://aws.amazon.com/documentation/) 
will help you get started using specific AWS services.

[SDK API Reference Documentation](https://docs.aws.amazon.com/sdk-for-go/v2/api/) - Use this
document to look up all API operation input and output parameters for AWS
services supported by the SDK. The API reference also includes documentation of
the SDK, and examples how to using the SDK, service client API operations, and
API operation require parameters.

[Service Documentation](https://aws.amazon.com/documentation/) - Use this
documentation to learn how to interface with AWS services. These guides are
great for getting started with a service, or when looking for more 
information about a service. While this document is not required for coding, 
services may supply helpful samples to look out for.

[SDK Examples](https://github.com/aws/aws-sdk-go-v2/tree/master/example) -
Included in the SDK's repo are several hand crafted examples using the SDK
features and AWS services.

[Forum](https://forums.aws.amazon.com/forum.jspa?forumID=293) - Ask questions, get help, and give feedback

[Issues](https://github.com/aws/aws-sdk-go-v2/issues) - Report issues, submit pull requests, and get involved
  (see [Apache 2.0 License][sdk-license])

[Dep]: https://github.com/golang/dep
[Issues]: https://github.com/aws/aws-sdk-go-v2/issues
[Projects]: https://github.com/aws/aws-sdk-go-v2/projects
[CHANGE_LOG]: https://github.com/aws/aws-sdk-go-v2/blob/master/CHANGELOG.md
[Amazon DynamoDB]: https://aws.amazon.com/dynamodb/
[Gitter channel]: https://gitter.im/aws/aws-sdk-go-v2
[design]: https://github.com/aws/aws-sdk-go-v2/blob/master/DESIGN.md  