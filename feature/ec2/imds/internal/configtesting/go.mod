module github.com/aws/aws-sdk-go-v2/feature/ec2/imds/internal/configtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/config v1.4.1
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.2.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../../../config/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/internal/ini => ../../../../../internal/ini/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/sso => ../../../../../service/sso/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../../service/sts/
