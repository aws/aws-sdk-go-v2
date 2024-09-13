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
  "context"
  "log"
  "github.com/aws/aws-sdk-go-v2/config"
)

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
  log.Fatalf("failed to load configuration, %v", err)
}
```

`config.LoadDefaultConfig(context.TODO())` will construct an [aws.Config]({{< apiref "aws#Config" >}})
using the AWS shared configuration sources. This includes configuring a credential provider, configuring the AWS Region,
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
cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
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
   2. Web Identity Token (`AWS_WEB_IDENTITY_TOKEN_FILE`)
1. Shared configuration files.
   1. SDK defaults to `credentials` file under `.aws` folder that is placed in the home folder on your computer.
   1. SDK defaults to `config` file under `.aws` folder that is placed in the home folder on your computer.
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

1. Use shared credentials or config files.
    
   The credentials and config files are shared across other AWS SDKs and {{% alias service=CLI %}}.
   As a security best practice, we recommend using credentials file for setting sensitive values 
   such as access key IDs and secret keys. Here are the 
   [formatting requirements](https://docs.aws.amazon.com/credref/latest/refdocs/file-format.html) for each of these files.

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

### Shared Credentials and Configuration

The shared credentials and config files can be used to share common configuration 
amongst AWS SDKs and other tools. If you use different credentials for different 
tools or applications, you can use *profiles* to configure multiple access keys 
in the same configuration file.

You can provide multiple credential or config files locations using 
`config.LoadOptions`, by default the SDK loads files stored at default 
locations mentioned in the [specifying credentials](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials) 
section.

```go
import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"    
)

// ...

cfg , err := config.LoadDefaultConfig(context.TODO(), 
    config.WithSharedCredentialsFiles(
	[]string{"test/credentials", "data/credentials"},
    ), 
    config.WithSharedConfigFiles(
        []string{"test/config", "data/config"},
    )	
) 

```

When working with shared credentials and config files, if duplicate profiles 
are specified they are merged to resolve a profile. In case of merge conflict,

1. If duplicate profiles are specified within a same credentials/config file,
   the profile properties specified in the latter profile takes precedence. 

1. If duplicate profiles are specified across either multiple credentials files 
   or across multiple config files, the profile properties are resolved as per 
   the order of file input to the `config.LoadOptions`. The profile properties 
   in the latter files take precedence.
   
1. If a profile exists in both credentials file and config file, the credentials file 
   properties take precedence.
   
If need be, you can enable `LogConfigurationWarnings` on `config.LoadOptions`, and 
log the profile resolution steps.  

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
[profile temp]
aws_access_key_id = <YOUR_TEMP_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_TEMP_SECRET_ACCESS_KEY>
aws_session_token = <YOUR_SESSION_TOKEN>
```
The section name for a non-default profile within a credentials file 
must not begin with the word `profile`. You can read more at 
[AWS Tools and SDKs Shared Configuration and Credentials Reference Guide](https://docs.aws.amazon.com/credref/latest/refdocs/file-format.html#file-format-creds).

#### Creating the Config File

If you don't have a shared credentials file (`.aws/config`), you
can use any text editor to create one in your home directory. Add the
following content to your config file, replacing `<REGION>` with the 
desired region.

```ini
[default]
region = <REGION>
```

The `[default]` heading defines config for the default profile,
which the SDK will use unless you configure it to use another profile.

You use named profiles, as shown in the following example:

```ini
[profile named-profile]
region = <REGION>
```

The section name for a non-default profile within a config file
must always begin with the word `profile `, followed by the 
intended profile name. You can read more at 
[AWS Tools and SDKs Shared Configuration and Credentials Reference Guide](https://docs.aws.amazon.com/credref/latest/refdocs/file-format.html#file-format-config).

#### Specifying Profiles

You can include multiple access keys in the same configuration file by
associating each set of access keys with a profile. For example, in your
credentials file, you can declare multiple profiles, as follows.

```ini
[default]
aws_access_key_id = <YOUR_DEFAULT_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_DEFAULT_SECRET_ACCESS_KEY>

[profile test-account]
aws_access_key_id = <YOUR_TEST_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_TEST_SECRET_ACCESS_KEY>

[profile prod-account]
; work profile
aws_access_key_id = <YOUR_PROD_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_PROD_SECRET_ACCESS_KEY>
```

By default, the SDK checks the `AWS_PROFILE` environment variable to
determine which profile to use. If no `AWS_PROFILE` variable is set,
the SDK uses the `default` profile.

Sometimes, you may want to use a different profile with your application.
For example let's say you want to use the `test-account` credentials with
your `myapp` application. You can, use this profile by using
the following command:

```
$ AWS_PROFILE=test-account myapp
```

You can also use instruct the SDK to select a profile by either
`os.Setenv("AWS_PROFILE", "test-account")` before calling `config.LoadDefaultConfig`,
or by passing an explicit profile as an argument as shown in the following example:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), 
	config.WithSharedConfigProfile("test-account"))
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
To pass an explicit credential provider when loading shared configuration use
[config.WithCredentialsProvider]({{< apiref "config#WithCredentialsProvider" >}}). For example, if `customProvider`
references an instance of `aws.CredentialProvider` implementation, it can be passed during configuration loading
like so:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), 
	config.WithCredentialsProvider(customProvider))
