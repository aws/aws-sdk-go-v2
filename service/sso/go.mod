module github.com/aws/aws-sdk-go-v2/service/sso

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.24
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.24
	github.com/aws/smithy-go v1.22.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
