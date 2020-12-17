module github.com/aws/aws-sdk-go-v2/service/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201217001905-4acf9c65b2d1
	github.com/aws/aws-sdk-go-v2/config v0.3.0
	github.com/aws/aws-sdk-go-v2/service/acm v0.30.0
	github.com/aws/aws-sdk-go-v2/service/apigateway v0.30.0
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v0.30.0
	github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/appstream v0.30.0
	github.com/aws/aws-sdk-go-v2/service/athena v0.30.0
	github.com/aws/aws-sdk-go-v2/service/autoscaling v0.30.0
	github.com/aws/aws-sdk-go-v2/service/batch v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudformation v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudfront v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudhsmv2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v0.30.0
	github.com/aws/aws-sdk-go-v2/service/codebuild v0.30.0
	github.com/aws/aws-sdk-go-v2/service/codecommit v0.30.0
	github.com/aws/aws-sdk-go-v2/service/codedeploy v0.30.0
	github.com/aws/aws-sdk-go-v2/service/codepipeline v0.30.0
	github.com/aws/aws-sdk-go-v2/service/codestar v0.30.0
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v0.30.0
	github.com/aws/aws-sdk-go-v2/service/configservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/costandusagereportservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/devicefarm v0.30.0
	github.com/aws/aws-sdk-go-v2/service/directconnect v0.30.0
	github.com/aws/aws-sdk-go-v2/service/directoryservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/docdb v0.30.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.30.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/ecr v0.30.0
	github.com/aws/aws-sdk-go-v2/service/ecs v0.30.0
	github.com/aws/aws-sdk-go-v2/service/efs v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v0.30.0
	github.com/aws/aws-sdk-go-v2/service/elastictranscoder v0.30.0
	github.com/aws/aws-sdk-go-v2/service/emr v0.30.0
	github.com/aws/aws-sdk-go-v2/service/eventbridge v0.30.0
	github.com/aws/aws-sdk-go-v2/service/firehose v0.30.0
	github.com/aws/aws-sdk-go-v2/service/gamelift v0.30.0
	github.com/aws/aws-sdk-go-v2/service/glacier v0.30.0
	github.com/aws/aws-sdk-go-v2/service/glue v0.30.0
	github.com/aws/aws-sdk-go-v2/service/health v0.30.0
	github.com/aws/aws-sdk-go-v2/service/iam v0.30.0
	github.com/aws/aws-sdk-go-v2/service/inspector v0.30.0
	github.com/aws/aws-sdk-go-v2/service/iot v0.30.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v0.30.0
	github.com/aws/aws-sdk-go-v2/service/kms v0.30.0
	github.com/aws/aws-sdk-go-v2/service/lambda v0.30.0
	github.com/aws/aws-sdk-go-v2/service/lightsail v0.30.0
	github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics v0.30.0
	github.com/aws/aws-sdk-go-v2/service/neptune v0.30.0
	github.com/aws/aws-sdk-go-v2/service/opsworks v0.30.0
	github.com/aws/aws-sdk-go-v2/service/pinpointemail v0.30.0
	github.com/aws/aws-sdk-go-v2/service/polly v0.30.0
	github.com/aws/aws-sdk-go-v2/service/rds v0.30.0
	github.com/aws/aws-sdk-go-v2/service/redshift v0.30.0
	github.com/aws/aws-sdk-go-v2/service/rekognition v0.30.0
	github.com/aws/aws-sdk-go-v2/service/route53 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/route53domains v0.30.0
	github.com/aws/aws-sdk-go-v2/service/route53resolver v0.30.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/s3control v0.30.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v0.30.0
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v0.30.0
	github.com/aws/aws-sdk-go-v2/service/ses v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sfn v0.30.0
	github.com/aws/aws-sdk-go-v2/service/shield v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sms v0.30.0
	github.com/aws/aws-sdk-go-v2/service/snowball v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sns v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sqs v0.30.0
	github.com/aws/aws-sdk-go-v2/service/ssm v0.30.0
	github.com/aws/aws-sdk-go-v2/service/sts v0.30.0
	github.com/aws/aws-sdk-go-v2/service/support v0.30.0
	github.com/aws/aws-sdk-go-v2/service/waf v0.30.0
	github.com/aws/aws-sdk-go-v2/service/wafregional v0.30.0
	github.com/aws/aws-sdk-go-v2/service/wafv2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/workspaces v0.30.0
	github.com/awslabs/smithy-go v0.4.1-0.20201216214517-20e212c92831
	github.com/google/go-cmp v0.5.4
)

go 1.15

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

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

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/support => ../../../service/support/

replace github.com/aws/aws-sdk-go-v2/service/waf => ../../../service/waf/

replace github.com/aws/aws-sdk-go-v2/service/wafregional => ../../../service/wafregional/

replace github.com/aws/aws-sdk-go-v2/service/wafv2 => ../../../service/wafv2/

replace github.com/aws/aws-sdk-go-v2/service/workspaces => ../../../service/workspaces/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
