module github.com/aws/aws-sdk-go-v2/service/kinesis/internal/testing

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.1
	github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting v1.3.5
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.27.1
	github.com/aws/smithy-go v1.20.1
	github.com/google/go-cmp v0.5.8
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting => ../../../../service/internal/eventstreamtesting/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/kinesis => ../../../../service/kinesis/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/ssooidc => ../../../../service/ssooidc/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/
