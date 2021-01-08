module github.com/aws/aws-sdk-go-v2/service/internal/benchmark

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.33
	github.com/aws/aws-sdk-go-v2 v0.31.1-0.20210108183639-b6b5057e2ab1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.31.0
	github.com/aws/aws-sdk-go-v2/service/lexruntimeservice v0.31.0
	github.com/aws/smithy-go v0.5.1-0.20210107224202-ae5323020d60
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
