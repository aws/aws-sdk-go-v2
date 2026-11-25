module github.com/aws/aws-sdk-go-v2/service/redshiftdata

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.48.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.13
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.13
	github.com/aws/smithy-go v1.28.10
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
