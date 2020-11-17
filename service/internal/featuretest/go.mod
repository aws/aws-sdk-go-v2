module github.com/aws/aws-sdk-go-v2/service/internal/unittest

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v0.29.0
	github.com/google/go-cmp v0.5.2
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../../service/s3/
)
