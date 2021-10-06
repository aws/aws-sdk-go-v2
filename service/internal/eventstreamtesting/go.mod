module github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.9.1
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/credentials v1.4.2
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/
