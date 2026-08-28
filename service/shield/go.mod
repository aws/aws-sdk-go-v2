module github.com/aws/aws-sdk-go-v2/service/shield

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.45.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.1
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.1
	github.com/aws/smithy-go v1.28.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
