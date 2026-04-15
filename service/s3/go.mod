module github.com/aws/aws-sdk-go-v2/service/s3

go 1.24

require (
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.13
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.21
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.21
	github.com/aws/smithy-go v1.24.2
	github.com/aws/aws-sdk-go-v2 v1.41.5
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.8
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.7
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
