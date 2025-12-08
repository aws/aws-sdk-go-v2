module github.com/aws/aws-sdk-go-v2/service/neptune

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.41.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16
	github.com/aws/smithy-go v1.24.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
