module github.com/aws/aws-sdk-go-v2/service/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.15.15
	github.com/aws/aws-sdk-go-v2/service/acm v1.22.7
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.21.7
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.25.7
	github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice v1.22.6
	github.com/aws/aws-sdk-go-v2/service/appstream v1.31.1
	github.com/aws/aws-sdk-go-v2/service/athena v1.38.0
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.37.0
	github.com/aws/aws-sdk-go-v2/service/batch v1.30.7
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.43.0
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.32.6
	github.com/aws/aws-sdk-go-v2/service/cloudhsmv2 v1.19.7
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.20.6
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.36.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.32.2
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.28.0
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.20.1
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.22.3
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.22.6
	github.com/aws/aws-sdk-go-v2/service/codestar v1.19.6
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.32.1
	github.com/aws/aws-sdk-go-v2/service/configservice v1.44.0
	github.com/aws/aws-sdk-go-v2/service/costandusagereportservice v1.21.6
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.35.7
	github.com/aws/aws-sdk-go-v2/service/devicefarm v1.20.6
	github.com/aws/aws-sdk-go-v2/service/directconnect v1.22.7
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.22.8
	github.com/aws/aws-sdk-go-v2/service/docdb v1.30.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.27.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.146.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.24.7
	github.com/aws/aws-sdk-go-v2/service/ecs v1.38.1
	github.com/aws/aws-sdk-go-v2/service/efs v1.26.6
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.34.7
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.20.6
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.21.7
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.27.0
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.25.0
	github.com/aws/aws-sdk-go-v2/service/elastictranscoder v1.19.6
	github.com/aws/aws-sdk-go-v2/service/emr v1.36.1
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.28.1
	github.com/aws/aws-sdk-go-v2/service/firehose v1.24.0
	github.com/aws/aws-sdk-go-v2/service/gamelift v1.28.1
	github.com/aws/aws-sdk-go-v2/service/glacier v1.19.6
	github.com/aws/aws-sdk-go-v2/service/glue v1.74.0
	github.com/aws/aws-sdk-go-v2/service/health v1.22.6
	github.com/aws/aws-sdk-go-v2/service/iam v1.28.7
	github.com/aws/aws-sdk-go-v2/service/inspector v1.19.6
	github.com/aws/aws-sdk-go-v2/service/iot v1.49.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.24.7
	github.com/aws/aws-sdk-go-v2/service/kms v1.27.9
	github.com/aws/aws-sdk-go-v2/service/lambda v1.49.7
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.34.0
	github.com/aws/aws-sdk-go-v2/service/marketplacecommerceanalytics v1.18.6
	github.com/aws/aws-sdk-go-v2/service/neptune v1.28.2
	github.com/aws/aws-sdk-go-v2/service/opsworks v1.19.6
	github.com/aws/aws-sdk-go-v2/service/pinpointemail v1.17.6
	github.com/aws/aws-sdk-go-v2/service/polly v1.36.6
	github.com/aws/aws-sdk-go-v2/service/rds v1.69.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.39.8
	github.com/aws/aws-sdk-go-v2/service/rekognition v1.36.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.37.1
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.20.6
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.25.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.42.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.26.2
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.26.0
	github.com/aws/aws-sdk-go-v2/service/ses v1.19.6
	github.com/aws/aws-sdk-go-v2/service/sfn v1.24.7
	github.com/aws/aws-sdk-go-v2/service/shield v1.23.6
	github.com/aws/aws-sdk-go-v2/service/snowball v1.24.7
	github.com/aws/aws-sdk-go-v2/service/sns v1.26.7
	github.com/aws/aws-sdk-go-v2/service/sqs v1.29.7
	github.com/aws/aws-sdk-go-v2/service/ssm v1.45.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7
	github.com/aws/aws-sdk-go-v2/service/support v1.19.6
	github.com/aws/aws-sdk-go-v2/service/timestreamwrite v1.23.7
	github.com/aws/aws-sdk-go-v2/service/transcribestreaming v1.15.6
	github.com/aws/aws-sdk-go-v2/service/waf v1.18.6
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.19.6
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.43.6
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.35.8
	github.com/aws/smithy-go v1.19.0
	github.com/google/go-cmp v0.5.8
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.8.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.7 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

go 1.19

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/feature/s3/manager => ../../../feature/s3/manager/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../../internal/v4a/

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

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

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

replace github.com/aws/aws-sdk-go-v2/service/snowball => ../../../service/snowball/

replace github.com/aws/aws-sdk-go-v2/service/sns => ../../../service/sns/

replace github.com/aws/aws-sdk-go-v2/service/sqs => ../../../service/sqs/

replace github.com/aws/aws-sdk-go-v2/service/ssm => ../../../service/ssm/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/support => ../../../service/support/

replace github.com/aws/aws-sdk-go-v2/service/timestreamwrite => ../../../service/timestreamwrite/

replace github.com/aws/aws-sdk-go-v2/service/transcribestreaming => ../../../service/transcribestreaming/

replace github.com/aws/aws-sdk-go-v2/service/waf => ../../../service/waf/

replace github.com/aws/aws-sdk-go-v2/service/wafregional => ../../../service/wafregional/

replace github.com/aws/aws-sdk-go-v2/service/wafv2 => ../../../service/wafv2/

replace github.com/aws/aws-sdk-go-v2/service/workspaces => ../../../service/workspaces/
