module github.com/aws/aws-sdk-go-v2/service/s3control

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.12.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.3
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.1.0
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.10.0
	github.com/aws/smithy-go v1.9.2-0.20220113020543-dec09760da63
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
