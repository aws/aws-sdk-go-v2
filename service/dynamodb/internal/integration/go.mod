module github.com/aws/aws-sdk-go-v2/service/dynamodb/internal/integration

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200915201900-9bd330b7bf22
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200915201900-9bd330b7bf22
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200915201900-9bd330b7bf22
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.0.0-20200915201900-9bd330b7bf22
	github.com/awslabs/smithy-go v0.0.0-20200828214850-b1c39f43623b
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../../config

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../
