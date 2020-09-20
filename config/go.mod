module github.com/aws/aws-sdk-go-v2/config

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200921180648-50b89d38c63c
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200922233035-113c46575ad6
	github.com/aws/aws-sdk-go-v2/ec2imds v0.0.0-20200922233035-113c46575ad6
	github.com/awslabs/smithy-go v0.0.0-20200922192056-dab44aa99759
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/credentials => ../credentials
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds
)
