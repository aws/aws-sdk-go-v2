module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v1.2.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.1.2
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v1.1.2
	github.com/aws/smithy-go v1.2.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
