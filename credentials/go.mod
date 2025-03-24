module github.com/Enflick/aws-sdk-go-v2/credentials

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.27.0
	github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds v1.16.3
	github.com/Enflick/aws-sdk-go-v2/service/sso v1.20.9
	github.com/Enflick/aws-sdk-go-v2/service/ssooidc v1.24.3
	github.com/Enflick/aws-sdk-go-v2/service/sts v1.28.10
	github.com/aws/smithy-go v1.20.2
)

require (
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v1.3.7 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v1.11.9 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../

replace github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/sso => ../service/sso/

replace github.com/Enflick/aws-sdk-go-v2/service/ssooidc => ../service/ssooidc/

replace github.com/Enflick/aws-sdk-go-v2/service/sts => ../service/sts/
