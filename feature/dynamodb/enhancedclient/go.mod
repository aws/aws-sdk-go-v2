module github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient

go 1.24.0

replace github.com/aws/aws-sdk-go-v2/config => ../../../config

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb

replace github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression => ../../../feature/dynamodb/expression

require (
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.4
	github.com/aws/smithy-go v1.22.4
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.19.3 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.36 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.25.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19 // indirect
)
