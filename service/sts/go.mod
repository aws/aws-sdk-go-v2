module github.com/aws/aws-sdk-go-v2/service/sts

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.7
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.38
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.38
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38
	github.com/aws/smithy-go v1.27.10
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
