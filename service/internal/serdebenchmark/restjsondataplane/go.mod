module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restjsondataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.5
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.5
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.6
	github.com/aws/smithy-go v1.28.5
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.5 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
