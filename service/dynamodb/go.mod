module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.14
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.35
	github.com/aws/smithy-go v1.14.2
	github.com/google/go-cmp v0.5.8
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
