module github.com/aws/aws-sdk-go-v2/service/s3

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.1-0.20210108183639-b6b5057e2ab1
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.4.1-0.20210108183639-b6b5057e2ab1
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.2.1-0.20210108183639-b6b5057e2ab1
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v0.4.1-0.20210108183639-b6b5057e2ab1
	github.com/aws/smithy-go v0.5.1-0.20210107224202-ae5323020d60
	github.com/google/go-cmp v0.5.4
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/

replace github.com/aws/aws-sdk-go-v2/service/internal/s3shared => ../../service/internal/s3shared/
