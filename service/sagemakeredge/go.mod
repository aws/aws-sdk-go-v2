module github.com/aws/aws-sdk-go-v2/service/sagemakeredge

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.14.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.5
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.3.0
	github.com/aws/smithy-go v1.11.1-0.20220308004241-26b41c98827d
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
