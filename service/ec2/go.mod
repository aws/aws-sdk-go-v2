module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10
	github.com/aws/smithy-go v1.19.0
	github.com/google/go-cmp v0.5.8
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
