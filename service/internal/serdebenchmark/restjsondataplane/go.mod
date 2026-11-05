module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restjsondataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.7
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.10
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.10
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.11
	github.com/aws/smithy-go v1.28.8
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.10 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
