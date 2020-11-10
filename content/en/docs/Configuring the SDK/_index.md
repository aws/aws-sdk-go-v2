---
title: "Configuring the AWS SDK for Go V2"
linkTitle: "Configuring the SDK"
weight: 3
---

In the AWS SDK for Go V2, you can configure common settings for service clients, such as the logger, log level, and
retry configuration. Most settings are optional. However, for each service client, you must specify an AWS Region and
your credentials. The SDK uses these values to send requests to the correct Region and sign requests with the correct
credentials. You can specify these values as programmatically in code, or via the execution environment.

## Loading AWS Shared Configuration

There are a number of ways to initialize a service API client, but the following is the most common pattern recommended
to users.

To configure the SDK to use the AWS shared configuration use the following code:

```go
import (
  "log"
  "github.com/aws/aws-sdk-go-v2/config"
)

// ...

cfg, err := config.LoadDefaultConfig()
if err != nil {
  log.Fatalf("failed to load configuration, %v", err)
}
```

`config.LoadDefaultConfig()` will construct an [aws.Config](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config)
using the AWS shared configuration sources. This includes configuring a credential provider. configuring the region, and
loading service specific configuration. Service clients can be constructed using the loaded `aws.Config`, providing a
consistent pattern for constructing clients in a uniform manner.

For more information about AWS Shared Configuration see the
[AWS Tools and SDKs Shared Configuration and Credentials Reference Guide ](https://docs.aws.amazon.com/credref/latest/refdocs/overview.html)

## Specifying the AWS Region

When you specify the Region, you specify where to send requests, such as us-west-2 or us-east-2. For a list of Regions
for each service, see Regions and Endpoints in the Amazon Web Services General Reference.

The SDK does not have a default Region. To specify a Region:

* Set the `AWS_REGION` environment variable to the default Region

* Set the region explicitly
  using [config.WithRegion](https://github.com/aws/aws-sdk-go-v2/blob/config/v0.2.2/config/provider.go#L127)
  as an argument to `config.LoadDefaultConfig` when loading configuration.

If you set a Region using all of these techniques, the SDK uses the Region you explicitly specified.

##### Configure Region with Environment Variable

###### Linux, macOS, or Unix

```bash
export AWS_REGION=us-west-2
```

###### Windows

```batchfile
set AWS_REGION=us-west-2
```

##### Specify Region Explicitly

```go
cfg, err := config.LoadDefaultConfig(config.WithRegion("us-west-2"))
```

## Specifying Credentials

The AWS SDK for Go requires credentials (an access key and secret access key) to sign requests to AWS. You can specify
your credentials in several locations, depending on your particular use case.

When you initialize a new service client without providing any credential arguments, the SDK uses the default credential
provider chain to find AWS credentials. The SDK uses the first provider in the chain that returns credentials without an
error. The default provider chain looks for credentials in the following order: 


