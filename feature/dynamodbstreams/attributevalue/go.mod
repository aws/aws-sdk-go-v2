module github.com/Enflick/aws-sdk-go-v2/feature/dynamodbstreams/attributevalue

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/service/dynamodb v0.0.0-20250325221836-b4195dfa2eb5
	github.com/Enflick/aws-sdk-go-v2/service/dynamodbstreams v0.0.0-20250325221836-b4195dfa2eb5
)

require (
	github.com/Enflick/smithy-go v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/Enflick/aws-sdk-go-v2/service/dynamodbstreams => ../../../service/dynamodbstreams/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/
