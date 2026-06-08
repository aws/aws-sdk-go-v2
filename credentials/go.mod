module github.com/aws/aws-sdk-go-v2/credentials

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.12
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.28
	github.com/aws/aws-sdk-go-v2/service/signin v1.1.4
	github.com/aws/aws-sdk-go-v2/service/sso v1.31.2
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.36.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.43.2
	github.com/aws/smithy-go v1.27.1
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.29 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.28 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/signin => ../service/signin/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../service/sts/
