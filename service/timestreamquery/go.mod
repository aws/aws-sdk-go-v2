module github.com/aws/aws-sdk-go-v2/service/timestreamquery

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.8
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.8
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.8
	github.com/aws/smithy-go v1.28.7
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
