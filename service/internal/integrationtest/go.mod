module github.com/aws/aws-sdk-go-v2/service/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v1.7.1
	github.com/aws/aws-sdk-go-v2/config v1.5.0
	github.com/aws/aws-sdk-go-v2/service/acm v1.5.0
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.5.1
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.4.1
	github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice v1.3.1
	github.com/aws/aws-sdk-go-v2/service/appstream v1.5.1
	github.com/aws/aws-sdk-go-v2/service/athena v1.4.1
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.10.0
	github.com/aws/aws-sdk-go-v2/service/batch v1.5.1
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.7.0
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.7.0
	github.com/aws/aws-sdk-go-v2/service/cloudhsmv2 v1.3.1
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.5.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.4.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.6.0
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.5.1
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.3.1
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.5.1
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.4.1
	github.com/aws/aws-sdk-go-v2/service/codestar v1.3.1
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.4.1
	github.com/aws/aws-sdk-go-v2/service/configservice v1.6.1
	github.com/aws/aws-sdk-go-v2/service/costandusagereportservice v1.4.1
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.6.0
	github.com/aws/aws-sdk-go-v2/service/devicefarm v1.4.1
	github.com/aws/aws-sdk-go-v2/service/directconnect v1.6.0
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.3.1
	github.com/aws/aws-sdk-go-v2/service/docdb v1.8.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.4.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.12.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.4.1
	github.com/aws/aws-sdk-go-v2/service/ecs v1.6.1
	github.com/aws/aws-sdk-go-v2/service/efs v1.5.1
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.8.0
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.5.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.5.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.5.0
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.5.1
	github.com/aws/aws-sdk-go-v2/service/elastictranscoder v1.4.1
	github.com/aws/aws-sdk-go-v2/service/emr v1.4.1
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.5.1
	github.com/aws/aws-sdk-go-v2/service/firehose v1.4.1
	github.com/aws/aws-sdk-go-v2/service/gamelift v1.5.1
	github.com/aws/aws-sdk-go-v2/service/glacier v1.4.1
	github.com/aws/aws-sdk-go-v2/service/glue v1.9.0
	github.com/aws/aws-sdk-go-v2/service/health v1.6.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.7.0
	github.com/aws/aws-sdk-go-v2/service/inspector v1.3.1
	github.com/aws/aws-sdk-go-v2/service/iot v1.8.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.5.1
	github.com/aws/aws-sdk-go-v2/service/kms v1.4.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.5.0
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.8.0
	github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics v1.3.1
	github.com/aws/aws-sdk-go-v2/service/neptune v1.7.0
	github.com/aws/aws-sdk-go-v2/service/opsworks v1.4.1
	github.com/aws/aws-sdk-go-v2/service/pinpointemail v1.3.1
	github.com/aws/aws-sdk-go-v2/service/polly v1.4.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.6.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.9.0
	github.com/aws/aws-sdk-go-v2/service/rekognition v1.6.1
	github.com/aws/aws-sdk-go-v2/service/route53 v1.7.1
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.3.1
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.5.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.11.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.8.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.4.1
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.4.1
	github.com/aws/aws-sdk-go-v2/service/ses v1.5.0
	github.com/aws/aws-sdk-go-v2/service/sfn v1.4.1
	github.com/aws/aws-sdk-go-v2/service/shield v1.5.1
	github.com/aws/aws-sdk-go-v2/service/sms v1.3.1
	github.com/aws/aws-sdk-go-v2/service/snowball v1.5.1
	github.com/aws/aws-sdk-go-v2/service/sns v1.7.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.7.0
	github.com/aws/aws-sdk-go-v2/service/ssm v1.8.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.6.0
	github.com/aws/aws-sdk-go-v2/service/support v1.4.1
	github.com/aws/aws-sdk-go-v2/service/timestreamwrite v1.3.1
	github.com/aws/aws-sdk-go-v2/service/waf v1.3.1
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.4.1
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.6.1
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.5.1
	github.com/aws/smithy-go v1.6.1-0.20210716220526-e488a706561f
	github.com/google/go-cmp v0.5.6
)

go 1.15

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/service/acm => ../../../service/acm/

replace github.com/aws/aws-sdk-go-v2/service/apigateway => ../../../service/apigateway/

replace github.com/aws/aws-sdk-go-v2/service/applicationautoscaling => ../../../service/applicationautoscaling/

replace github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice => ../../../service/applicationdiscoveryservice/

replace github.com/aws/aws-sdk-go-v2/service/appstream => ../../../service/appstream/

replace github.com/aws/aws-sdk-go-v2/service/athena => ../../../service/athena/

replace github.com/aws/aws-sdk-go-v2/service/autoscaling => ../../../service/autoscaling/

replace github.com/aws/aws-sdk-go-v2/service/batch => ../../../service/batch/

replace github.com/aws/aws-sdk-go-v2/service/cloudformation => ../../../service/cloudformation/

replace github.com/aws/aws-sdk-go-v2/service/cloudfront => ../../../service/cloudfront/

replace github.com/aws/aws-sdk-go-v2/service/cloudhsmv2 => ../../../service/cloudhsmv2/

replace github.com/aws/aws-sdk-go-v2/service/cloudsearch => ../../../service/cloudsearch/

replace github.com/aws/aws-sdk-go-v2/service/cloudtrail => ../../../service/cloudtrail/

replace github.com/aws/aws-sdk-go-v2/service/cloudwatch => ../../../service/cloudwatch/

replace github.com/aws/aws-sdk-go-v2/service/codebuild => ../../../service/codebuild/

replace github.com/aws/aws-sdk-go-v2/service/codecommit => ../../../service/codecommit/

replace github.com/aws/aws-sdk-go-v2/service/codedeploy => ../../../service/codedeploy/

