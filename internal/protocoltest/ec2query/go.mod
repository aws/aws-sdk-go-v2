module github.com/aws/aws-sdk-go-v2/internal/protocoltest/ec2query

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.12.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.3
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.1.0
	github.com/aws/smithy-go v1.10.0
	github.com/google/go-cmp v0.5.6
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/
