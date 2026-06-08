module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.12
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.28
	github.com/aws/smithy-go v1.27.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
