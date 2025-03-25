module github.com/Enflick/aws-sdk-go-v2/service/s3

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/internal/v4a v1.3.7
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/service/internal/checksum v1.3.9
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20250325155711-0a4bf6fdbeb3
	github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared v1.17.7
	github.com/Enflick/smithy-go v1.3.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/checksum => ../../service/internal/checksum/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
