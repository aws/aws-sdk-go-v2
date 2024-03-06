module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.1
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.3
	github.com/aws/smithy-go v1.20.1
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
