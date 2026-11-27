module github.com/aws/aws-sdk-go-v2/service/internal/integrationtest

require (
	github.com/aws/aws-sdk-go-v2 v1.48.2
	github.com/aws/aws-sdk-go-v2/config v1.33.23
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.23.30
	github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.64.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.74.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.350.2
	github.com/aws/aws-sdk-go-v2/service/lambda v1.118.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.116.2
	github.com/aws/aws-sdk-go-v2/service/s3control v1.81.3
	github.com/aws/aws-sdk-go-v2/service/sqs v1.53.2
	github.com/aws/aws-sdk-go-v2/service/sts v1.53.2
	github.com/aws/aws-sdk-go-v2/service/transcribestreaming v1.44.16
	github.com/aws/smithy-go v1.28.11
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.29 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.21.5 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.14.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.39.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.44.7 // indirect
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
