module github.com/aws/aws-sdk-go-v2/service/codecommit

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.3
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.19
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.19
	github.com/aws/smithy-go v1.24.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
