module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.43.7
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38
	github.com/aws/smithy-go v1.27.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
