module github.com/aws/aws-sdk-go-v2/service/s3/internal/configtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/config v0.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.3-0.20201217001905-4acf9c65b2d1
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config/
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../../../service/s3/
)

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
