module github.com/aws/aws-sdk-go-v2/service/timestreamquery

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.10.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.0.7
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.2.0
	github.com/aws/smithy-go v1.8.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
