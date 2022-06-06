module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.16.4
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.11
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.5
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.2
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.1
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.6
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.5
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.5
	github.com/aws/smithy-go v1.11.3-0.20220606214609-8c1eac595edb
	github.com/google/go-cmp v0.5.8
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
