module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restxmldataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.48.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.14
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.14
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.15
	github.com/aws/smithy-go v1.28.11
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.14 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
