module github.com/aws/aws-sdk-go-v2/config

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.1-0.20201113222241-726e4a15683d
	github.com/aws/aws-sdk-go-v2/credentials v0.1.4
	github.com/aws/aws-sdk-go-v2/ec2imds v0.1.4
	github.com/aws/aws-sdk-go-v2/service/sts v0.29.0
	github.com/awslabs/smithy-go v0.4.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/credentials => ../credentials/
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
