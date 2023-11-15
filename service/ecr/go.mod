module github.com/aws/aws-sdk-go-v2/service/ecr

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.23.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.3
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.3
	github.com/aws/smithy-go v1.17.0
	github.com/google/go-cmp v0.5.8
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
