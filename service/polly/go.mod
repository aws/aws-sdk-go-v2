module github.com/aws/aws-sdk-go-v2/service/polly

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.8
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.26
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.11
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.11
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.24
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.11
	github.com/aws/smithy-go v1.28.9
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
