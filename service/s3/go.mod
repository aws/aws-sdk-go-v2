module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.30.1-0.20201221101722-677dd4a81dad
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.3.2-0.20201217001905-4acf9c65b2d1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.1.3-0.20201217001905-4acf9c65b2d1
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.3.3-0.20201217001905-4acf9c65b2d1
	github.com/aws/smithy-go v0.4.1-0.20201222001052-74df8ddd8c79
	github.com/google/go-cmp v0.5.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
