module github.com/aws/aws-sdk-go-v2/service/rds

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.35
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.35
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.16
	github.com/aws/smithy-go v1.22.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
