module github.com/aws/aws-sdk-go-v2/service/sesv2

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.42.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.25
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.25
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.26
	github.com/aws/smithy-go v1.27.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/
