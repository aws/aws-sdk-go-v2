module github.com/aws/aws-sdk-go-v2/service/deadline

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.41.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.18
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.18
	github.com/aws/smithy-go v1.24.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
