module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.28.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.0
	github.com/awslabs/smithy-go v0.2.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
