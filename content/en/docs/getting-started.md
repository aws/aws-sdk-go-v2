---
title: "Getting Started with the AWS SDK for Go V2"
linkTitle: "Getting Started"
weight: 2
---

The AWS SDK for Go requires Go 1.15 or later. You can view your current version of Go by running the following command.

```bash
go version
```

For information about installing or upgrading your version of Go, see https://golang.org/doc/install.

## Get an Amazon Account

Before you can use the AWS SDK for Go V2, you must have an Amazon account.
See [How do I create and activate a new AWS account?](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
for details.

## Install the AWS SDK for Go V2

The AWS SDK for Go V2 uses Go Modules, which was a feature introduced in Go 1.11. To get started initialize your local
project by running the following Go command.

```bash
go mod init example
```

After initializing your Go Module project you will be able to retrieve the SDK, and it's required dependencies using
the `go get` command. These dependencies will be recorded in the `go.mod` file which was created by the previous
command.

The following commands show how to retrieve the standard set of SDK modules that you will use in your application.

```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
```

This will retrieve the core SDK module, and the config module which is used for loading the AWS shared configuration.

Next you can install one or more AWS service API clients required by your application. All API clients are located
under `github.com/aws/aws-sdk-go-v2/service` import hierarchy. A complete set of currently supported API clients can be
found
[here](https://pkg.go.dev/mod/github.com/aws/aws-sdk-go-v2?tab=packages). To install a service client, execute the
following command to retrieve the module and record the dependency in your `go.mod` file. In this example we retrieve
the Amazon S3 API client.

```bash
go get github.com/aws/aws-sdk-go-v2/service/s3
```

## Get your AWS access keys

Access keys consist of an access key ID and secret access key, which are used to sign programmatic requests that you
make to AWS. If you don’t have access keys, you can create them by using
the [AWS Management Console](https://console.aws.amazon.com/console/home). We recommend that you use IAM access keys
instead of AWS root account access keys. IAM lets you securely control access to AWS services and resources in your AWS
account. Note

{{% pageinfo color="info" %}} To create access keys, you must have permissions to perform the required IAM actions. For
more information,
see [Granting IAM User Permission to Manage Password Policy and Credentials](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_delegate-permissions.html)
in the IAM User Guide. {{% /pageinfo %}}

### To get your access key ID and secret access key.

1. Open the [IAM console](https://console.aws.amazon.com/iam/home)
1. On the navigation menu, choose **Users**.
1. Choose your IAM user name (not the check box).
1. Open the **Security credentials** tab, and then choose **Create access key**.
1. To see the new access key, choose **Show**. Your credentials resemble the following:
    * Access key ID: `AKIAIOSFODNN7EXAMPLE`
    * Secret access key: `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`
1. To download the key pair, choose **Download .csv file**. Store the keys in a secure location.

{{% pageinfo color="warning" %}} Keep the keys confidential to protect your AWS account, and never email them. Do not
share them outside your organization, even if an inquiry appears to come from AWS or Amazon.com. No one who legitimately
represents Amazon will ever ask you for your secret key. {{% /pageinfo %}}

### Related topics

* [What Is IAM?](https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html)
  in IAM User Guide.
* [AWS Security Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html)
  in Amazon Web Services General Reference.

## Import Packages

After you have installed the SDK, you import AWS packages into your Go applications to use the SDK, as shown in the
following example, which imports the AWS, Config, and Amazon S3 libraries:

```go
import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)
```

