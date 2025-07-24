module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.6
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.18
	github.com/aws/smithy-go v1.22.5
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
