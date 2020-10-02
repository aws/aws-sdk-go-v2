module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20201001231852-1fc1ab173989
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20200930084954-897dfb99530c
	github.com/awslabs/smithy-go v0.1.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/
