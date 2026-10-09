module github.com/aws/aws-sdk-go-v2/service/bedrockruntime

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.3
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.22
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.6
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.6
	github.com/aws/smithy-go v1.28.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
