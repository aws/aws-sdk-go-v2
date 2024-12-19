module github.com/aws/aws-sdk-go-v2/service/kinesis

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.7
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.7
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.26
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.26
	github.com/aws/smithy-go v1.22.1
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
