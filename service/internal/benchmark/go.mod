module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v1.3.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.2.1
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v1.2.1
	github.com/aws/smithy-go v1.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.0.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
