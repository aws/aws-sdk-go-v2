module github.com/aws/aws-sdk-go-v2/credentials

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.1-0.20210108204630-4822f3195720
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v0.1.0
	github.com/aws/aws-sdk-go-v2/service/sts v0.31.0
	github.com/aws/smithy-go v0.5.1-0.20210115041537-09631dea532e
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/
	github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
)

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/
