module github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.31.0
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.18
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.18
	github.com/aws/smithy-go v1.21.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
