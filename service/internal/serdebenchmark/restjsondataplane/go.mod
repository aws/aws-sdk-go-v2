module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restjsondataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.3
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.3
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.4
	github.com/aws/smithy-go v1.28.1
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.3 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
