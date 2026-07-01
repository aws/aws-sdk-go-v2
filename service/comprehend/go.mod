module github.com/aws/aws-sdk-go-v2/service/comprehend

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.42.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.30
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.30
	github.com/aws/smithy-go v1.27.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
