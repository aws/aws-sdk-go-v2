module github.com/Enflick/aws-sdk-go-v2/service/internal/checksum

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v1.11.9
	github.com/Enflick/smithy-go v1.3.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
