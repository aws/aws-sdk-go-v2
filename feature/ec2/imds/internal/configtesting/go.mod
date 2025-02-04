module github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/configtesting

go 1.21

require (
	github.com/aws/aws-sdk-go-v2/config v1.29.5
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.27
)

require (
	github.com/aws/aws-sdk-go-v2 v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.58 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.13 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../../../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../../service/sts/
