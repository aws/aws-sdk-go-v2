module github.com/aws/aws-sdk-go-v2/service/lexruntimeservice/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.16
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924195609-61ca34a7860e
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.0.0-00010101000000-000000000000
	github.com/awslabs/smithy-go v0.0.0-20200924163652-fc0366622e14
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../
