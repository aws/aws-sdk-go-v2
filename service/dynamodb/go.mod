module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.29.1-0.20201113222241-726e4a15683d
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.3.1-0.20201113222241-726e4a15683d
	github.com/awslabs/smithy-go v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/
