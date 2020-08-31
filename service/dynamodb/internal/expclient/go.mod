module github.com/aws/aws-sdk-go-v2/service/dynamodb/internal/expclient

go 1.14

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200828232643-9cee6f194dae
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.0.0-00010101000000-000000000000
	github.com/awslabs/smithy-go v0.0.0-20200831191241-f0896471bee5
)

replace github.com/awslabs/smithy-go => ../../../../../smithy-go
