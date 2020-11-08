module github.com/aws/aws-sdk-go-v2/feature/s3/manager

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.0
	github.com/aws/aws-sdk-go-v2/config v0.2.2
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v0.29.0
	github.com/awslabs/smithy-go v0.3.1-0.20201108010311-62c2a93810b4
	github.com/google/go-cmp v0.4.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/ec2imds => ../../../ec2imds/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/
