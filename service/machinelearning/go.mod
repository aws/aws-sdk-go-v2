module github.com/aws/aws-sdk-go-v2/service/machinelearning

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.39.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.8
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.8
	github.com/aws/smithy-go v1.23.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
