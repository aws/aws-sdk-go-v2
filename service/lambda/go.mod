module github.com/aws/aws-sdk-go-v2/service/lambda

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.6
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.11
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.37
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.37
	github.com/aws/smithy-go v1.22.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
