module github.com/aws/aws-sdk-go-v2/service/sts

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.13
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.13
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.3
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.15
	github.com/aws/smithy-go v1.20.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
