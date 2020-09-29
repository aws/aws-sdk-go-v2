module github.com/aws/aws-sdk-go-v2/service/s3/internal/configtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/config v0.1.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.1.0
	github.com/aws/aws-sdk-go-v2 v0.25.0
	github.com/aws/aws-sdk-go-v2/credentials v0.1.0
	github.com/aws/aws-sdk-go-v2/ec2imds v0.1.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.1.0
	github.com/aws/aws-sdk-go-v2/service/sts v0.1.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config/
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../
)
