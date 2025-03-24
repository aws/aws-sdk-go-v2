module github.com/Enflick/aws-sdk-go-v2/service/bedrockruntime

go 1.20

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream

replace github.com/Enflick/aws-sdk-go-v2 => ../../

require (
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v0.0.0-20250324190212-47af6a3c9b8a
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250324190212-47af6a3c9b8a
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250324190212-47af6a3c9b8a
	github.com/Enflick/smithy-go v1.3.0
)
