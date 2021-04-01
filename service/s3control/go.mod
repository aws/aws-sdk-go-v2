module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.3.1
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.2.1
	github.com/aws/smithy-go v1.3.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
