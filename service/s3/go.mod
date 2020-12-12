module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.3.1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.1.2
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.2
	github.com/awslabs/smithy-go v0.4.1-0.20201208232924-b8cdbaa577ff
	github.com/google/go-cmp v0.5.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
