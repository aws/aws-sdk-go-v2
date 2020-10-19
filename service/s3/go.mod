module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.27.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20200930084954-897dfb99530c
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.26.1-0.20201016111247-66b2791dafc4
	github.com/awslabs/smithy-go v0.2.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
