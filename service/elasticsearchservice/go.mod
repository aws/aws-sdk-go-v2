module github.com/aws/aws-sdk-go-v2/service/elasticsearchservice

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.8
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.27
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.27
	github.com/aws/smithy-go v1.22.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
