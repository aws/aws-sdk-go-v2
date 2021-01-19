module github.com/aws/aws-sdk-go-v2/service/sts

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.1-0.20210108204630-4822f3195720
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.2.1-0.20210108204630-4822f3195720
	github.com/aws/smithy-go v1.0.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
