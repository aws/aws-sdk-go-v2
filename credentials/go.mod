module github.com/aws/aws-sdk-go-v2/credentials

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.41.1
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.17
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.5
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.9
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.13
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.6
	github.com/aws/smithy-go v1.24.0
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.17 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/signin => ../service/signin/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
