module github.com/aws/aws-sdk-go-v2/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924224914-965b6782bf3d
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/acm v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/emr v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/eventbridge v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/firehose v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/gamelift v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/glacier v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/glue v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/health v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/iam v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/inspector v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/iot v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/kinesis v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/kms v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/lambda v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/lightsail v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/neptune v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/opsworks v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/pinpointemail v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/polly v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/rds v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/redshift v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/rekognition v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/route53 v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/route53domains v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/route53resolver v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/s3 v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/ses v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/sfn v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/shield v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/sms v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/snowball v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/sns v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/sqs v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/ssm v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200924224731-ccbcb2eb486d
	github.com/aws/aws-sdk-go-v2/service/support v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/waf v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/wafregional v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/wafv2 v0.0.0-20200924231056-a86f2aee748d
	github.com/aws/aws-sdk-go-v2/service/workspaces v0.0.0-20200924231056-a86f2aee748d
	github.com/awslabs/smithy-go v0.0.0-20200924210334-28773c6e7960
)

replace github.com/aws/aws-sdk-go-v2/service/kms => ./../../service/kms

go 1.15
