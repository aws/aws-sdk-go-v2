---
title: "Amazon EC2 Instance Metadata Service"
linkTitle: "Amazon EC2 IMDS"
description: "Using the AWS SDK for Go V2 Amazon EC2 Instance Metadata Service Client"
---

You can use the {{% alias sdk-go %}} to access the
[{{% alias service=EC2 %}} Instance Metadata Service](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).
The [feature/ec2/imds]({{% apiref "feature/ec2/imds" %}}) Go package provides a
[Client]({{% apiref "feature/ec2/imds#Client" %}}) type that can be used to access the {{% alias service=EC2 %}}
Instance Metadata Service. The `Client` and associated operations can be used similar to the other AWS service clients
provided by the SDK. To learn more information on how to configure the SDK, and use service clients see
[Configuring the SDK]({{% ref "configuring-sdk" %}}) and [Using AWS Services]({{% ref "making-requests.md" %}}).

The client can help you easily retrieve information about instances on which your applications run, such as its AWS
Region or local IP address. Typically, you must create and submit HTTP requests to retrieve instance metadata. Instead,
create an `imds.Client` to access the {{% alias service=EC2 %}} Instance Metadata Service using a programmatic client
like other AWS Services.

For example to construct a client:
```go
import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	log.Printf("error: %v", err)
	return
}

client := imds.NewFromConfig(cfg)
```

Then use the service client to retrieve information from a metadata category such as `local-ipv4`
(the private IP address of the instance).

```go
localip, err := client.GetMetadata(context.TODO(), &imds.GetMetadataInput{
	Path: "local-ipv4",
})
if err != nil {
    log.Printf("Unable to retrieve the private IP address from the EC2 instance: %s\n", err)
    return
}

fmt.Printf("local-ip: %v\n", localip)
```

For a list of all metadata categories, see
[Instance Metadata Categories](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html#dynamic-data-categories)
in the {{% alias service="EC2" %}} User Guide.

### Retrieving an Instance's Region

There's no instance metadata category that returns only the Region of an
instance. Instead, use the included `Region` method to easily return
an instance's Region.

```go
response, err := client.GetRegion(context.TODO(), &imds.GetRegionInput{})
if err != nil {
    log.Printf("Unable to retrieve the region from the EC2 instance %v\n", err)
}

fmt.Printf("region: %v\n", response.Region)
```

For more information about the EC2 metadata utility, see the [feature/ec2/imds]({{% apiref "feature/ec2/imds" %}}) package in the
{{% alias sdk-api %}}.

