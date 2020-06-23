module github.com/aws/aws-sdk-go-v2/internal/protocoltest/awsrestjson

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200622232612-3b8c27eea891
	github.com/awslabs/smithy-go v0.0.0-20200621224423-542121c5d59c
	github.com/google/go-cmp v0.4.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/awslabs/smithy-go => ../../../../smithy-go
