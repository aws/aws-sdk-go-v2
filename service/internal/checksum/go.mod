module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.40.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.14
	github.com/aws/smithy-go v1.24.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
