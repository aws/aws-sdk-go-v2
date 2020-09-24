module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200924165317-475229a2f000
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.0.0-20200924165317-475229a2f000
	github.com/awslabs/smithy-go v0.0.0-20200924163652-fc0366622e14
)

replace (
   github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ./../internal/s3shared
   github.com/aws/aws-sdk-go-v2/service/s3control => ./../s3control
)
