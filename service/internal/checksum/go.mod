module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.33.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9
	github.com/aws/smithy-go v1.22.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
