module github.com/Enflick/aws-sdk-go-v2/service/kinesis

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/smithy-go v1.3.0
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
