module github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.2.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.1.1
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.1.1
	github.com/google/go-cmp v0.5.4
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.0.1
)

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
