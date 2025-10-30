module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.39.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.12
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.12
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.2
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.12
	github.com/aws/smithy-go v1.23.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
