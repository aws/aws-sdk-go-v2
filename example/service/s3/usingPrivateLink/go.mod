module github.com/Enflick/aws-sdk-go-v2/example/service/s3/usingPrivateLink

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.5.0
	github.com/Enflick/aws-sdk-go-v2/config v0.0.0-20250325221836-b4195dfa2eb5
	github.com/Enflick/aws-sdk-go-v2/service/s3 v1.54.3
	github.com/Enflick/aws-sdk-go-v2/service/s3control v1.44.11
)

require (
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/credentials v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/ini v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/v4a v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/checksum v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/sso v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/ssooidc v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/sts v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/smithy-go v1.3.0 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/config => ../../../../config/

replace github.com/Enflick/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/internal/ini => ../../../../internal/ini/

replace github.com/Enflick/aws-sdk-go-v2/internal/v4a => ../../../../internal/v4a/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared => ../../../../service/internal/s3shared/

replace github.com/Enflick/aws-sdk-go-v2/service/s3 => ../../../../service/s3/

replace github.com/Enflick/aws-sdk-go-v2/service/s3control => ../../../../service/s3control/

replace github.com/Enflick/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/Enflick/aws-sdk-go-v2/service/ssooidc => ../../../../service/ssooidc/

replace github.com/Enflick/aws-sdk-go-v2/service/sts => ../../../../service/sts/
