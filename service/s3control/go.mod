module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.2
	github.com/awslabs/smithy-go v0.4.1-0.20201208232924-b8cdbaa577ff
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
