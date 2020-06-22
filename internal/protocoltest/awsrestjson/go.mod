module github.com/aws/aws-sdk-go-v2/internal/protocoltest/awsrestjson

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200608172716-6b3036355dc9
	github.com/awslabs/smithy-go v0.0.0-20200621224423-542121c5d59c
	github.com/google/go-cmp v0.4.1
)

replace (
    github.com/aws/aws-sdk-go-v2 => /volumes/brazil/workspace/sdk/aws-sdk-go-v2/sdk
    github.com/awslabs/smithy-go => /volumes/brazil/workspace/sdk/aws-sdk-go-v2/smithy-go
)
