module github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.6
	github.com/aws/aws-sdk-go-v2/config v1.33.14
	github.com/aws/aws-sdk-go-v2/credentials v1.20.14
	github.com/aws/aws-sdk-go-v2/service/s3 v1.114.7
	github.com/aws/aws-sdk-go-v2/service/sts v1.51.7
	github.com/aws/smithy-go v1.28.8
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.26 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.39.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.44.2 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/aws/aws-sdk-go-v2/service/signin => ../../../service/signin/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/
