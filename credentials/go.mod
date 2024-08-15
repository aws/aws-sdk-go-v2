module github.com/aws/aws-sdk-go-v2/credentials

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.30.4
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.12
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.5
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.4
	github.com/aws/smithy-go v1.20.4
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18 // indirect
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
