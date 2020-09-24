module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200924095642-be147c6e7568
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.0.0-20200924095642-be147c6e7568
	github.com/awslabs/smithy-go v0.0.0-20200924081159-7ac2e6483c86
)

replace (
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ./../internal/s3shared
	github.com/aws/aws-sdk-go-v2/service/s3 => ./../s3
)