replace github.com/aws/aws-sdk-go-v2/service/codepipeline => ../../../service/codepipeline/

replace github.com/aws/aws-sdk-go-v2/service/codestar => ../../../service/codestar/

replace github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider => ../../../service/cognitoidentityprovider/

replace github.com/aws/aws-sdk-go-v2/service/configservice => ../../../service/configservice/

replace github.com/aws/aws-sdk-go-v2/service/costandusagereportservice => ../../../service/costandusagereportservice/

replace github.com/aws/aws-sdk-go-v2/service/databasemigrationservice => ../../../service/databasemigrationservice/

replace github.com/aws/aws-sdk-go-v2/service/devicefarm => ../../../service/devicefarm/

replace github.com/aws/aws-sdk-go-v2/service/directconnect => ../../../service/directconnect/

replace github.com/aws/aws-sdk-go-v2/service/directoryservice => ../../../service/directoryservice/

replace github.com/aws/aws-sdk-go-v2/service/docdb => ../../../service/docdb/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/ec2 => ../../../service/ec2/

replace github.com/aws/aws-sdk-go-v2/service/ecr => ../../../service/ecr/

replace github.com/aws/aws-sdk-go-v2/service/ecs => ../../../service/ecs/

replace github.com/aws/aws-sdk-go-v2/service/efs => ../../../service/efs/

replace github.com/aws/aws-sdk-go-v2/service/elasticache => ../../../service/elasticache/

replace github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk => ../../../service/elasticbeanstalk/

replace github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing => ../../../service/elasticloadbalancing/

replace github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 => ../../../service/elasticloadbalancingv2/

replace github.com/aws/aws-sdk-go-v2/service/elasticsearchservice => ../../../service/elasticsearchservice/

replace github.com/aws/aws-sdk-go-v2/service/elastictranscoder => ../../../service/elastictranscoder/

replace github.com/aws/aws-sdk-go-v2/service/emr => ../../../service/emr/

replace github.com/aws/aws-sdk-go-v2/service/eventbridge => ../../../service/eventbridge/

replace github.com/aws/aws-sdk-go-v2/service/firehose => ../../../service/firehose/

replace github.com/aws/aws-sdk-go-v2/service/gamelift => ../../../service/gamelift/

replace github.com/aws/aws-sdk-go-v2/service/glacier => ../../../service/glacier/

replace github.com/aws/aws-sdk-go-v2/service/glue => ../../../service/glue/

replace github.com/aws/aws-sdk-go-v2/service/health => ../../../service/health/

replace github.com/aws/aws-sdk-go-v2/service/iam => ../../../service/iam/

replace github.com/aws/aws-sdk-go-v2/service/inspector => ../../../service/inspector/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/iot => ../../../service/iot/

replace github.com/aws/aws-sdk-go-v2/service/kinesis => ../../../service/kinesis/

replace github.com/aws/aws-sdk-go-v2/service/kms => ../../../service/kms/

replace github.com/aws/aws-sdk-go-v2/service/lambda => ../../../service/lambda/

replace github.com/aws/aws-sdk-go-v2/service/lightsail => ../../../service/lightsail/

replace github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics => ../../../service/marketplacecommerceanalytics/

replace github.com/aws/aws-sdk-go-v2/service/neptune => ../../../service/neptune/

replace github.com/aws/aws-sdk-go-v2/service/opsworks => ../../../service/opsworks/

replace github.com/aws/aws-sdk-go-v2/service/pinpointemail => ../../../service/pinpointemail/

replace github.com/aws/aws-sdk-go-v2/service/polly => ../../../service/polly/

replace github.com/aws/aws-sdk-go-v2/service/rds => ../../../service/rds/

replace github.com/aws/aws-sdk-go-v2/service/redshift => ../../../service/redshift/

replace github.com/aws/aws-sdk-go-v2/service/rekognition => ../../../service/rekognition/

replace github.com/aws/aws-sdk-go-v2/service/route53 => ../../../service/route53/

replace github.com/aws/aws-sdk-go-v2/service/route53domains => ../../../service/route53domains/

replace github.com/aws/aws-sdk-go-v2/service/route53resolver => ../../../service/route53resolver/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/aws/aws-sdk-go-v2/service/s3control => ../../../service/s3control/

replace github.com/aws/aws-sdk-go-v2/service/secretsmanager => ../../../service/secretsmanager/

replace github.com/aws/aws-sdk-go-v2/service/servicecatalog => ../../../service/servicecatalog/

replace github.com/aws/aws-sdk-go-v2/service/ses => ../../../service/ses/

replace github.com/aws/aws-sdk-go-v2/service/sfn => ../../../service/sfn/

replace github.com/aws/aws-sdk-go-v2/service/shield => ../../../service/shield/

replace github.com/aws/aws-sdk-go-v2/service/sms => ../../../service/sms/

replace github.com/aws/aws-sdk-go-v2/service/snowball => ../../../service/snowball/

replace github.com/aws/aws-sdk-go-v2/service/sns => ../../../service/sns/

replace github.com/aws/aws-sdk-go-v2/service/sqs => ../../../service/sqs/

replace github.com/aws/aws-sdk-go-v2/service/ssm => ../../../service/ssm/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/support => ../../../service/support/

replace github.com/aws/aws-sdk-go-v2/service/timestreamwrite => ../../../service/timestreamwrite/

replace github.com/aws/aws-sdk-go-v2/service/waf => ../../../service/waf/

replace github.com/aws/aws-sdk-go-v2/service/wafregional => ../../../service/wafregional/

replace github.com/aws/aws-sdk-go-v2/service/wafv2 => ../../../service/wafv2/

replace github.com/aws/aws-sdk-go-v2/service/workspaces => ../../../service/workspaces/
