module github.com/Enflick/aws-sdk-go-v2/credentials

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/service/sso v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/service/ssooidc v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/service/sts v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/smithy-go v1.3.0
)

require (
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20250325155711-0a4bf6fdbeb3 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../

replace github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds => ../feature/ec2/imds/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/sso => ../service/sso/

replace github.com/Enflick/aws-sdk-go-v2/service/ssooidc => ../service/ssooidc

replace github.com/Enflick/aws-sdk-go-v2/service/sts => ../service/sts/
