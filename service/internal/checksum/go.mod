module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.47.5
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.8
	github.com/aws/smithy-go v1.28.7
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
