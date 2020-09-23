module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200921180648-50b89d38c63c
	github.com/aws/aws-sdk-go-v2/ec2imds v0.0.0-20200923000803-25ae3780f7e8
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200923231307-31335607d195
	github.com/awslabs/smithy-go v0.0.0-20200923183614-866bcae027e6
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/ec2imds => ../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)
