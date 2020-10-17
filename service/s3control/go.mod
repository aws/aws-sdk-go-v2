module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.26.1-0.20201016111247-66b2791dafc4
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.26.1-0.20201016111247-66b2791dafc4
	github.com/awslabs/smithy-go v0.2.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
