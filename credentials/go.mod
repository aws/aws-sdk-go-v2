module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200928205708-45b061121bf8
	github.com/aws/aws-sdk-go-v2/ec2imds v0.0.0-20200928205708-45b061121bf8
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200928205708-45b061121bf8
	github.com/awslabs/smithy-go v0.1.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
