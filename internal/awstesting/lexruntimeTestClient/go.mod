module github.com/aws/aws-sdk-go-v2/internal/awstesting/lexruntimeTestClient

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v0.21.1-0.20200622163451-5b12b74d16d9
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.0.0-20200622163451-5b12b74d16d9
	github.com/awslabs/smithy-go v0.0.0-20200622173248-cd334615e99f
)

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/awslabs/smithy-go => ../../../../smithy-go
