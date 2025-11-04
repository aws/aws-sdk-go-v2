module github.com/aws/aws-sdk-go-v2/service/inspector2

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.13
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.13
	github.com/aws/smithy-go v1.23.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
