module github.com/aws/aws-sdk-go-v2/service/ssm

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.33
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.33
	github.com/aws/smithy-go v1.22.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
