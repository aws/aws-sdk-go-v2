module github.com/aws/aws-sdk-go-v2/service/ecrpublic

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.19.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.36
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.30
	github.com/aws/smithy-go v1.13.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
