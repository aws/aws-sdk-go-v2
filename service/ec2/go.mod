module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.18.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.33
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.27
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.27
	github.com/aws/smithy-go v1.13.5
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
