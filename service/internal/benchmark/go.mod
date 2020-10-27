module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v0.28.1-0.20201027184009-8eb8fc303e7c
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.28.0
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.28.0
	github.com/awslabs/smithy-go v0.2.2-0.20201026231331-345290040c23
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
