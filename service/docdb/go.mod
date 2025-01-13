module github.com/aws/aws-sdk-go-v2/service/docdb

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.32.8
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.27
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.27
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.8
	github.com/aws/smithy-go v1.22.1
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
