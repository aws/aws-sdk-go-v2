module github.com/aws/aws-sdk-go-v2/service/kinesisvideosignaling

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.7
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.7
	github.com/aws/smithy-go v1.28.7
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
