module github.com/aws/aws-sdk-go-v2/credentials

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.2
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.1
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.1
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.1
	github.com/aws/smithy-go v1.20.1
	github.com/google/go-cmp v0.5.8
)

require (
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.2 // indirect
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
