module github.com/aws/aws-sdk-go-v2/service/health

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.37.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.2
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.2
	github.com/aws/smithy-go v1.22.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
