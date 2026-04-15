module github.com/aws/aws-sdk-go-v2/service/dynamodb

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.5
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.7
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.21
	github.com/aws/smithy-go v1.24.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => ../../service/internal/accept-encoding/

replace github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../service/internal/endpoint-discovery/
