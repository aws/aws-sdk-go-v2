---
title: "Amazon RDS Utilities"
linkTitle: "Amazon RDS"
date: "2021-04-16"
description: "Using the AWS SDK for Go V2 Amazon RDS Utilities"
---

## IAM Authentication

The [auth]({{< apiref "feature/rds/auth" >}}) package provides utilities for generating authentication tokens for
connecting to {{% alias service="RDS" %}} MySQL and PostgreSQL database instances. Using the [BuildAuthToken]({{<
apiref "feature/rds/auth#BuildAuthToken" >}}) method, you generate a database authorization token by providing the
database endpoint, AWS Region, username, and a [aws.CredentialProvider]({{< apiref "aws#CredentialsProvider" >}})
implementation that returns IAM credentials with permission to connect to the database using {{< alias service="IAM" >}}
database authentication. To learn more about configuring {{< alias service="RDS" >}} with {{< alias service="IAM" >}}
authentication see the following {{< alias service="RDS" >}} Developer Guide resources:

* [Enabling and disabling IAM database authentication](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.Enabling.html)
* [Creating and using an IAM policy for IAM database access](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.IAMPolicy.html)
* [Creating a database account using IAM authentication](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.DBAccounts.html)

The following examples shows how to generate an authentication token to connect to an {{< alias service="RDS" >}}
database:

```go
import "context"
import "github.com/aws/aws-sdk-go-v2/config"
import "github.com/aws/aws-sdk-go-v2/feature/rds/auth"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
	panic("configuration error: " + err.Error())
}

authenticationToken, err := auth.BuildAuthToken(
	context.TODO(),
	"mydb.123456789012.us-east-1.rds.amazonaws.com:3306", // Database Endpoint (With Port)
	"us-east-1", // AWS Region
	"jane_doe", // Database Account
	cfg.Credentials,
)
if err != nil {
	panic("failed to create authentication token: " + err.Error())
}
```
