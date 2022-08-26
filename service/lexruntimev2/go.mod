module github.com/aws/aws-sdk-go-v2/service/lexruntimev2

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.16.11
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.18
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.12
	github.com/aws/smithy-go v1.12.2-0.20220825202034-eec4392b3806
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
