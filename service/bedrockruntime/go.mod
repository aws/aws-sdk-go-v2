module github.com/Enflick/aws-sdk-go-v2/service/bedrockruntime

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.0.0
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250321221718-8b300ea454e4
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34
	github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.26.1
	github.com/aws/smithy-go v1.22.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/smithy-go => github.com/aws/smithy-go v1.20.2
