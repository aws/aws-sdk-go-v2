module github.com/aws/aws-sdk-go-v2/service/transcribestreaming/internal/testing

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.9.1
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.6.1
	github.com/aws/smithy-go v1.8.1-0.20211009065425-564a5297a774
	github.com/google/go-cmp v0.5.6
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting => ../../../../service/internal/eventstreamtesting

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/kinesis => ../../../../service/kinesis/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/
