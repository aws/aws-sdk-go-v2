module github.com/aws/aws-sdk-go-v2/service/marketplaceentitlementservice

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.17.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.25
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.19
	github.com/aws/smithy-go v1.13.4
	github.com/stretchr/testify v1.8.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
