module github.com/aws/aws-sdk-go-v2/service/lambda

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.19.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.35
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.29
	github.com/aws/smithy-go v1.13.5
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
