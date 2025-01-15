module github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.33.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28
	github.com/aws/smithy-go v1.22.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
