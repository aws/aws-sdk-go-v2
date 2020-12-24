module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201222223005-ee883de66531
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20201222223005-ee883de66531
	github.com/aws/aws-sdk-go-v2/service/sts v0.30.0
	github.com/aws/smithy-go v0.5.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
