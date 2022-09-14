module github.com/aws/aws-sdk-go-v2/service/sqs

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.16.15
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.22
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.16
	github.com/aws/smithy-go v1.13.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
