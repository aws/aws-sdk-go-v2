module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.38.3
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.6
	github.com/aws/smithy-go v1.23.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
