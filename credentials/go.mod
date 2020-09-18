module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200918204210-51fd26a6f991
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200918170804-ed5c823eb142
	github.com/awslabs/smithy-go v0.0.0-20200914213924-b41e7bef5d4f
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts
)
