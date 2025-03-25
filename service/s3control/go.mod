module github.com/Enflick/aws-sdk-go-v2/service/s3control

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.27.0
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v1.3.7
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7
	github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared v1.17.7
	github.com/Enflick/smithy-go v1.3.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
