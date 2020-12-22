module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201221101722-677dd4a81dad
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.30.0
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.30.0
	github.com/aws/smithy-go v0.4.1-0.20201222001052-74df8ddd8c79
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
