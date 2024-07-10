module github.com/aws/aws-sdk-go-v2/example/service/dynamodb/listItems

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.2
	github.com/aws/aws-sdk-go-v2/config v1.27.25
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.14.8
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.34.2
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.25 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.22.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.2 // indirect
	github.com/aws/smithy-go v1.20.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue => ../../../../feature/dynamodb/attributevalue/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/dynamodbstreams => ../../../../service/dynamodbstreams/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../../service/internal/endpoint-discovery/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/
