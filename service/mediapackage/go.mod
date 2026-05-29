module github.com/aws/aws-sdk-go-v2/service/mediapackage

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.9
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.25
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.25
	github.com/aws/smithy-go v1.26.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
