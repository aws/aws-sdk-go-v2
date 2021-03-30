module github.com/aws/aws-sdk-go-v2/service/docdb

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.3.0
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.0.4
	github.com/aws/smithy-go v1.2.1-0.20210330191137-b290458ce30a
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
