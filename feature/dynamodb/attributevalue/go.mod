module github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.9.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.5.0
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.4.0
	github.com/aws/smithy-go v1.8.0
	github.com/google/go-cmp v0.5.6
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/
