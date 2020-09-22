module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200921214008-bd52d2b3309e
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200921214008-bd52d2b3309e
	github.com/awslabs/smithy-go v0.0.0-20200921163940-68b34a89572d
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts
)
