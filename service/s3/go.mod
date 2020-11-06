module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.1-0.20201027184009-8eb8fc303e7c
	github.com/awslabs/smithy-go v0.3.1-0.20201106213844-55a49a94218d
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
