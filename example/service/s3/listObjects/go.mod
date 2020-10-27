module github.com/aws/aws-sdk-go-v2/example/service/s3/listObjects

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/config v0.2.1
	github.com/aws/aws-sdk-go-v2/service/s3 v0.28.0
)

replace github.com/aws/aws-sdk-go-v2/config => ../../../../config/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../../service/s3/

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/aws/aws-sdk-go-v2/ec2imds => ../../../../ec2imds/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../../service/internal/s3shared/
