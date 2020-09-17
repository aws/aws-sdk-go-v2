module github.com/aws/aws-sdk-go-v2/service/s3/internal/configtesting

go 1.15

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/config => ../../../../config

replace github.com/aws/aws-sdk-go-v2/service/s3 => ../../

require (
	github.com/aws/aws-sdk-go-v2/config v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go-v2/service/s3 v0.0.0-00010101000000-000000000000
)
