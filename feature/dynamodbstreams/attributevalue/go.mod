module github.com/aws/aws-sdk-go-v2/feature/dynamodbstreams/attributevalue

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.27.0
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.18.7
	github.com/aws/smithy-go v1.19.0
	github.com/google/go-cmp v0.5.8
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/
