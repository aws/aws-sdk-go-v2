module github.com/Enflick/aws-sdk-go-v2/service/lambda

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.27.0
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v1.3.7
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7
	github.com/aws/smithy-go v1.20.2
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
