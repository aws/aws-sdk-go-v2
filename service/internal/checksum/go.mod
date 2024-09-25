module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.31.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.20
	github.com/aws/smithy-go v1.21.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
