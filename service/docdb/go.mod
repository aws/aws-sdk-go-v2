module github.com/aws/aws-sdk-go-v2/service/docdb

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.39.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.7
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.7
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.7
	github.com/aws/smithy-go v1.23.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
