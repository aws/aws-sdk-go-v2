module github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.3.1
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.0.5
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.2.1
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.1.4
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.0.3
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue => ../../../feature/dynamodb/attributevalue/
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/
)

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
