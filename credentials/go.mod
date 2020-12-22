module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201221101722-677dd4a81dad
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20201216182737-05d6a8e2a8df
	github.com/aws/aws-sdk-go-v2/service/sts v0.30.0
	github.com/aws/smithy-go v0.4.1-0.20201222001052-74df8ddd8c79
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
