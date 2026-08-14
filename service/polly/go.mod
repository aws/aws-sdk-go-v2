module github.com/aws/aws-sdk-go-v2/service/polly

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.5
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.17
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.36
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.36
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.16
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.36
	github.com/aws/smithy-go v1.27.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
