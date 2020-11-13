module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.1-0.20201112231730-2d786000ccee
	github.com/aws/aws-sdk-go-v2/ec2imds v0.1.4
	github.com/aws/aws-sdk-go-v2/service/sts v0.29.0
	github.com/awslabs/smithy-go v0.3.1-0.20201108010311-62c2a93810b4
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
