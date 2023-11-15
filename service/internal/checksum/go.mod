module github.com/aws/aws-sdk-go-v2/service/internal/checksum

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.23.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.3
	github.com/aws/smithy-go v1.17.0
	github.com/google/go-cmp v0.5.8
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
