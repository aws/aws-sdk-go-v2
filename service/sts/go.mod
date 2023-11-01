module github.com/aws/aws-sdk-go-v2/service/sts

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.22.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.1
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.1
	github.com/aws/smithy-go v1.16.0
	github.com/google/go-cmp v0.5.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
