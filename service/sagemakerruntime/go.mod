module github.com/aws/aws-sdk-go-v2/service/sagemakerruntime

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.1
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.20
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.20
	github.com/aws/smithy-go v1.22.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
