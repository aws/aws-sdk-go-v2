module canaries

go 1.23.0

replace github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient => ../..

require (
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/config v1.31.17
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.52.2
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.18.21 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.19 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.8.19 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.39.1 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
)
