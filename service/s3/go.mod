module github.com/aws/aws-sdk-go-v2/service/s3

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.40.1
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.15
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.15
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.15
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.6
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.15
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.15
	github.com/aws/smithy-go v1.24.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
