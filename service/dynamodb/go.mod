module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.7
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.38
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.38
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.15
	github.com/aws/smithy-go v1.27.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
