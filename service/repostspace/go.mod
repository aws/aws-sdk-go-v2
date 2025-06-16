module github.com/aws/aws-sdk-go-v2/service/repostspace

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.35
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.35
	github.com/aws/smithy-go v1.22.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
