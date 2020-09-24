module github.com/aws/aws-sdk-go-v2/config

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924003723-be31bb8b3c58
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200924003723-be31bb8b3c58
	github.com/aws/aws-sdk-go-v2/ec2imds v0.0.0-20200924003723-be31bb8b3c58
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200924003723-be31bb8b3c58
	github.com/awslabs/smithy-go v0.0.0-20200923224119-6c637abd889a
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/credentials => ../credentials/
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
