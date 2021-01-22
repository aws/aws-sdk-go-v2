module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v1.0.1-0.20210122214637-6cf9ad2f8e2f
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.0.0
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v1.0.0
	github.com/aws/smithy-go v1.0.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
