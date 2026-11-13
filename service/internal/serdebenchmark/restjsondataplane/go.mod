module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restjsondataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.48.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.12
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.12
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.13
	github.com/aws/smithy-go v1.28.9
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.12 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
