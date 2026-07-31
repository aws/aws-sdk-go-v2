module github.com/aws/aws-sdk-go-v2/service/internal/serdebenchmark/restjsondataplane

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.33
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.33
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.26
	github.com/aws/smithy-go v1.27.6
)

require github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.33 // indirect

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/service/internal/checksum => ../../../../service/internal/checksum/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/
