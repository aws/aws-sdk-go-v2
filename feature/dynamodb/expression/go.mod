module github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v0.1.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.31.0
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v0.31.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.4.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue => ../../../feature/dynamodb/attributevalue/
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/
)

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
