module github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.11.2
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.0.0
	github.com/aws/aws-sdk-go-v2/credentials v1.6.5
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/
