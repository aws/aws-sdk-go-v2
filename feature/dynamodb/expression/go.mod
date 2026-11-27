module github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.48.1
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.21.24
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.74.2
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.48.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.13 // indirect
	github.com/aws/smithy-go v1.28.11 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue => ../../../feature/dynamodb/attributevalue/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/
