module github.com/aws/aws-sdk-go-v2/service/resiliencehubv2

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.35
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.35
	github.com/aws/smithy-go v1.27.6
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
