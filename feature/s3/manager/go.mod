module github.com/aws/aws-sdk-go-v2/feature/s3/manager

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.27.0
	github.com/aws/aws-sdk-go-v2/config v0.2.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.27.0
	github.com/awslabs/smithy-go v0.2.1-0.20201023021502-2e30e33fd215
	github.com/google/go-cmp v0.4.1
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../config/
	github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/
	github.com/aws/aws-sdk-go-v2/ec2imds => ../../../ec2imds/
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/
	github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/
)
