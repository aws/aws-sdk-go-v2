module github.com/aws/aws-sdk-go-v2/config

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200921191144-68cec59c7dee
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200921191144-68cec59c7dee
	github.com/awslabs/smithy-go v0.0.0-20200914213924-b41e7bef5d4f
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/aws/aws-sdk-go-v2/credentials => ../credentials
