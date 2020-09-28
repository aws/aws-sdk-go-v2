module github.com/aws/aws-sdk-go-v2/service/lexruntimeservice/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.16
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924225551-a2b886903b8b
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.0.0-20200924225551-a2b886903b8b
	github.com/awslabs/smithy-go v0.1.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../
