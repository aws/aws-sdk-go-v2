module github.com/aws/aws-sdk-go-v2/service/dynamodb/internal/integration

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924201813-c4bcc170fc79
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200923001701-7b95ccd95ed9
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.0.0-00010101000000-000000000000
	github.com/awslabs/smithy-go v0.0.0-20200924163652-fc0366622e14
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config/
	github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../
)
