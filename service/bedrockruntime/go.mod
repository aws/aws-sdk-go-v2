module github.com/aws/aws-sdk-go-v2/service/bedrockruntime

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.38.1
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.4
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.4
	github.com/aws/smithy-go v1.23.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
