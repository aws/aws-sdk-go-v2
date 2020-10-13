module github.com/aws/aws-sdk-go-v2/feature/s3/manager

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.26.0
	github.com/aws/aws-sdk-go-v2/config v0.1.1
	github.com/aws/aws-sdk-go-v2/service/s3 v0.26.0
	github.com/awslabs/smithy-go v0.1.2-0.20201012175301-b4d8737f29d1
	github.com/google/go-cmp v0.4.1
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../config/
	github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/
	github.com/aws/aws-sdk-go-v2/ec2imds => ../../../ec2imds
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/
	github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts
)
