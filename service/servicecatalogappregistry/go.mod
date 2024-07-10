module github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.14
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.14
	github.com/aws/smithy-go v1.20.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
