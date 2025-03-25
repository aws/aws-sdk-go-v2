module github.com/Enflick/aws-sdk-go-v2/service/dynamodb

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.27.0
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v1.3.7
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2
	github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.8
	github.com/Enflick/smithy-go v1.3.0
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
