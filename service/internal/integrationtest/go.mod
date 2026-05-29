module github.com/aws/aws-sdk-go-v2/service/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v1.41.8
	github.com/aws/aws-sdk-go-v2/config v1.32.19
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.22.21
	github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.53.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.5
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.304.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.90.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.102.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.71.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.28
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.2
	github.com/aws/aws-sdk-go-v2/service/transcribestreaming v1.34.8
	github.com/aws/smithy-go v1.26.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.18 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.36.1 // indirect
)

go 1.24

replace github.com/aws/aws-sdk-go-v2/service/codestar => ../../../service/codestar/

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/feature/s3/manager => ../../../feature/s3/manager/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/bedrockruntime => ../../../service/bedrockruntime/

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/aws/aws-sdk-go-v2/service/ec2 => ../../../service/ec2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/lambda => ../../../service/lambda/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/aws/aws-sdk-go-v2/service/s3control => ../../../service/s3control/

replace github.com/aws/aws-sdk-go-v2/service/signin => ../../../service/signin/

replace github.com/aws/aws-sdk-go-v2/service/sqs => ../../../service/sqs/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/transcribestreaming => ../../../service/transcribestreaming/
