module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.21
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.21
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.7
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.21
	github.com/aws/smithy-go v1.24.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
