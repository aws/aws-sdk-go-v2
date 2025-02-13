module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.32
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.32
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.13
	github.com/aws/smithy-go v1.22.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
