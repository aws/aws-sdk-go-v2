module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.21.2
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.14
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.43
	github.com/aws/aws-sdk-go-v2/internal/endpoints/config v1.0.0-preview
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.37
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.6
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.15
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.38
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.37
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.15.6
	github.com/aws/smithy-go v1.15.0
	github.com/google/go-cmp v0.5.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/config => ../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/config => ../../internal/endpoints/config/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../service/sts/
