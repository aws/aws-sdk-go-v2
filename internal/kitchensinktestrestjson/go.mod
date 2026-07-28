module github.com/aws/aws-sdk-go-v2/internal/kitchensinktestrestjson

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.32
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.32
	github.com/aws/smithy-go v1.27.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
