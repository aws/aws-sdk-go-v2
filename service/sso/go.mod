module github.com/aws/aws-sdk-go-v2/service/sso

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.3
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.15
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.15
	github.com/aws/smithy-go v1.20.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
