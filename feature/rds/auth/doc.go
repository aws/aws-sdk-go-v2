// Package auth is used to generate authentication tokens used to
// connect to a given Amazon Relational Database Service (RDS) database.
//
// Before using the authentication please visit the docs here to ensure
// the database has the proper policies to allow for IAM token authentication.
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html#UsingWithRDS.IAMDBAuth.Availability
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/sdk-utilities-rds.html
package auth
