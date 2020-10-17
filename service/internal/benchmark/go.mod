module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v0.27.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.27.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.2.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.27.0
	github.com/awslabs/smithy-go v0.2.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
