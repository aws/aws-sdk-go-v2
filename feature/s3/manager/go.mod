module github.com/aws/aws-sdk-go-v2/feature/s3/manager

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201217001905-4acf9c65b2d1
	github.com/aws/aws-sdk-go-v2/config v0.3.0
	github.com/aws/aws-sdk-go-v2/service/s3 v0.30.0
	github.com/awslabs/smithy-go v0.4.1-0.20201216214517-20e212c92831
	github.com/google/go-cmp v0.5.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../config/

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/aws/aws-sdk-go-v2/credentials => ../../../credentials/

replace github.com/aws/aws-sdk-go-v2/feature/ec2/imds => ../../../feature/ec2/imds/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../service/sts/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
