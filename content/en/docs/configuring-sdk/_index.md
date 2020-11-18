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

`config.LoadDefaultConfig()` will construct an [aws.Config]({{< apiref "aws#Config" >}})
using the AWS shared configuration sources. This includes configuring a credential provider. configuring the AWS Region,
and loading service specific configuration. Service clients can be constructed using the loaded `aws.Config`, providing
a consistent pattern for constructing clients.

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

```
export AWS_REGION=us-west-2
```

###### Windows

```batchfile
set AWS_REGION=us-west-2
```

##### Specify Region Programmatically

```go
cfg, err := config.LoadDefaultConfig(config.WithRegion("us-west-2"))
```

## Specifying Credentials

The {{% alias sdk-go %}} requires credentials (an access key and secret access
key) to sign requests to AWS. You can specify your credentials in
several locations, depending on your particular use case. For
information about obtaining credentials, see [Getting Started]({{% relref "getting-started.md" %}}).

When you initialize an `aws.Config` instance using `config.LoadDefaultConfig`,
the SDK uses its default credential chain to find AWS credentials. This
default credential chain looks for credentials in the following order:

1. Environment variables.
   1. Static Credentials (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
   1. Web Identity Token (`AWS_WEB_IDENTITY_TOKEN_FILE`)
1. Shared credentials file.
   1. ~/.aws/credentials
   1. ~/.aws/config
1. If your application uses an ECS task definition or RunTask API operation,
   {{% alias service=IAM %}} role for tasks.
1. If your application is running on an {{% alias service=EC2 %}} instance, {{% alias service=IAM %}} role for {{% alias service=EC2 %}}.

The SDK detects and uses the built-in providers automatically, without
requiring manual configurations. For example, if you use {{% alias service=IAM %}} roles for
{{% alias service=EC2 %}} instances, your applications automatically use the
instance's credentials. You don't need to manually configure credentials
in your application.

As a best practice, AWS recommends that you specify credentials in the
following order:

1. Use {{% alias service=IAM %}} roles for tasks if your application uses an ECS task definition or RunTask API operation.
   
1. Use {{% alias service=IAM %}} roles for {{% alias service=EC2 %}} (if your application is running on an
   {{% alias service=EC2 %}} instance).

   {{% alias service=IAM %}} roles provide applications on the instance temporary security
   credentials to make AWS calls. {{% alias service=IAM %}} roles provide an easy way to
   distribute and manage credentials on multiple {{% alias service=EC2 %}} instances.

1. Use a shared credentials file.

   This credentials file is the same one used by other SDKs and the {{% alias tools=CLI %}}.
   If you're already using a shared credentials file, you can also use
   it for this purpose.

4. Use environment variables.

   Setting environment variables is useful if you're doing development
   work on a machine other than an {{% alias service=EC2 %}} instance.


### {{% alias service=IAM %}} Roles for Tasks

If your application uses an {{% alias service=ECS %}} task definition or `RunTask` operation,
use [IAM Roles for Tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html)
to specify an IAM role that can be used by the containers in a task.

### {{% alias service=IAM %}} Roles for {{% alias service=EC2 %}} Instances

If you are running your application on an {{% alias service=EC2 %}} instance,
use the instance's [IAM role](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html)
to get temporary security credentials to make calls to AWS.

If you have configured your instance to use {{% alias service=IAM %}} roles, the SDK uses
these credentials for your application automatically. You don't need to
manually specify these credentials.

### Shared Credentials File

A credential file is a plaintext file that contains your access keys.
The file must be on the same machine on which you're running your
application. The file must be named `credentials` and located in the
`.aws/` folder in your home directory. The home directory can vary by
operating system. In Windows, you can refer to your home directory by
using the environment variable :code:`%UserProfile%`. In Unix-like systems, you
can use the environment variable :code:`$HOME` or :code:`~` (tilde).

If you already use this file for other SDKs and tools (like the {{% alias tools=CLI %}}),
you don't need to change anything to use the files in this SDK. If
you use different credentials for different tools or applications, you
can use *profiles* to configure multiple access keys in the same
configuration file.

#### Creating the Credentials File

If you don't have a shared credentials file (`.aws/credentials`), you
can use any text editor to create one in your home directory. Add the
following content to your credentials file, replacing
`<YOUR_ACCESS_KEY_ID>` and `<YOUR_SECRET_ACCESS_KEY>` with your
credentials.

```ini
[default]
aws_access_key_id = <YOUR_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_SECRET_ACCESS_KEY>
```

The `[default]` heading defines credentials for the default profile,
which the SDK will use unless you configure it to use another profile.

You can also use temporary security credentials by adding the session
tokens to your profile, as shown in the following example:

```ini
[temp]
aws_access_key_id = <YOUR_TEMP_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_TEMP_SECRET_ACCESS_KEY>
aws_session_token = <YOUR_SESSION_TOKEN>
```

#### Specifying Profiles

You can include multiple access keys in the same configuration file by
associating each set of access keys with a profile. For example, in your
credentials file, you can declare multiple profiles, as follows.

```ini
[default]
aws_access_key_id = <YOUR_DEFAULT_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_DEFAULT_SECRET_ACCESS_KEY>

[test-account]
aws_access_key_id = <YOUR_TEST_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_TEST_SECRET_ACCESS_KEY>

[prod-account]
; work profile
aws_access_key_id = <YOUR_PROD_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_PROD_SECRET_ACCESS_KEY>
```

By default, the SDK checks the `AWS_PROFILE` environment variable to
determine which profile to use. If no `AWS_PROFILE` variable is set,
the SDK uses the default profile.

Sometimes, you may to want to use a different profile with your application.
For example let's say you want to use the `test-account` credentials with
your `myapp` application. You can your application and use this profile by using
the following command:

```
$ AWS_PROFILE=test-account myapp
```

You can also use instruct the SDK to select a profile by either
`os.Setenv("AWS_PROFILE", "test-account")` before calling `config.LoadDefaultConfig`,
or by passing an explicit profile as an argument as shown in the following example:

```go
cfg, err := config.LoadDefaultConfig(config.WithSharedConfigProfile("test-account"))
```

{{% pageinfo color="info" %}}
If you specify credentials in environment variables, the SDK
always uses those credentials, no matter which profile you specify.
{{% /pageinfo %}}

### Environment Variables

By default, the SDK detects AWS credentials set in your environment and
uses them to sign requests to AWS. That way you don't need to manage
credentials in your applications.

The SDK looks for credentials in the following environment variables:

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_SESSION_TOKEN` (optional)

The following examples show how you configure the environment variables.

#### Linux, OS X, or Unix
```
$ export AWS_ACCESS_KEY_ID=YOUR_AKID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
$ export AWS_SESSION_TOKEN=TOKEN
```

#### Windows

```batch
> set AWS_ACCESS_KEY_ID=YOUR_AKID
> set AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
> set AWS_SESSION_TOKEN=TOKEN
```

### Specify Credentials Programmatically
`config.LoadDefaultConfig` allows you to provide an explicit
[aws.CredentialProvider]({{< apiref "aws#CredentialsProvider" >}}) when loading the shared configuration sources.
To pass an explicity credential provider when loading shared configuration use
[config.WithCredentialsProvider]({{< apiref "config#WithCredentialsProvider" >}}). For example, if `customProvider`
references an instance of `aws.CredentialProvider` implementation, it can be passed during configuration loading
like so:

```go
cfg, err := config.LoadDefaultConfig(config.WithCredentialsProvider(customProvider))
```

If you explicitly provide credentials, as in this example, the SDK uses only those credentials.

{{% pageinfo color="info" %}}
All credential providers passed to or returned by `LoadDefaultConfig` are wrapped in a
[CredentialsCache]({{< apiref "aws#CredentialsCache" >}}) automatically. This enables caching and concurrency safe 
credential access. If you explicitly configure a provider on `aws.Config` directly you must explicitly wrap the provider
with this type.
{{% /pageinfo %}}

#### Static Credentials

You can hard-code credentials in your application by using the [credentials.StaticCredentialsProvider]({{< apiref "credentials#StaticCredentialsProvider" >}})
credential provider to explicitly set the access keys to be used. For example:

```go
cfg, err := config.LoadDefaultConfig(
	config.WithCredentialsProvider(aws.StaticCredentialsProvider("AKID", "SECRET_KEY", "TOKEN")),
)
```

{{% pageinfo color="warning" %}}
Do not embed credentials inside an application. Use this
method only for testing purposes.
{{% /pageinfo %}}

#### Other Credentials Providers

The SDK provides other methods for retrieving credentials in the
[credentials]({{< apiref credentials >}}) module. For example, you can retrieve temporary security credentials from {{%
alias service=STSlong %}} or credentials from encrypted storage.

**Available Credential Providers**:

* [ec2rolecreds]({{< apiref "credentials/ec2rolecreds" >}}) &ndash; Retrieve Credentials from {{< alias service=EC2 >}}
  Instances Roles via {{< alias service=EC2 >}} IMDS.

* [endpointcreds]({{< apiref "credentials/endpointcreds" >}}) &ndash; Retrieve Credentials from an arbitrary HTTP
  endpoint.

* [processcreds]({{< apiref "credentials/processcreds" >}}) &ndash; Retrieve Credentials from an external process that
  will be invoked by the host environment's shell.
  
* [stscreds]({{< apiref "credentials/stscreds" >}}) &ndash; Retrieve Credentials from {{% alias service=STS %}}
