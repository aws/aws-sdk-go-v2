module github.com/aws/aws-sdk-go-v2/service/transcribestreaming

go 1.24

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.8
	github.com/aws/smithy-go v1.24.2
	github.com/aws/aws-sdk-go-v2 v1.41.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/
