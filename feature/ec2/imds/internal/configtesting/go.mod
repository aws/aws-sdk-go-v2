module github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/configtesting

go 1.24

require (
	github.com/aws/aws-sdk-go-v2/config v1.32.10
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.18
)

require (
	github.com/aws/aws-sdk-go-v2 v1.41.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.7 // indirect
	github.com/aws/smithy-go v1.24.2 // indirect
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

replace github.com/aws/aws-sdk-go-v2/service/signin => ../../../../../service/signin/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../../service/sts/