```

If you explicitly provide credentials, as in this example, the SDK uses only those credentials.

{{% pageinfo color="info" %}}
All credential providers passed to or returned by `LoadDefaultConfig` are wrapped in a
[CredentialsCache]({{< apiref "aws#CredentialsCache" >}}) automatically. This enables caching, and credential rotation that is concurrency safe. If you explicitly configure a provider on `aws.Config` directly you must also explicitly wrap the provider
with this type using [NewCredentialsCache]({{< apiref "aws#NewCredentialsCache" >}}).
{{% /pageinfo %}}

#### Static Credentials

You can hard-code credentials in your application by using the [credentials.NewStaticCredentialsProvider]({{< apiref "credentials#NewStaticCredentialsProvider" >}})
credential provider to explicitly set the access keys to be used. For example:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(), 
	config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET_KEY", "TOKEN")),
)
```

{{% pageinfo color="warning" %}}
Do not embed credentials inside an application. Use this
method only for testing purposes.
{{% /pageinfo %}}

#### Single Sign-on Credentials

The SDK provides a credential provider for retrieving temporary AWS credentials using {{% alias service=SSOlong %}}.
Using the {{% alias service=CLI %}}, you authenticate with the AWS access portal and authorize access to temporary
AWS credentials. You then configure your application to load the single sign-on (SSO) profile, and the SDK uses your
SSO credentials to retrieve temporary AWS credentials that will be automatically renewed if expired.
If your SSO credentials expire, you must explicitly renew them by logging in to your
{{% alias service=SSO %}} account again using the {{% alias service=CLI %}}.

For example, you can create a profile, `dev-profile`, authenticate and authorize that profile using the
{{% alias service=CLI %}}, and configure your application as shown below.

1. First create the `profile` and `sso-session`

```
[profile dev-profile]
sso_session = dev-session
sso_account_id = 012345678901
sso_role_name = Developer
region = us-east-1

[sso-session dev-session]
sso_region = us-west-2
sso_start_url = https://company-sso-portal.awsapps.com/start
sso_registration_scopes = sso:account:access

```
2. Login using the {{% alias service=CLI %}} to authenticate and authorize the SSO profile.
```
$ aws --profile dev-profile sso login 
Attempting to automatically open the SSO authorization page in your default browser.
If the browser does not open or you wish to use a different device to authorize this request, open the following URL:

https://device.sso.us-west-2.amazonaws.com/

Then enter the code:

ABCD-EFGH
Successully logged into Start URL: https://company-sso-portal.awsapps.com/start
```
3. Next configure your application to use the SSO profile.
```go
import "github.com/aws/aws-sdk-go-v2/config"

// ...

cfg, err := config.LoadDefaultConfig(
	context.Background(),
	config.WithSharedConfigProfile("dev-profile"),
)
if err != nil {
	return err
}
```

For more information on configuring
SSO profiles and authenticating using the {{% alias service=CLI %}} see
[Configuring the {{% alias service=CLI %}} to use {{% alias service=SSOlong %}}](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html)
in the {{% alias service=CLI %}} User Guide. For more information on programmatically constructing the
SSO credential provider see the [ssocreds]({{< apiref "credentials/ssocreds" >}}) API reference
documentation.

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
