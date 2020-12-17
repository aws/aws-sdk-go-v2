module github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201216221327-f18ebfdeb472
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v0.0.0-20201217001131-4ae90bb70aa7
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.30.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue => ../../../feature/dynamodb/attributevalue/
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/
)

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
