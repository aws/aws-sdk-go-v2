module github.com/aws/aws-sdk-go-v2/service/neptune

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.0
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.0
	github.com/aws/smithy-go v1.20.0
	github.com/google/go-cmp v0.5.8
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
