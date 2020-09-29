module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924224914-965b6782bf3d
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.0.0-20200928200900-9b4f334f82b2
	github.com/awslabs/smithy-go v0.1.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
