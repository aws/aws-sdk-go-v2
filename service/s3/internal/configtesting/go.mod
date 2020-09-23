module github.com/aws/aws-sdk-go-v2/service/s3/internal/configtesting

go 1.15

require (
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200923001701-7b95ccd95ed9
	github.com/aws/aws-sdk-go-v2/service/s3 v0.0.0-00010101000000-000000000000
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config
	github.com/aws/aws-sdk-go-v2/service/s3 => ../../
)
