module github.com/aws/aws-sdk-go-v2/service/dynamodbstreams

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.23.4
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.7
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.7
	github.com/aws/smithy-go v1.18.1
	github.com/google/go-cmp v0.5.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
