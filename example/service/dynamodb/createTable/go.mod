module github.com/Enflick/aws-sdk-go-v2/example/service/dynamodb/createTable

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.5.0
	github.com/Enflick/aws-sdk-go-v2/config v0.0.0-20250325221836-b4195dfa2eb5
	github.com/Enflick/aws-sdk-go-v2/service/dynamodb v0.0.0-20250325221836-b4195dfa2eb5
)

require (
	github.com/Enflick/aws-sdk-go-v2/credentials v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/ini v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/sso v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/ssooidc v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/sts v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/smithy-go v1.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../../

replace github.com/Enflick/aws-sdk-go-v2/config => ../../../../config/

replace github.com/Enflick/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/internal/ini => ../../../../internal/ini/

replace github.com/Enflick/aws-sdk-go-v2/service/dynamodb => ../../../../service/dynamodb/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../../service/internal/endpoint-discovery/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/Enflick/aws-sdk-go-v2/service/ssooidc => ../../../../service/ssooidc/

replace github.com/Enflick/aws-sdk-go-v2/service/sts => ../../../../service/sts/
