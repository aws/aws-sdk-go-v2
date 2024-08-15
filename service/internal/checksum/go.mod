module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.30.4
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.18
	github.com/aws/smithy-go v1.20.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
