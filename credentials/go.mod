module github.com/aws/aws-sdk-go-v2/credentials

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.39.0
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.7
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.2
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.34.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.3
	github.com/aws/smithy-go v1.23.0
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.7 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
