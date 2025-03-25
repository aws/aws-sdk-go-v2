module github.com/Enflick/aws-sdk-go-v2/service/kinesis/internal/testing

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2
	github.com/Enflick/aws-sdk-go-v2/service/internal/eventstreamtesting v1.4.3
	github.com/Enflick/aws-sdk-go-v2/service/kinesis v1.27.8
	github.com/Enflick/smithy-go v1.3.0
)

require (
	github.com/Enflick/aws-sdk-go-v2/credentials v1.17.16 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/Enflick/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/eventstreamtesting => ../../../../service/internal/eventstreamtesting/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/kinesis => ../../../../service/kinesis/

replace github.com/Enflick/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/Enflick/aws-sdk-go-v2/service/ssooidc => ../../../../service/ssooidc/

replace github.com/Enflick/aws-sdk-go-v2/service/sts => ../../../../service/sts/
