module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201217001905-4acf9c65b2d1
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20201216182737-05d6a8e2a8df
	github.com/aws/aws-sdk-go-v2/service/sts v0.30.0
	github.com/awslabs/smithy-go v0.4.1-0.20201216214517-20e212c92831
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
