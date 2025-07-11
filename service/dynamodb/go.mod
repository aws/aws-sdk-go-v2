module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.36
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.36
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.4
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.17
	github.com/aws/smithy-go v1.22.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
