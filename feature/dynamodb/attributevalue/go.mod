module github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.28
	github.com/aws/aws-sdk-go-v2 v0.29.1-0.20201115205015-a82264590e72
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.29.1-0.20201115205015-a82264590e72
	github.com/google/go-cmp v0.5.3
)

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb
