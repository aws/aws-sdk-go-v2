module github.com/aws/aws-sdk-go-v2/service/timestreamwrite

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.9.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.0.6
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.1.2
	github.com/aws/smithy-go v1.8.1-0.20211020161917-191f636375d1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
