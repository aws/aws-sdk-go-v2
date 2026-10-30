module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restxmldataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.5
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.8
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.8
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.9
	github.com/aws/smithy-go v1.28.7
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.8 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
