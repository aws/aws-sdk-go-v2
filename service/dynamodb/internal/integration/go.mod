module github.com/aws/aws-sdk-go-v2/service/dynamodb/internal/integration

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200923001701-7b95ccd95ed9
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200923001701-7b95ccd95ed9
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200923001701-7b95ccd95ed9
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.0.0-00010101000000-000000000000
	github.com/awslabs/smithy-go v0.0.0-20200914213924-b41e7bef5d4f
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config
	github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../
)
