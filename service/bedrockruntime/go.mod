module github.com/Enflick/aws-sdk-go-v2/service/bedrockruntime

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.0.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/smithy-go v1.20.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/smithy-go => github.com/aws/smithy-go v1.20.0
