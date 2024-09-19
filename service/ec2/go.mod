module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.30.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.17
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.17
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.4
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.19
	github.com/aws/smithy-go v1.21.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
