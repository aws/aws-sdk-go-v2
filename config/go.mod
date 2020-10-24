module github.com/aws/aws-sdk-go-v2/config

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.27.1-0.20201022222834-4451b4af620e
	github.com/aws/aws-sdk-go-v2/credentials v0.1.2
	github.com/aws/aws-sdk-go-v2/ec2imds v0.1.2
	github.com/aws/aws-sdk-go-v2/service/sts v0.27.0
	github.com/awslabs/smithy-go v0.2.1-0.20201023220843-5834338b6151
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/credentials => ../credentials/
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
