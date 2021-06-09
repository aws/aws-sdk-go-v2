module github.com/aws/aws-sdk-go-v2/service/timestreamwrite

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.6.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v0.0.0-00010101000000-000000000000
	github.com/aws/smithy-go v1.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
