module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200921195903-c1c5f2a1c25c
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200921195903-c1c5f2a1c25c
	github.com/awslabs/smithy-go v0.0.0-20200917082847-627658712072
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts
)
