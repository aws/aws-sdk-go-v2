module github.com/aws/aws-sdk-go-v2/service/proton

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.23.3
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.6
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.6
	github.com/aws/smithy-go v1.18.0
	github.com/google/go-cmp v0.5.8
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
